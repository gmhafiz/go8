package presentation

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func GetURLParamInt64(r *http.Request, param string) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, param), 10, 64)
}
