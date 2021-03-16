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
			write(w, data)
		} else if reflect.TypeOf(payload).Kind() == reflect.Slice {
			var errors []string
			s := reflect.ValueOf(payload)
			for i := 0; i < s.Len(); i++ {
				errors = append(errors, s.Index(i).String())
			}
			p := map[string][]string{
				"error": errors,
			}
			data, err := json.Marshal(p)
			if err != nil {
				log.Println(err)
			}
			write(w, data)
		} else {
			data, err := json.Marshal(payload)
			if err != nil {
				log.Println(err)
			}
			write(w, data)
		}
	}
}

func write(w http.ResponseWriter, data []byte) {
	_, err := w.Write(data)
	if err != nil {
		log.Println(err)
	}
}
