package param

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/gmhafiz/go8/internal/utility/respond"
)

func UInt64(w http.ResponseWriter, r *http.Request, param string) uint64 {
	val, err := strconv.ParseInt(chi.URLParam(r, param), 10, 64)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
	}

	return uint64(val)
}

func Int64(w http.ResponseWriter, r *http.Request, param string) int64 {
	val, err := strconv.ParseInt(chi.URLParam(r, param), 10, 64)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err)
	}

	return val
}

func String(r *http.Request, param string) string {
	return chi.URLParam(r, param)
}
