package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"

	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/models"
	"github.com/gmhafiz/go8/internal/utility/respond"
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
	var bookRequest book.Request
	err := json.NewDecoder(r.Body).Decode(&bookRequest)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	errs := respond.Validate(h.validate, bookRequest)
	if errs != nil {
		respond.Error(w, http.StatusBadRequest, map[string][]string{"errors": errs})
		return
	}

	bk, err := h.useCase.Create(context.Background(), bookRequest)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	b, err := book.Book(bk)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}
	respond.Render(w, http.StatusCreated, b)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	bookID := respond.GetURLParamInt64(w, r, "bookID")

	var bookRequest book.Request
	err := json.NewDecoder(r.Body).Decode(&bookRequest)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}
	bookRequest.BookID = strconv.FormatInt(bookID, 10)

	errs := respond.Validate(h.validate, bookRequest)
	if errs != nil {
		respond.Error(w, http.StatusBadRequest, map[string][]string{"errors": errs})
		return
	}

	resp, err := h.useCase.Update(context.Background(), book.ToBook(&bookRequest))
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	res, err := book.Book(resp)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	respond.Render(w, http.StatusOK, res)
}

func (h *Handler) All(w http.ResponseWriter, r *http.Request) {
	filters, search := book.Filter(r.URL.Query())

	var books []*models.Book
	if search {
		resp, err := h.useCase.Search(context.Background(), filters)
		if err != nil {
			respond.Error(w, http.StatusInternalServerError, err)
			return
		}
		books = resp
	} else {
		resp, err := h.useCase.All(r.Context())
		if err != nil {
			respond.Error(w, http.StatusInternalServerError, err)
			return
		}
		books = resp
	}

	list, err := book.Books(books)
	if err != nil {
		respond.Render(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	} else if list == nil {
		respond.Error(w, http.StatusNoContent, nil)
		return
	}

	respond.Render(w, http.StatusOK, list)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	bookID := respond.GetURLParamInt64(w, r, "bookID")

	b, err := h.useCase.Find(context.Background(), bookID)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, nil)
		return
	}
	list, err := book.Book(b)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, nil)
		return
	}
	respond.Render(w, http.StatusOK, list)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	bookID := respond.GetURLParamInt64(w, r, "bookID")

	err := h.useCase.Delete(context.Background(), bookID)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, nil)
		return
	}

	respond.Render(w, http.StatusOK, nil)
}
