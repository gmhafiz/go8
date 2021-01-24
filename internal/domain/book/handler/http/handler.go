package http

import (
	"context"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"

	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/resource"
	"github.com/gmhafiz/go8/internal/utility"
)

type Handler struct {
	useCase  book.UseCase
	validate *validator.Validate
}

func NewHandler(useCase book.UseCase) *Handler {
	return &Handler{
		useCase:  useCase,
		validate: validator.New(),
	}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var bookRequest resource.BookRequest
	err := json.NewDecoder(r.Body).Decode(&bookRequest)
	if err != nil {
		utility.Render(w, http.StatusBadRequest, nil)
		return
	}

	errs := utility.Validate(h.validate, bookRequest)
	if errs != nil {
		utility.Render(w, http.StatusBadRequest, map[string][]string{"errors": errs})
		return
	}

	bk, err := h.useCase.Create(context.Background(), bookRequest.Title, bookRequest.Description, bookRequest.ImageURL, bookRequest.PublishedDate)
	if err != nil {
		utility.Render(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.Render(w, http.StatusCreated, bk)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	bookID, err  := utility.GetURLParamInt64(r, "bookID")
	if err != nil {
		utility.Render(w, http.StatusBadRequest, nil)
		return
	}

	var bookRequest resource.BookRequest
	err = json.NewDecoder(r.Body).Decode(&bookRequest)
	if err != nil {
		utility.Render(w, http.StatusBadRequest, nil)
		return
	}
	bookRequest.BookID = strconv.FormatInt(bookID, 10)

	errs := utility.Validate(h.validate, bookRequest)
	if errs != nil {
		utility.Render(w, http.StatusBadRequest, map[string][]string{"errors": errs})
		return
	}

	resp, err := h.useCase.Update(context.Background(), resource.ToBook(&bookRequest))
	if err != nil {
		utility.Render(w, http.StatusInternalServerError, err.Error())
		return
	}

	res, err := resource.Book(resp)
	if err != nil {
		utility.Render(w, http.StatusInternalServerError, err.Error())
		return
	}

	utility.Render(w, http.StatusOK, res)
}

func (h *Handler) All(w http.ResponseWriter, r *http.Request) {
	resp, err := h.useCase.All(r.Context())
	if err != nil {
		utility.Render(w, http.StatusInternalServerError, err.Error())
		return
	}

	list, err := resource.Books(resp)
	if err != nil {
		utility.Render(w, http.StatusInternalServerError, err.Error())
		return
	} else if list == nil {
		utility.Render(w, http.StatusNoContent, nil)
		return
	}

	utility.Render(w, http.StatusOK, list)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	bookID, err  := utility.GetURLParamInt64(r, "bookID")
	if err != nil {
		utility.Render(w, http.StatusInternalServerError, nil)
		return
	}

	err = h.useCase.Delete(context.Background(), bookID)
	if err != nil {
		utility.Render(w, http.StatusInternalServerError, nil)
		return
	}

	utility.Render(w, http.StatusOK, nil)
}
