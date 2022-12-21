package respond

import (
	"encoding/json"
	"log"
	"net/http"
)

func Errors(w http.ResponseWriter, statusCode int, errors []string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(statusCode)

	if errors == nil {
		write(w, nil)
		return
	}

	p := map[string][]string{
		"message": errors,
	}
	data, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
	}

	if string(data) == "null" {
		return
	}

	write(w, data)
}

func Error(w http.ResponseWriter, statusCode int, message error) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(statusCode)

	var p map[string]string
	if message == nil {
		write(w, nil)
		return
	}

	p = map[string]string{
		"message": message.Error(),
	}
	data, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
	}

	if string(data) == "null" {
		return
	}

	write(w, data)
}

func write(w http.ResponseWriter, data []byte) {
	_, err := w.Write(data)
	if err != nil {
		log.Println(err)
	}
}
