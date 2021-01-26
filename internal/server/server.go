package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"

	"github.com/gmhafiz/go8/configs"
	bookHandler "github.com/gmhafiz/go8/internal/domain/book/handler/http"
	bookRepo "github.com/gmhafiz/go8/internal/domain/book/repository/postgres"
	bookUseCase "github.com/gmhafiz/go8/internal/domain/book/usecase"
	healthHandler "github.com/gmhafiz/go8/internal/domain/health/handler/http"
	healthRepo "github.com/gmhafiz/go8/internal/domain/health/repository/postgres"
	healthUseCase "github.com/gmhafiz/go8/internal/domain/health/usecase"
	"github.com/gmhafiz/go8/internal/middleware"
	"github.com/gmhafiz/go8/third_party/database"
)

type Server struct {
	cfg        *configs.Configs
	db         *sqlx.DB
	router     *chi.Mux
	httpServer *http.Server
}

func New(version string) *Server {
	log.Printf("staring API version: %s\n", version)
	return &Server{}
}

func (s *Server) Init() {
	s.newConfig()
	s.newDatabase()
	s.newRouter()
	s.initDomains()
}

func (s *Server) newConfig() {
	s.cfg = configs.New()
}

func (s *Server) newDatabase() {
	s.db = database.NewSqlx(s.cfg)
}

func (s *Server) newRouter() {
	s.router = chi.NewRouter()
	s.router.Use(middleware.Cors)
	if s.cfg.Api.RequestLog {
		s.router.Use(chiMiddleware.Logger)
	}
	s.router.Use(chiMiddleware.Recoverer)
}

func (s *Server) initDomains() {
	healthHandler.RegisterHTTPEndPoints(s.router, healthUseCase.NewHealthUseCase(healthRepo.NewHealthRepository(s.db)))
	bookHandler.RegisterHTTPEndPoints(s.router, bookUseCase.NewBookUseCase(bookRepo.NewBookRepository(s.db)))
}

func (s *Server) Migrate() {
	source := "file://database/migrations/"
	if s.cfg.DockerTest.Driver == "postgres" {
		driver, err := postgres.WithInstance(s.db.DB, &postgres.Config{})
		if err != nil {
			log.Fatalf("error instantiating database: %v", err)
		}
		m, err := migrate.NewWithDatabaseInstance(
			source, s.cfg.DockerTest.Driver, driver,
		)
		if err != nil {
			log.Fatalf("error connecting to database: %v", err)
		}
		log.Println("migrating...")
		err = m.Up()
		if err != nil {
			if err != migrate.ErrNoChange {
				log.Panicf("error migrating: %v", err)
			}
		}

		log.Println("done migration.")
	}
}

func (s *Server) Run() error {
	s.httpServer = &http.Server{
		Addr:           ":" + s.cfg.Api.Port,
		Handler:        s.router,
		ReadTimeout:    s.cfg.Api.ReadTimeout * time.Second,
		WriteTimeout:   s.cfg.Api.WriteTimeout * time.Second,
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
		log.Printf("serving at %s:%s\n", s.cfg.Api.Host, s.cfg.Api.Port)
		printAllRegisteredRoutes(s.router)
		err := s.httpServer.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), s.cfg.Api.IdleTimeout*time.Second)
	defer shutdown()

	return s.httpServer.Shutdown(ctx)
}

func (s *Server) GetConfig() *configs.Configs {
	return s.cfg
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
