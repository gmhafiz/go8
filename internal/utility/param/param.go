package param

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func UInt64(r *http.Request, param string) (uint64, error) {
	val, err := strconv.ParseInt(chi.URLParam(r, param), 10, 64)
	if err != nil {
		return 0, err
	}

	return uint64(val), nil
}

func Int(r *http.Request, param string) (int, error) {
	val, err := strconv.Atoi(chi.URLParam(r, param))
	if err != nil {
		return 0, err
	}

	return val, nil
}

func String(r *http.Request, param string) string {
	return chi.URLParam(r, chi.URLParam(r, param))
}
