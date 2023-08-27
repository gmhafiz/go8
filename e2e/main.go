package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gmhafiz/go8/config"
	"github.com/gmhafiz/go8/internal/domain/book"
)

// Version is injected using ldflags during build time
const Version = "v0.1.0"

var url = ""

func main() {
	log.Printf("Starting e2e API version: %s\n", Version)
	cfg := config.New()

	url = fmt.Sprintf("http://%s:%s", cfg.Api.Host, cfg.Api.Port)

	waitForApi(fmt.Sprintf("%s/api/health/readiness", url))

	run()
}

func run() {
	testBook()

	log.Println("all tests have passed.")
}

func testBook() {
	testEmptyBook()
	id := testAddOneBook()
	id = testGetOneBook(id)
	testUpdateBook(id)
	testDeleteOneBook(id)
}

func testEmptyBook() {
	resp, err := http.Get(fmt.Sprintf("%s/api/v1/book", url))
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	got, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if status := resp.StatusCode; status != http.StatusOK {
		log.Printf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected, _ := json.Marshal(make([]*book.Res, 0))

	if !bytes.Equal(expected, got) {
		log.Printf("handler returned unexpected body: got %v want %v", string(got), expected)
	}

	log.Println("testEmptyBook passes")
}

func testAddOneBook() int {
	want := &book.CreateRequest{
		Title:         "test01",
		PublishedDate: "2020-02-02",
		ImageURL:      "https://example.com/image.png",
		Description:   "test01",
	}

	bR, _ := json.Marshal(want)

	resp, err := http.Post(
		fmt.Sprintf("%s/api/v1/book", url),
		"Content-Type: application/json",
		bytes.NewBuffer(bR),
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	gotBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	got := book.Res{}
	err = json.Unmarshal(gotBody, &got)
	if err != nil {
		log.Println(err)
	}

	if resp.StatusCode != http.StatusCreated {
		log.Printf("error code want %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	if want.Title != got.Title && want.Description != got.Description && want.
		ImageURL != got.ImageURL && want.PublishedDate != got.PublishedDate.String() {
		log.Printf("want %v, got %v\n", want, got)
	}

	log.Println("testAddOneBook passes")
	return got.ID
}

func testGetOneBook(id int) int {
	client := &http.Client{}

	url := fmt.Sprintf("%s/api/v1/book/%d", url, id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("error code fail, want %d, got %d\n", http.StatusOK, resp.StatusCode)
	}

	got := book.Res{}
	err = json.Unmarshal(respBody, &got)
	if err != nil {
		log.Println(err)
	}

	log.Println("testGetBook passes")

	return got.ID
}

func testUpdateBook(bookID int) {
	newBook := book.CreateRequest{
		Title:         "updated title",
		PublishedDate: "2020-07-31T15:04:05.123499999Z",
		ImageURL:      "https://example.com/image.png",
		Description:   "test description",
	}

	client := &http.Client{}

	bR, err := json.Marshal(&newBook)
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("%s/api/v1/book/%d", url, bookID)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(bR))
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("error code fail, want %d, got %d\n", http.StatusOK, resp.StatusCode)
	}

	got := book.Res{}
	err = json.Unmarshal(respBody, &got)
	if err != nil {
		log.Println(err)
	}

	if got.ID != bookID && got.Title != newBook.Title && got.Description != newBook.Description && got.ImageURL != newBook.ImageURL {
		if err != nil {
			log.Fatalf("returned resource does not match. want %v, got %v", respBody, got)
		}
	}

	log.Println("testUpdateBook passes")
}

func testDeleteOneBook(id int) {
	client := &http.Client{}

	req, err := http.NewRequest(
		http.MethodDelete, fmt.Sprintf("%s/api/v1/book/%d", url, id),
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

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("error code fail, want %d, got %d\n", http.StatusOK, resp.StatusCode)
	}
	log.Println("testDeleteOneBook passes")
}

func waitForApi(readinessURL string) {
	log.Println("Connecting to api with exponential backoff... ")
	for {
		//nolint:gosec
		_, err := http.Get(readinessURL)
		if err == nil {
			log.Println("api is up")
			return
		}

		base, capacity := time.Second, time.Minute
		for backoff := base; err != nil; backoff <<= 1 {
			if backoff > capacity {
				backoff = capacity
			}

			// A pseudo-random number generator here is fine. No need to be
			// cryptographically secure. Ignore with the following comment:
			/* #nosec */
			jitter := rand.Int63n(int64(backoff * 3))
			sleep := base + time.Duration(jitter)
			time.Sleep(sleep)
			//nolint:gosec
			_, err := http.Get(readinessURL)
			if err == nil {
				log.Println("api is up")
				return
			}
		}
	}
}
