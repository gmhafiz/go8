package authors

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	errorsUtil "go8ddd/internal/utility/errors_handling"
	"net/http"
)

type AuthorHandler struct {
	AuthorUseCase AuthorUseCase
	Validator     *validator.Validate
	Router        *chi.Mux
}

func NewHandler(router *chi.Mux, validate *validator.Validate, useCase AuthorUseCase) {
	handler := &AuthorHandler{
		AuthorUseCase: useCase,
		Validator:     validate,
		Router:        router,
	}
	initRoutes(router, handler)
}

// FetchArticle will fetch the article based on given params
func (handler *AuthorHandler) All() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list, err := handler.AuthorUseCase.All(r.Context())
		if err != nil {
			render.Status(r, errorsUtil.GetStatusCode(err))
			render.JSON(w, r, errorsUtil.ResponseError{Error: err.Error()})
			return
		}
		render.JSON(w, r, list)
	}
}
