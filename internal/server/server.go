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
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"

	"github.com/gmhafiz/go8/configs"
	"github.com/gmhafiz/go8/internal/middleware"
	db "github.com/gmhafiz/go8/third_party/database"
	"github.com/gmhafiz/go8/third_party/validate"
)

const (
	databaseMigrationPath = "file://database/migrations/"
	swaggerDocsAssetPath  = "./docs"
)

type Server struct {
	cfg        *configs.Configs
	db         *sqlx.DB
	router     *chi.Mux
	httpServer *http.Server
	validator  *validator.Validate
}

func New(version string) *Server {
	log.Printf("Starting API version: %s\n", version)
	return &Server{}
}

func (s *Server) Init() {
	s.newConfig()
	s.newDatabase()
	s.newValidator()
	s.newRouter()
	s.setGlobalMiddleware()
	s.initDomains()
	s.startSwagger()
}

func (s *Server) startSwagger() {
	if s.cfg.Api.RunSwagger {
		swaggerServer(s.router)
	}
}

func (s *Server) newConfig() {
	s.cfg = configs.New()
}

func (s *Server) newDatabase() {
	if s.cfg.Database.Driver == "" {
		log.Fatal("please fill in database credentials in .env file")
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
	s.router.Use(middleware.Json)
	s.router.Use(middleware.Cors)
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

func (s *Server) Run() error {
	s.httpServer = &http.Server{
		Addr:           s.cfg.Api.Host.String() + ":" + s.cfg.Api.Port,
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
		printAllRegisteredRoutes(s.router)
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

	return s.httpServer.Shutdown(ctx)
}

func (s *Server) Config() *configs.Configs {
	return s.cfg
}

func (s *Server) DB() *sqlx.DB {
	return s.db
}

func swaggerServer(router *chi.Mux) {
	fs := http.FileServer(http.Dir(swaggerDocsAssetPath))

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(swaggerDocsAssetPath + r.RequestURI); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	})
}

func printAllRegisteredRoutes(router *chi.Mux) {
	walkFunc := func(method string, path string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("path: %s method: %s ", path, method)
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Print(err)
	}
}
