package server

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"syscall"
	"time"

	"encoding/json"
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/jwalton/gchalk"
	"golang.org/x/mod/modfile"

	"github.com/gmhafiz/go8/config"
	_ "github.com/gmhafiz/go8/docs"
	"github.com/gmhafiz/go8/ent/gen"
	"github.com/gmhafiz/go8/internal/middleware"
	db "github.com/gmhafiz/go8/third_party/database"
	redisLib "github.com/gmhafiz/go8/third_party/redis"
	"github.com/gmhafiz/go8/third_party/validate"
)

const (
	databaseMigrationPath = "file://database/migrations/"
	swaggerDocsAssetPath  = "./docs/"
)

type Server struct {
	Version string
	cfg     *config.Config

	db    *sqlx.DB
	ent   *gen.Client
	cache *redis.Client

	validator  *validator.Validate
	router     *chi.Mux
	httpServer *http.Server

	Domain
}

type Options func(opts *Server) error

func defaultServer() *Server {
	return &Server{
		cfg:    config.New(),
		router: chi.NewRouter(),
	}
}

func New(opts ...Options) *Server {
	s := defaultServer()

	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return s
}

func (s *Server) Init(version string) {
	s.Version = version
	s.newConfig()
	s.newRedis()
	s.newDatabase()
	s.newValidator()
	s.newRouter()
	s.setGlobalMiddleware()
	s.InitDomains()
}

func (s *Server) newConfig() {
	s.cfg = config.New()
}

func (s *Server) newRedis() {
	s.cache = redisLib.New(s.cfg.Cache)
}

func (s *Server) newDatabase() {
	if s.cfg.Database.Driver == "" {
		log.Fatal("please fill in database credentials in .env file or set in environment variable")
	}
	s.db = db.NewSqlx(s.cfg)
	s.db.SetMaxOpenConns(s.cfg.Database.MaxConnectionPool)
	s.db.SetMaxIdleConns(s.cfg.Database.MaxIdleConnections)
	s.db.SetConnMaxLifetime(s.cfg.Database.ConnectionsMaxLifeTime)

	dsn := fmt.Sprintf("postgres://%s/%s?sslmode=%s&user=%s&password=%s",
		s.cfg.Database.Host,
		s.cfg.Database.Name,
		s.cfg.Database.SslMode,
		s.cfg.Database.User,
		s.cfg.Database.Pass,
	)

	s.newEnt(dsn)
}

func (s *Server) newValidator() {
	s.validator = validate.New()
}

func (s *Server) newRouter() {
	s.router = chi.NewRouter()
}

func (s *Server) setGlobalMiddleware() {
	s.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error": "endpoint not found"}`))
	})
	s.router.Use(middleware.Json)
	s.router.Use(middleware.AuthN())
	s.router.Use(middleware.Audit)
	s.router.Use(middleware.CacheByURL)
	s.router.Use(middleware.CORS)
	if s.cfg.Api.RequestLog {
		s.router.Use(chiMiddleware.Logger)
	}
	s.router.Use(middleware.Recovery)
}

func (s *Server) Migrate() {
	log.Println("migrating...")

	var driver database.Driver
	switch s.cfg.Database.Driver {
	case "postgres":
		d, err := postgres.WithInstance(s.DB().DB, &postgres.Config{})
		if err != nil {
			log.Fatalf("error instantiating database: %v", err)
		}
		driver = d
	case "mysql":
		d, err := mysql.WithInstance(s.DB().DB, &mysql.Config{})
		if err != nil {
			log.Fatalf("error instantiating database: %v", err)
		}
		driver = d
	}

	m, err := migrate.NewWithDatabaseInstance(
		databaseMigrationPath, s.cfg.Database.Driver, driver,
	)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	err = m.Up()
	if err != nil {
		if err != migrate.ErrNoChange {
			log.Panicf("error migrating: %v", err)
		}
	}

	log.Println("done migration.")
}

func (s *Server) Run() {
	s.httpServer = &http.Server{
		Addr:    s.cfg.Api.Host + ":" + s.cfg.Api.Port,
		Handler: s.router,
	}

	fmt.Println(`            .,*/(#####(/*,.                               .,*((###(/*.
        .*(%%%%%%%%%%%%%%#/.                           .*#%%%%####%%%%#/.
      ./#%%%%#(/,,...,,***.           .......          *#%%%#*.   ,(%%%#/.
     .(#%%%#/.                    .*(#%%%%%%%##/,.     ,(%%%#*    ,(%%%#*.
    .*#%%%#/.    ..........     .*#%%%%#(/((#%%%%(,     ,/#%%%#(/#%%%#(,
    ./#%%%(*    ,#%%%%%%%%(*   .*#%%%#*     .*#%%%#,      *(%%%%%%%#(,.
    ./#%%%#*    ,(((##%%%%(*   ,/%%%%/.      .(%%%#/   .*#%%%#(*/(#%%%#/,
     ,#%%%#(.        ,#%%%(*   ,/%%%%/.      .(%%%#/  ,/%%%#/.    .*#%%%(,
      *#%%%%(*.      ,#%%%(*   .*#%%%#*     ./#%%%#,  ,(%%%#*      .(%%%#*
       ,(#%%%%%##(((##%%%%(*    .*#%%%%#(((##%%%%(,   .*#%%%##(///(#%%%#/.
         .*/###%%%%%%%###(/,      .,/##%%%%%##(/,.      .*(##%%%%%%##(*,
              .........                ......                .......`)
	go func() {
		start(s)
	}()

	_ = gracefulShutdown(context.Background(), s)
}

func (s *Server) Config() *config.Config {
	return s.cfg
}

func (s *Server) DB() *sqlx.DB {
	return s.db
}

func (s *Server) Cache() *redis.Client {
	return s.cache
}

// PrintAllRegisteredRoutes prints all registered routes from Chi router.
// definitely can be an extension to the router instead.
func (s *Server) PrintAllRegisteredRoutes(exceptions ...string) {
	exceptions = append(exceptions, "/swagger")

	walkFunc := func(method string, path string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {

		for _, val := range exceptions {
			if strings.HasPrefix(path, val) {
				return nil
			}
		}

		switch method {
		case "GET":
			fmt.Printf("%s", gchalk.Green(fmt.Sprintf("%-8s", method)))
		case "POST", "PUT", "PATCH":
			fmt.Printf("%s", gchalk.Yellow(fmt.Sprintf("%-8s", method)))
		case "DELETE":
			fmt.Printf("%s", gchalk.Red(fmt.Sprintf("%-8s", method)))
		default:
			fmt.Printf("%s", gchalk.White(fmt.Sprintf("%-8s", method)))
		}

		//fmt.Printf("%-25s %60s\n", path, getHandler(getModName(), handler))
		fmt.Printf("%s", strPad(path, 25, "-", "RIGHT"))
		fmt.Printf("%s\n", strPad(getHandler(getModName(), handler), 60, "-", "LEFT"))

		return nil
	}
	if err := chi.Walk(s.router, walkFunc); err != nil {
		fmt.Print(err)
	}

	if s.cfg.Api.RunSwagger {
		fmt.Printf("%s", gchalk.Green(fmt.Sprintf("%-8s", "GET")))
		fmt.Printf("/swagger\n")
	}
}

func (s *Server) newEnt(dsn string) {
	client, err := gen.Open(dialect.Postgres, dsn)
	if err != nil {
		log.Fatal(err)
	}

	client.Use(func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, mutation ent.Mutation) (ent.Value, error) {
			meta, ok := ctx.Value(middleware.AuditID).(middleware.Event)
			if !ok {
				return next.Mutate(ctx, mutation)
			}

			val, err := next.Mutate(ctx, mutation)

			meta.Table = mutation.Type()
			meta.Action = middleware.Action(mutation.Op().String())

			newValues, _ := json.Marshal(val)
			meta.NewValues = string(newValues)
			log.Println(meta)

			return val, err
		})
	})

	s.ent = client
}

// StrPad returns the input string padded on the left, right or both sides using padType to the specified padding length padLength.
//
// Example:
// input := "Codes";
// StrPad(input, 10, " ", "RIGHT")        // produces "Codes     "
// StrPad(input, 10, "-=", "LEFT")        // produces "=-=-=Codes"
// StrPad(input, 10, "_", "BOTH")         // produces "__Codes___"
// StrPad(input, 6, "___", "RIGHT")       // produces "Codes_"
// StrPad(input, 3, "*", "RIGHT")         // produces "Codes"
// taken from // https://gist.github.com/asessa/3aaec43d93044fc42b7c6d5f728cb039
func strPad(input string, padLength int, padString string, padType string) string {
	var output string

	inputLength := len(input)
	padStringLength := len(padString)

	if inputLength >= padLength {
		return input
	}

	repeat := math.Ceil(float64(1) + (float64(padLength-padStringLength))/float64(padStringLength))

	switch padType {
	case "RIGHT":
		output = input + strings.Repeat(padString, int(repeat))
		output = output[:padLength]
	case "LEFT":
		output = strings.Repeat(padString, int(repeat)) + input
		output = output[len(output)-padLength:]
	case "BOTH":
		length := (float64(padLength - inputLength)) / float64(2)
		repeat = math.Ceil(length / float64(padStringLength))
		output = strings.Repeat(padString, int(repeat))[:int(math.Floor(float64(length)))] + input + strings.Repeat(padString, int(repeat))[:int(math.Ceil(float64(length)))]
	}

	return output
}

func getHandler(projectName string, handler http.Handler) (funcName string) {
	// https://github.com/go-chi/chi/issues/424
	funcName = runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	base := filepath.Base(funcName)

	nameSplit := strings.Split(funcName, "")
	names := nameSplit[len(projectName):]
	path := strings.Join(names, "")

	pathSplit := strings.Split(path, "/")
	path = strings.Join(pathSplit[:len(pathSplit)-1], "/")

	sFull := strings.Split(base, ".")
	s := sFull[len(sFull)-1:]

	s = strings.Split(s[0], "")
	if len(s) <= 4 && len(sFull) >= 3 {
		s = sFull[len(sFull)-3 : len(sFull)-2]
		return "@" + gchalk.Blue(strings.Join(s, ""))
	}
	s = s[:len(s)-3]
	funcName = strings.Join(s, "")

	return path + "@" + gchalk.Blue(funcName)
}

// adapted from https://stackoverflow.com/a/63393712/1033134
func getModName() string {
	goModBytes, err := os.ReadFile("go.mod")
	if err != nil {
		os.Exit(0)
	}
	return modfile.ModulePath(goModBytes)
}

func start(s *Server) {
	log.Printf("Serving at %s:%s\n", s.cfg.Api.Host, s.cfg.Api.Port)
	err := s.httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func gracefulShutdown(ctx context.Context, s *Server) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-quit

	log.Println("Shutting down...")

	ctx, shutdown := context.WithTimeout(ctx, s.Config().Api.GracefulTimeout*time.Second)
	defer shutdown()

	_ = s.DB().Close()
	s.cache.Shutdown(ctx)

	return s.httpServer.Shutdown(ctx)
}
