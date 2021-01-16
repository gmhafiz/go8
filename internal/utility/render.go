package utility

import (
	"encoding/json"
	"log"
	"net/http"
)

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
