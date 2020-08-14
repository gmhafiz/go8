package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	nh "net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"

	"eight/internal/api"
	"eight/internal/configs"
	"eight/internal/datastore"
	"eight/internal/domain/books"
	"eight/internal/models"
	"eight/internal/server/http"
)

//func (suite *TestSuite) setupTest() (*chi.Mux, *http.Handlers, *sql.DB) {
func setupTest() (*chi.Mux, *http.Handlers, *sql.DB) {
	cfg, err := configs.NewService("test")
	if err != nil {
		log.Fatal(err)
	}

	dataStoreCfg, err := cfg.DataStore()
	if err != nil {
		log.Panic(err)
	}
	db, err := datastore.NewService(dataStoreCfg)
	if err != nil {
		log.Panic(err)
	}
	bookService, err := books.NewService(db, nil)
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
	r := http.Router(h, nil)

	return r, h, db
}

// tearDown Truncates all tables in the database
func tearDown() {
	_, _, db := setupTest()

	v := reflect.ValueOf(models.TableNames)

	tableName := make([]string, v.NumField())

	var tableNames string
	for i := 0; i < v.NumField(); i++ {
		tableName[i] = v.Field(i).String()
		log.Println(v.Field(i).Interface())
		tableNames += tableName[i] + ","
	}

	length := len(tableNames)
	if length > 0 && tableNames[length-1] == ',' {
		tableNames = tableNames[:length-1]
	}

	query := fmt.Sprintf("TRUNCATE TABLE %s;", tableNames)
	_, err := db.Exec(query)
	if err != nil {
		log.Println(err)
	}
}

func TestAPI_GetAllBooks(t *testing.T) {
	assert.True(t, true, "True is true!")

	r, _, _ := setupTest()

	resp := httptest.NewRecorder()

	req := httptest.NewRequest("GET", "/api/v1/books", nil)
	req.Header.Add("Authorization", "Bearer token")

	r.ServeHTTP(resp, req)

	if _, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Fail()
	}

	if resp.Code != nh.StatusOK {
		t.Fail()
	}

	if resp.Body.Len() != 0 {
		t.Fail()
	}

}

func TestAPI_CreateBook(t *testing.T) {
	r, _, _ := setupTest()

	resp := httptest.NewRecorder()

	type bookRequest struct {
		Title         string `json:"title"`
		PublishedDate string `json:"published_date"`
		ImageURL      string `json:"image_url"`
		Description   string `json:"description"`
	}

	bookR := bookRequest{
		Title:         "Test Title",
		PublishedDate: "2020-07-31 15:04:05.123499999",
		ImageURL:      "https://example.com/image.png",
		Description:   "Test Description",
	}

	payload, err := json.Marshal(bookR)
	if err != nil {
		log.Fatal(err)
	}

	req := httptest.NewRequest("POST", "/api/v1/book", bytes.NewReader(payload))
	req.Header.Add("Authorization", "Bearer token")

	r.ServeHTTP(resp, req)

	var bookResponse models.Book
	if p, err := ioutil.ReadAll(resp.Body); err != nil {
		t.Fail()
	} else {
		log.Println(p)
		log.Println(resp.Code)

		err = json.Unmarshal(p, &bookResponse)
		if err != nil {
			log.Println(err)
		}
	}

	if resp.Code != nh.StatusCreated {
		t.Fail()
	}

	if bookR.Title != bookResponse.Title {
		t.Fail()
	}

	assert.Equal(t, bookR, bookResponse)

	tearDown()
}
