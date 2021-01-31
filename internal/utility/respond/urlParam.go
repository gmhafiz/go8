package respond

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func GetURLParamInt64(w http.ResponseWriter, r *http.Request, param string) int64 {
	val, err := strconv.ParseInt(chi.URLParam(r, param), 10, 64)
	if err != nil {
		Error(w, http.StatusBadRequest, err)
	}

	return val
}

func GetURLParam(r *http.Request, param string) string {
	return chi.URLParam(r, param)
}
