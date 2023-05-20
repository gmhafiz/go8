package authentication

import (
	"github.com/gmhafiz/scs/v2"
	"github.com/go-chi/chi/v5"

	"github.com/gmhafiz/go8/internal/middleware"
)

func RegisterHTTPEndPoints(router *chi.Mux, session *scs.SessionManager, repo Repo) {
	h := NewHandler(session, repo)

	router.Post("/api/v1/login", h.Login)
	router.Post("/api/v1/register", h.Register)

	router.Route("/api/v1/logout", func(router chi.Router) {
		router.Post("/", h.Logout)
	})

	router.Route("/api/v1/restricted", func(router chi.Router) {
		router.Use(middleware.Authenticate(session))
		router.Get("/csrf", h.Csrf)
		router.Get("/", h.Protected)
		router.Get("/me", h.Me)
		router.Post("/logout/{userID}", h.ForceLogout)
	})
}
