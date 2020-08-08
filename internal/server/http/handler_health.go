package http

import "net/http"

func (h *Handlers) HandleLive() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Api.HandleLive()
	}
}

func (h *Handlers) HandleReady() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Api.HandleReady()
	}
}

