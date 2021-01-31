package respond

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
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

func Error(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.WriteHeader(statusCode)

	if payload == nil {
		_, err := w.Write(nil)
		if err != nil {
			log.Println(err)
		}
	} else {
		if reflect.ValueOf(payload).Type().String() == "*errors.errorString" {
			err := fmt.Errorf("%v", payload)
			p := map[string]string{
				"error": err.Error(),
			}
			data, err := json.Marshal(p)
			if err != nil {
				log.Println(err)
			}
			_, err = w.Write(data)
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
}
