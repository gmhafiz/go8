package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"

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
	version    string
	cfg        *config.Config
	db         *sqlx.DB
	router     *chi.Mux
	httpServer *http.Server
	validator  *validator.Validate
	cache      *redis.Client
	ent        *gen.Client
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
	s.version = version
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

	dsn := fmt.Sprintf("%s://%s/%s?sslmode=%s&user=%s&password=%s",
		s.cfg.Database.Driver,
		s.cfg.Database.Host,
		s.cfg.Database.Name,
		s.cfg.Database.SslMode,
		s.cfg.Database.User,
		s.cfg.Database.Pass)
	client, err := gen.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	s.ent = client
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

func (s *Server) PrintAllRegisteredRoutes() {
	walkFunc := func(method string, path string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%-7s %s\n", method, path)

		return nil
	}
	if err := chi.Walk(s.router, walkFunc); err != nil {
		fmt.Print(err)
	}
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

	return s.httpServer.Shutdown(ctx)
}
