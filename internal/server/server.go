package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
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

	"github.com/gmhafiz/go8/configs"
	_ "github.com/gmhafiz/go8/docs"
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
	cfg        *configs.Configs
	db         *sqlx.DB
	router     *chi.Mux
	httpServer *http.Server
	validator  *validator.Validate
	cache      *redis.Client
}

type Options func(opts *Server) error

func New(opts ...Options) *Server {
	s := &Server{}
	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return s
}

func WithConfig() func(s *Server) error {
	return func(s *Server) error {
		s.cfg = configs.New()
		return nil
	}
}

func WithRouter() func(s *Server) error {
	return func(s *Server) error {
		s.router = chi.NewRouter()
		return nil
	}
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
	s.StartSwagger()
}

func (s *Server) StartSwagger() {
	if s.cfg.Api.RunSwagger {
		swaggerServer(s.router)
	}
}

func (s *Server) newConfig() {
	s.cfg = configs.New()
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
		Addr:           s.cfg.Api.Host + ":" + s.cfg.Api.Port,
		Handler:        s.router,
		ReadTimeout:    s.cfg.Api.ReadTimeout,
		WriteTimeout:   s.cfg.Api.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
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
		log.Printf("Serving at %s:%s\n", s.cfg.Api.Host, s.cfg.Api.Port)
		err := s.httpServer.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), s.cfg.Api.IdleTimeout*time.Second)
	defer shutdown()

	_ = s.DB().Close()
	_ = s.httpServer.Shutdown(ctx)
}

func (s *Server) Config() *configs.Configs {
	return s.cfg
}

func (s *Server) DB() *sqlx.DB {
	return s.db
}

func (s *Server) Cache() *redis.Client {
	return s.cache
}

func swaggerServer(router *chi.Mux) {
	fileServer := http.FileServer(http.Dir(swaggerDocsAssetPath))
	router.Handle("/swagger/*", http.StripPrefix("/swagger", fileServer))
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
