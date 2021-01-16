package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/volatiletech/null/v8"

	"github.com/gmhafiz/go8/internal/resource"
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
	testDeleteOneBook(t, id)
}

func testDeleteOneBook(t *E2eTest, id int64) {

}

func testAddOneBook(t *E2eTest) int64 {
	want := resource.BookResource{
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

	got := resource.BookResource{}
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

	expected, _ := json.Marshal(make([]*resource.BookResource, 0))

	if bytes.Compare(expected, got) != 0 {
		log.Printf("handler returned unexpected body: got %v want %v", string(got), expected)
	}

	log.Println("testEmptyBook passes2")
}
