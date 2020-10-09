package respond

import (
	"net/http"

	"github.com/go-chi/render"
)

func Error(w http.ResponseWriter, r *http.Request, errorCode int, responseError interface{}) {
	render.Status(r, errorCode)
	render.JSON(w, r, responseError)
}

func Success(w http.ResponseWriter, r *http.Request, successCode int, message interface{}) {
	render.Status(r, successCode)
	render.JSON(w, r, message)
}
