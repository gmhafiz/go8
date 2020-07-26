package main

import (
	"eight/internal/api"
	"eight/internal/configs"
	"eight/internal/platform/datastore"
	"eight/internal/server/http"
	"eight/internal/service/books"
	"io/ioutil"
	"log"
	nh "net/http"
	"net/http/httptest"
	"testing"
)

func TestAPI_GetAllBooks(t *testing.T) {
	resp := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/api/v1/books", nil)
	req.Header.Add("Authorization", "Bearer token")

	cfg, err := configs.NewService("test")
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

	a, err := api.NewService(bookService, nil)
	if err != nil {
		log.Fatal(err)
	}

	h := &http.Handlers{
		Api: a,
	}
	r := http.Router(h)

	r.ServeHTTP(resp, req)

	if p, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Fail()
	} else {
		t.Log(p)
	}

	if resp.Code != nh.StatusOK {
		t.Fail()
	}
}
