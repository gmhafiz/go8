package app

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (s *Server) handleLive() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("."))
	}
}

func (s *Server) handleReady() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("."))
	}
}

func (s *Server) getAllBooks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("books"))
	}
}

func (s *Server) getBook() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("book"))
	}
}

func (s *Server) getAllContact() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contacts := map[string]string{
			"message": "all ok",
		}
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(contacts)
	}
}

func (s *Server) handleSomething() http.HandlerFunc {
	//thing := prepareThing()
	return func(w http.ResponseWriter, r *http.Request) {
		payload := Response{
			Message: "something",
		}
		Respond(w, &payload)
	}
}

func Respond(w http.ResponseWriter, r *Response) {
	_ = json.NewEncoder(w).Encode(r)
}

func (s *Server) handleAPI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("api"))
	}
}

func (s *Server) handleAbout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("about"))
	}
}

func (s *Server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("index"))
	}
}

func (s *Server) handleAdminIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("admin index"))
	}
}
