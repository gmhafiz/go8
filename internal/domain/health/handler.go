package health

import (
	"database/sql"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
	"net/http"
)

type Handler struct {
	Router *chi.Mux
	Log    zerolog.Logger
	DB     *sql.DB
}

func New(router *chi.Mux, log zerolog.Logger, db *sql.DB) {
	handler := &Handler{
		Router: router,
		Log: log,
		DB: db,
	}

	initRoutes(router, handler)
}

func (h *Handler) Liveness() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.Status(r, http.StatusOK)
		return
	}
}

func (h *Handler) Readiness() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h.DB.Ping()
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			return
		}
		render.Status(r, http.StatusOK)
		return
	}
}