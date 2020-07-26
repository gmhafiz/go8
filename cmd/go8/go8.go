package main

import (
	"log"

	"eight/internal/api"
	"eight/internal/configs"
	"eight/internal/platform/datastore"
	"eight/internal/server/http"
	"eight/internal/service/authors"
	"eight/internal/service/books"
)

func main() {
	cfg, err := configs.NewService("dev")
	if err != nil {
		log.Fatal(err)
	}

	dataStoreCfg, err := cfg.DataStore()
	if err != nil {
		log.Panic(err)
	}

	pqdriver, db, err := datastore.NewService(dataStoreCfg)
	if err != nil {
		log.Panic(err)
	}

	bookService, err := books.NewService(pqdriver, db)
	if err != nil {
		log.Panic(err)
	}

	authorService, err := authors.NewService(pqdriver, db)
	if err != nil {
		log.Panic(err)
	}

	// additional microservice added here
	a, err := api.NewService(bookService, authorService)
	if err != nil {
		log.Fatal(err)
	}

	httpCfg, err := cfg.HTTP()
	if err != nil {
		log.Fatal(err)
	}

	h, err := http.NewService(httpCfg, a)
	if err != nil {
		log.Panic(err)
	}

	h.Start()
}
