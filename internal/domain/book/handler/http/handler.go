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
	"github.com/gmhafiz/go8/internal/utility/param"
	"github.com/gmhafiz/go8/internal/utility/respond"
	"github.com/gmhafiz/go8/internal/utility/validate"
)

type Handler struct {
	useCase  book.UseCase
	validate *validator.Validate
}

func NewHandler(useCase book.UseCase, validate *validator.Validate) *Handler {
	return &Handler{
		useCase:  useCase,
		validate: validate,
	}
}

// Create creates a new book record
// @Summary Create a Book
// @Description Create a book using JSON payload
// @Accept json
// @Produce json
// @Param Book body book.Request true "Create a book using the following format"
// @Success 201 {object} book.Res
// @Failure 400 {string} Bad Request
// @Failure 500 {string} Internal Server Error
// @router /api/v1/book [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var bookRequest book.Request
	err := book.Bind(r.Body, &bookRequest)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}

	errs := validate.Validate(h.validate, bookRequest)
	if errs != nil {
		respond.Errors(w, http.StatusBadRequest, errs)
		return
	}

	bk, err := h.useCase.Create(r.Context(), book.ToBook(&bookRequest))
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
// @router /api/v1/book/{bookID} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	bookID := param.Int64(w, r, "bookID")

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
// @Summary Shows all books
// @Description Lists all books. By default, it gets first page with 30 items.
// @Accept json
// @Produce json
// @Param page query string false "page number"
// @Param size query string false "size of result"
// @Param title query string false "search by title"
// @Param description query string false "search by description"
// @Success 200 {object} []book.Res
// @Failure 500 {string} Internal Server Error
// @router /api/v1/book [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	filters := book.Filters(r.URL.Query())

	var books []*models.Book
	ctx := r.Context()

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
// @router /api/v1/book/{bookID} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	bookID := param.Int64(w, r, "bookID")

	var bookRequest book.Request
	err := json.NewDecoder(r.Body).Decode(&bookRequest)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, nil)
		return
	}
	bookRequest.BookID = bookID

	errs := validate.Validate(h.validate, bookRequest)
	if errs != nil {
		respond.Errors(w, http.StatusBadRequest, errs)
		return
	}

	resp, err := h.useCase.Update(r.Context(), book.ToBook(&bookRequest))
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
// @router /api/v1/book/{bookID} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	bookID := param.Int64(w, r, "bookID")

	err := h.useCase.Delete(r.Context(), bookID)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, message.ErrInternalError)
		return
	}

	respond.Render(w, http.StatusOK, nil)
}
