package respond

import (
	"encoding/json"
	"log"
	"net/http"

	"codeberg.org/gmhafiz/go8/internal/utility/message"
)

type Standard struct {
	Data any  `json:"data"`
	Meta Meta `json:"meta,omitzero"`
}

type Meta struct {
	Size  int `json:"size"`
	Total int `json:"total"`
}

func JSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if payload == nil {
		return
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		Error(w, http.StatusInternalServerError, message.ErrInternalError)
		return
	}

	if string(data) == "null" {
		_, _ = w.Write([]byte("[]"))
		return
	}

	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
		Error(w, http.StatusInternalServerError, message.ErrInternalError)
		return
	}
}
