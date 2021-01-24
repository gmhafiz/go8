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
	p "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"

	"github.com/gmhafiz/go8/configs"
	"github.com/gmhafiz/go8/internal/domain/book"
	bookHTTP "github.com/gmhafiz/go8/internal/domain/book/handler/http"
	bookRepo "github.com/gmhafiz/go8/internal/domain/book/repository/postgres"
	bookUseCase "github.com/gmhafiz/go8/internal/domain/book/usecase"
	"github.com/gmhafiz/go8/internal/domain/health"
	healthHTTP "github.com/gmhafiz/go8/internal/domain/health/handler/http"
	"github.com/gmhafiz/go8/internal/domain/health/repository/postgres"
	"github.com/gmhafiz/go8/internal/domain/health/usecase"
	"github.com/gmhafiz/go8/internal/middleware"
	"github.com/gmhafiz/go8/third_party/database"
)

type Server struct {
	httpServer *http.Server
	db         *sqlx.DB
	cfg        *configs.Configs
	Domain
}

type Domain struct {
	BookUC   book.UseCase
	HealthUC health.UseCase
}

func New(version string) *Server {
	log.Printf("staring API version: %s\n", version)
	return &Server{}
}

func (s *Server) Init() {
	s.NewConfig()
	s.NewDatabase()
	s.InitDomains()
}

func (s *Server) NewConfig() {
	s.cfg = configs.New()
}

func (s *Server) NewDatabase() {
	s.db = database.NewSqlx(s.cfg)
}

func (s *Server) InitDomains() {
	s.BookUC = bookUseCase.NewBookUseCase(bookRepo.NewBookRepository(s.db))
	s.HealthUC = usecase.NewHealthUseCase(postgres.NewHealthRepository(s.db))
}

func (s *Server) Migrate() {
	source := "file://database/migrations/"
	if s.cfg.DockerTest.Driver == "postgres" {
		driver, err := p.WithInstance(s.db.DB, &p.Config{})
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
	router := chi.NewRouter()
	router.Use(middleware.Cors)
	if s.cfg.Api.RequestLog == true {
		router.Use(chiMiddleware.Logger)
	}
	router.Use(chiMiddleware.Recoverer)

	healthHTTP.RegisterHTTPEndPoints(router, s.HealthUC)
	bookHTTP.RegisterHTTPEndPoints(router, s.BookUC)

	s.httpServer = &http.Server{
		Addr:           ":" + s.cfg.Api.Port,
		Handler:        router,
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
		printAllRegisteredRoutes(router)
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

func (s *Server) GetHTTPServer() *http.Server {
	return s.httpServer
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
