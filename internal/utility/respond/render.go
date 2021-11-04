package respond

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

var (
	ErrInternalServerError = errors.New("internal server error")
)

type Standard struct {
	Data interface{} `json:"data"`
	Meta Meta        `json:"meta,omitempty"`
}

type Meta struct {
	Size  int   `json:"size"`
	Total int64 `json:"total"`
}

func Json(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.WriteHeader(statusCode)

	if payload == nil {
		return
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		Error(w, http.StatusInternalServerError, ErrInternalServerError)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
		Error(w, http.StatusInternalServerError, ErrInternalServerError)
		return
	}
}

func Status(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}

func Render(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.WriteHeader(statusCode)

	if payload == nil {
		_, err := w.Write(nil)
		if err != nil {
			log.Println(err)
		}
	} else {
		data, err := json.Marshal(payload)
		if err != nil {
			log.Println(err)
		}
		_, err = w.Write(data)
		if err != nil {
			log.Println(err)
		}
	}
}
