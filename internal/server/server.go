package server

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
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

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"github.com/gmhafiz/scs/v2"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	_ "github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	"github.com/jwalton/gchalk"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
	"golang.org/x/mod/modfile"

	"github.com/gmhafiz/go8/config"
	"github.com/gmhafiz/go8/database"
	//_ "github.com/gmhafiz/go8/docs"
	"github.com/gmhafiz/go8/ent/gen"
	"github.com/gmhafiz/go8/internal/middleware"
	db "github.com/gmhafiz/go8/third_party/database"
	"github.com/gmhafiz/go8/third_party/postgresstore"
	redisLib "github.com/gmhafiz/go8/third_party/redis"
	"github.com/gmhafiz/go8/third_party/validate"
)

//go:embed docs/*
var swaggerDocsAssetPath embed.FS

type Server struct {
	Version string
	cfg     *config.Config

	db   *sql.DB
	sqlx *sqlx.DB
	ent  *gen.Client

	cache   *redis.Client
	cluster *redis.ClusterClient

	session       *scs.SessionManager
	sessionCloser *postgresstore.PostgresStore

	validator *validator.Validate
	cors      *cors.Cors
	router    *chi.Mux

	httpServer *http.Server
}

type Options func(opts *Server) error

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

func WithVersion(version string) Options {
	return func(opts *Server) error {
		log.Printf("Starting API version: %s\n", version)
		opts.Version = version
		return nil
	}
}

func defaultServer() *Server {
	return &Server{
		cfg:    config.New(),
		router: chi.NewRouter(),
	}
}

func (s *Server) Init() {
	s.setCors()
	s.newRedis()
	s.NewDatabase()
	s.newValidator()
	s.newAuthentication()
	s.newRouter()
	s.setGlobalMiddleware()
	s.InitDomains()
}

func (s *Server) setCors() {
	s.cors = cors.New(
		cors.Options{
			AllowedOrigins: s.cfg.Cors.AllowedOrigins,
			AllowedMethods: []string{
				http.MethodHead,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
			},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
		})
}

func (s *Server) newRedis() {
	if !s.cfg.Cache.Enable {
		return
	}

	if len(s.cfg.Cache.Hosts) > 0 {
		s.cluster = redisLib.NewCluster(s.cfg.Cache)
	} else {
		s.cache = redisLib.New(s.cfg.Cache)
	}
}

func (s *Server) NewDatabase() {
	if s.cfg.Database.Driver == "" {
		log.Fatal("please fill in database credentials in .env file or set in environment variable")
	}

	s.sqlx = db.NewSqlx(s.cfg.Database)
	s.sqlx.SetMaxOpenConns(s.cfg.Database.MaxConnectionPool)
	s.sqlx.SetMaxIdleConns(s.cfg.Database.MaxIdleConnections)
	s.sqlx.SetConnMaxLifetime(s.cfg.Database.ConnectionsMaxLifeTime)

	dsn := fmt.Sprintf("postgres://%s:%d/%s?sslmode=%s&user=%s&password=%s",
		s.cfg.Database.Host,
		s.cfg.Database.Port,
		s.cfg.Database.Name,
		s.cfg.Database.SslMode,
		s.cfg.Database.User,
		s.cfg.Database.Pass,
	)
	s.db = s.sqlx.DB
	s.newEnt(dsn)
}

func (s *Server) newValidator() {
	s.validator = validate.New()
}

func (s *Server) newAuthentication() {
	manager := scs.New()
	manager.Store = postgresstore.New(s.sqlx.DB)
	manager.CtxStore = postgresstore.New(s.sqlx.DB)
	manager.Lifetime = s.cfg.Session.Duration
	manager.Cookie.Name = s.cfg.Session.Name
	manager.Cookie.Domain = s.cfg.Session.Domain
	manager.Cookie.HttpOnly = s.cfg.Session.HttpOnly
	manager.Cookie.Path = s.cfg.Session.Path
	manager.Cookie.Persist = true
	manager.Cookie.SameSite = http.SameSite(s.cfg.Session.SameSite)
	manager.Cookie.Secure = s.cfg.Session.Secure

	s.sessionCloser = postgresstore.NewWithCleanupInterval(s.sqlx.DB, 30*time.Minute)

	s.session = manager
}

func (s *Server) newRouter() {
	s.router = chi.NewRouter()
}

func (s *Server) setGlobalMiddleware() {
	s.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message": "endpoint not found"}`))
	})
	s.router.Use(s.cors.Handler)
	s.router.Use(middleware.Json)
	s.router.Use(middleware.LoadAndSave(s.session))
	s.router.Use(middleware.Audit)
	if s.cfg.Api.RequestLog {
		s.router.Use(chiMiddleware.Logger)
	}
	s.router.Use(middleware.Recovery)
}

func (s *Server) Migrate() {
	log.Println("migrating...")

	var databaseUrl string
	switch s.cfg.Database.Driver {
	case "postgres":
		databaseUrl = fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s",
			s.cfg.Database.Driver,
			s.cfg.Database.User,
			s.cfg.Database.Pass,
			s.cfg.Database.Host,
			s.cfg.Database.Port,
			s.cfg.Database.Name,
			s.cfg.Database.SslMode,
		)
	case "mysql":
		databaseUrl = fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true",
			s.cfg.Database.User,
			s.cfg.Database.Pass,
			s.cfg.Database.Host,
			s.cfg.Database.Port,
			s.cfg.Database.Name,
		)
	}

	migrator := database.Migrator(s.db, database.WithDSN(databaseUrl))
	migrator.Up()

	log.Println("done migration.")
}

func (s *Server) Run() {
	s.httpServer = &http.Server{
		Addr:              s.cfg.Api.Host + ":" + s.cfg.Api.Port,
		Handler:           s.router,
		ReadHeaderTimeout: s.cfg.Api.ReadHeaderTimeout,
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
			meta, ok := ctx.Value(middleware.KeyAuditID).(middleware.Event)
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
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Shutting down...")

	ctx, shutdown := context.WithTimeout(ctx, s.Config().Api.GracefulTimeout*time.Second)
	defer shutdown()

	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		log.Println(err)
	}
	s.closeResources(ctx)

	return nil
}

func (s *Server) closeResources(ctx context.Context) {
	_ = s.sqlx.Close()
	_ = s.ent.Close()
	s.cluster.Shutdown(ctx)
	s.cache.Shutdown(ctx)
	s.sessionCloser.StopCleanup()
}
