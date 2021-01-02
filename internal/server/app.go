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

type App struct {
	httpServer *http.Server
	bookUC     book.UseCase
	healthUC   health.UseCase
}

func NewApp(cfg *configs.Configs) *App {
	db := database.NewSqlx(cfg)

	return &App{
		bookUC:   bookUseCase.NewBookUseCase(bookRepo.NewBookRepository(db)),
		healthUC: usecase.NewHealthUseCase(postgres.NewHealthRepository(db)),
	}
}

func (a *App) Run(cfg *configs.Configs, version string) error {
	router := chi.NewRouter()
	router.Use(middleware.Cors)
	router.Use(chiMiddleware.Logger)
	router.Use(chiMiddleware.Recoverer)

	healthHTTP.RegisterHTTPEndPoints(router, a.healthUC)
	bookHTTP.RegisterHTTPEndPoints(router, a.bookUC)

	a.httpServer = &http.Server{
		Addr:           ":" + cfg.Api.Port,
		Handler:        router,
		ReadTimeout:    cfg.Api.ReadTimeout * time.Second,
		WriteTimeout:   cfg.Api.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
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
		log.Printf("API version: %s\n", version)
		log.Printf("serving at %s:%s\n", cfg.Api.Host, cfg.Api.Port)
		printAllRegisteredRoutes(router)
		err := a.httpServer.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), cfg.Api.IdleTimeout*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
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
