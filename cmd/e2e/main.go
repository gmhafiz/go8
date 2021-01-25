package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/volatiletech/null/v8"

	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/server"
)

const Version = "v0.4.0-test"

func main() {
	s := server.New(Version)
	s.Init()
	s.Migrate()

	t := NewE2eTest(s)
	t.Run()
}

type E2eTest struct {
	server *server.Server
}

func NewE2eTest(server *server.Server) *E2eTest {
	return &E2eTest{
		server: server,
	}
}

func (t *E2eTest) Run() {
	testEmptyBook(t)
	id := testAddOneBook(t)
	testUpdateBook(t, strconv.FormatInt(id, 10))
	testDeleteOneBook(t, id)

	log.Println("all tests passed.")
}

func testUpdateBook(t *E2eTest, bookID string) {
	newBook := book.Request{
		BookID:        bookID,
		Title:         "updated title",
		PublishedDate: "2020-07-31T15:04:05.123499999Z",
		ImageURL:      "https://example.com/image.png",
		Description:   "test description",
	}

	client := &http.Client{}

	bR, _ := json.Marshal(&newBook)

	req, err := http.NewRequest(
		http.MethodPut, fmt.Sprintf("http://localhost:%s/api/v1/books/%s",
			t.server.GetConfig().Api.Port, newBook.BookID), bytes.NewBuffer(bR),
	)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("error code fail, want %d, got %d\n", http.StatusOK, resp.StatusCode)
	}

	got := book.Resource{}
	err = json.Unmarshal(respBody, &got)
	if err != nil {
		log.Println(err)
	}

	if strconv.FormatInt(got.BookID, 10) != newBook.BookID && got.Title != newBook.Title && got.Description.String != newBook.Description && got.ImageURL.String != newBook.ImageURL {
		if err != nil {
			log.Fatalf("returned resource does not match. want %v, got %v", respBody, got)
		}
	}

	log.Println("testUpdateBook passes")
}

func testDeleteOneBook(t *E2eTest, id int64) {
	client := &http.Client{}

	req, err := http.NewRequest(
		http.MethodDelete, fmt.Sprintf("http://localhost:%s/api/v1/books/%d", t.server.GetConfig().Api.Port, id),
		nil,
	)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("error code fail, want %d, got %d\n", http.StatusOK, resp.StatusCode)
	}
	log.Println("testDeleteOneBook passes")
}

func testAddOneBook(t *E2eTest) int64 {
	want := book.Resource{
		BookID:        1,
		Title:         "test title",
		PublishedDate: time.Time{},
		ImageURL: null.String{
			String: "https://example.com/image.png",
			Valid:  true,
		},
		Description: null.String{
			String: "test description",
			Valid:  true,
		},
	}

	bR, _ := json.Marshal(&want)

	resp, err := http.Post(
		fmt.Sprintf("http://localhost:%s/api/v1/books",
			t.server.GetConfig().Api.Port),
		"Content-Type: application/json",
		bytes.NewBuffer(bR),
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	gotBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	got := book.Resource{}
	err = json.Unmarshal(gotBody, &got)
	if err != nil {
		log.Println(err)
	}

	if resp.StatusCode != http.StatusCreated {
		log.Printf("error code want %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	if want.Title != got.Title && want.Description.String != got.Description.String && want.
		ImageURL.String != got.ImageURL.String && !want.PublishedDate.Equal(got.PublishedDate) {
		log.Printf("want %v, got %v\n", want, got)
	}

	log.Println("testAddOneBook passes")
	return got.BookID
}

func testEmptyBook(t *E2eTest) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%s/api/v1/books",
		t.server.GetConfig().Api.Port))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	got, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if status := resp.StatusCode; status != http.StatusOK {
		log.Printf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected, _ := json.Marshal(make([]*book.Resource, 0))

	if !bytes.Equal(expected, got) {
		log.Printf("handler returned unexpected body: got %v want %v", string(got), expected)
	}

	log.Println("testEmptyBook passes")
}
