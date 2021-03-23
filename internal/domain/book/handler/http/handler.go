package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/friendsofgo/errors"
	"github.com/go-playground/validator/v10"

	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/models"
	"github.com/gmhafiz/go8/internal/utility/message"
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

// Create creates a new book record
// @Summary Create a Book
// @Description Get a book with JSON payload
// @Accept json
// @Produce json
// @Param Book body book.Request true "Book Request"
// @Success 201 {object} book.Res
// @Failure 400 {string} Bad Request
// @Failure 500 {string} Internal Server Error
// @Router /api/v1/books [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var bookRequest book.Request
	err := book.Decode(r.Body, &bookRequest)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	errs := respond.Validate(h.validate, bookRequest)
	if errs != nil {
		respond.Error(w, http.StatusBadRequest, errs)
		return
	}

	bk, err := h.useCase.Create(context.Background(), book.ToBook(&bookRequest))
	if err != nil {
		if err == sql.ErrNoRows {
			respond.Error(w, http.StatusBadRequest, message.ErrBadRequest)
			return
		}
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	b, err := book.Resource(bk)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}
	respond.Render(w, http.StatusCreated, b)
}

// Get a book by its ID
// @Summary Get a Book
// @Description Get a book by its id.
// @Accept json
// @Produce json
// @Param bookID path int true "book ID"
// @Success 200 {object} book.Res
// @Failure 400 {string} Bad Request
// @Failure 500 {string} Internal Server Error
// @Router /api/v1/books/{bookID} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	bookID := respond.GetURLParamInt64(w, r, "bookID")

	b, err := h.useCase.Read(context.Background(), bookID)
	if err != nil {
		if err == sql.ErrNoRows {
			respond.Error(w, http.StatusBadRequest, sql.ErrNoRows)
			return
		}
		respond.Error(w, http.StatusInternalServerError, nil)
		return
	}
	list, err := book.Resource(b)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, nil)
		return
	}
	respond.Render(w, http.StatusOK, list)
}

// List will fetch the article based on given params
// @Summary Show all books
// @Description Get all books. By default it gets first page with 10 items.
// @Accept json
// @Produce json
// @Param page query string false "page"
// @Param size query string false "size"
// @Param title query string false "term"
// @Param description query string false "term"
// @Success 200 {object} []book.Res
// @Failure 500 {string} Internal Server Error
// @Router /api/v1/books [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	filters := book.Filters(r.URL.Query())

	var books []*models.Book
	ctx := context.Background()

	switch filters.Base.Search {
	case true:
		resp, err := h.useCase.Search(ctx, filters)
		if err != nil {
			if errors.Is(err, message.ErrFetchingBook) {
				respond.Error(w, http.StatusInternalServerError, err)
				return
			}
			respond.Error(w, http.StatusInternalServerError, err)
			return
		}
		books = resp
	default:
		resp, err := h.useCase.List(ctx, filters)
		if err != nil {
			if errors.Is(err, message.ErrFetchingBook) {
				respond.Error(w, http.StatusInternalServerError, err)
				return
			}
			respond.Error(w, http.StatusInternalServerError, err)
			return
		}
		books = resp
	}

	list, err := book.Resources(books)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, message.ErrFormingResponse)
		return
	} else if list == nil {
		respond.Error(w, http.StatusNoContent, nil)
		return
	}

	respond.Render(w, http.StatusOK, list)
}

// Update a book
// @Summary Update a Book
// @Description Update a book by its model.
// @Accept json
// @Produce json
// @Param Book body book.Request true "Book Request"
// @Success 200 {object} []book.Res
// @Failure 400 {string} Bad Request
// @Failure 500 {string} Internal Server Error
// @Router /api/v1/books/{bookID} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	bookID := respond.GetURLParamInt64(w, r, "bookID")

	var bookRequest book.Request
	err := json.NewDecoder(r.Body).Decode(&bookRequest)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}
	bookRequest.BookID = bookID

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

	res, err := book.Resource(resp)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, err)
		return
	}

	respond.Render(w, http.StatusOK, res)
}

// Delete a book by its ID
// @Summary Delete a Book
// @Description Delete a book by its id.
// @Accept json
// @Produce json
// @Param id path int true "book ID"
// @Success 200 "Ok"
// @Failure 500 {string} Internal Server Error
// @Router /api/v1/books/{bookID} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	bookID := respond.GetURLParamInt64(w, r, "bookID")

	err := h.useCase.Delete(context.Background(), bookID)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, message.ErrInternalError)
		return
	}

	respond.Render(w, http.StatusOK, nil)
}
