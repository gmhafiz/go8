package books

import (
	"encoding/json"
	"fmt"
	"go8ddd/internal/utility/respond"
	"net/http"
	"reflect"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"github.com/volatiletech/null/v8"

	"go8ddd/internal/model"
	"go8ddd/internal/utility/converter"
	"go8ddd/internal/utility/errors_handling"
	"go8ddd/internal/utility/success"
)

type BookHandler struct {
	BookUseCase BookUseCase
	Validator   *validator.Validate
	Router      *chi.Mux
}

type bookRequest struct {
	Title         string      `json:"title" validate:"required"`
	PublishedDate string      `json:"published_date" validate:"required"`
	ImageURL      null.String `json:"image_url" validate:"url"`
	Description   null.String `json:"description" validate:"required"`
}

func NewHandler(router *chi.Mux, validate *validator.Validate, bu BookUseCase) {
	handler := &BookHandler{
		BookUseCase: bu,
		Validator:   validate,
		Router:      router,
	}

	initRoutes(router, handler)
}

// All will fetch the article based on given params
// @Summary Show all books
// @Description Get all books. By default it gets first page with 10 items.
// @Accept json
// @Produce json
// @Param page query string false "page number"
// @Param size query string false "size"
// @Success 200 {object} []models.Book
// @Router /books [get]
func (handler *BookHandler) All() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := handler.BookUseCase.All(r.Context())
		if err != nil {
			render.Status(r, errors_handling.GetStatusCode(err))
			render.JSON(w, r, errors_handling.ResponseError{Error: err.Error()})
			return
		}
		list, err := resourcesReflection(resp)
		if err != nil {
			render.Status(r, errors_handling.GetStatusCode(err))
			render.JSON(w, r, errors_handling.ResponseError{Error: err.Error()})
			return
		}
		render.JSON(w, r, list)
	}
}

// Create creates a new bookr ecord
// @Summary Create a Book
// @Description Get a book with JSON payload
// @Accept json
// @Produce json
// @Param Book body bookRequest true "Book Request"
// @Success 201 {object} models.Book
// @Router /book [post]
func (handler *BookHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var bookRequest bookRequest
		err := json.NewDecoder(r.Body).Decode(&bookRequest)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, errors_handling.ResponseError{Error: errors_handling.ErrBadRequest.Error()})
			return
		}

		err = handler.Validator.Struct(bookRequest)
		if err != nil {
			var errs []string
			for _, err := range err.(validator.ValidationErrors) {
				errs = append(errs, fmt.Sprintf("%s is %s", err.StructNamespace(), err.Tag()))
			}
			respond.Error(w, r, http.StatusBadRequest, errors_handling.ResponseErrors{Errors: errs})
			//render.Status(r, http.StatusBadRequest)
			//render.JSON(w, r, errors_handling.ResponseErrors{Errors: errs})
			return
		}

		convertedTime, err := converter.StringToTime(bookRequest.PublishedDate)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{
				"error": err.Error(),
			})
			return
		}

		book, err := handler.BookUseCase.Create(r.Context(), &model.Book{
			Title:         bookRequest.Title,
			PublishedDate: convertedTime,
			ImageURL:      bookRequest.ImageURL,
			Description:   bookRequest.Description,
		})

		if err != nil {
			respond.Error(w, r, http.StatusInternalServerError, errors_handling.ResponseError{Error: errors_handling.ErrInternalServerError.Error()})

			//render.Status(r, http.StatusInternalServerError)
			//render.JSON(w, r, errors_handling.ResponseError{Error: errors_handling.ErrInternalServerError.Error()})
			//return
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, book)
	}
}

// Get a book by its ID
// @Summary Get a Book
// @Description Get a book by its id.
// @Accept json
// @Produce json
// @Param id path int true "book ID"
// @Success 200 {object} models.Book
// @Router /book/{bookID} [get]
func (handler *BookHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value("id").(int64)

		book, err := handler.BookUseCase.Get(r.Context(), id)
		if err != nil {
			render.Status(r, errors_handling.GetStatusCode(err))
			render.JSON(w, r, errors_handling.ResponseError{Error: err.Error()})
			return
		}

		resp, err := resourceReflection(book)
		if err != nil {
			render.Status(r, errors_handling.GetStatusCode(err))
			render.JSON(w, r, errors_handling.ResponseError{Error: err.Error()})
			return
		}

		render.JSON(w, r, resp)
	}
}

// Update a book
// @Summary Update a Book
// @Description Update a book by its model.
// @Accept json
// @Produce json
// @Param Book body bookRequest true "Book Request"
// @Success 200 "Ok"
// @Failure 500 "Internal Server error"
// @Router /book/{bookID} [put]
func (handler *BookHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var bookRequest *model.Book
		err := json.NewDecoder(r.Body).Decode(&bookRequest)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, errors_handling.ResponseError{Error: errors_handling.ErrBadRequest.Error()})
			return
		}

		book, err := handler.BookUseCase.Update(r.Context(), bookRequest)
		if err != nil {
			render.Status(r, errors_handling.GetStatusCode(err))
			render.JSON(w, r, errors_handling.ResponseError{Error: err.Error()})
			return
		}
		render.JSON(w, r, book)
	}
}

// Delete a book by its ID
// @Summary Delete a Book
// @Description Delete a book by its id.
// @Accept json
// @Produce json
// @Param id path int true "book ID"
// @Success 200 "Ok"
// @Failure 500 "Internal Server error"
// @Router /book/{bookID} [delete]
func (handler *BookHandler) Delete() http.HandlerFunc {
	type deleteType struct {
		HardDelete bool `json:"hard_delete,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value("id").(int64)

		var deleteType deleteType
		err := json.NewDecoder(r.Body).Decode(&deleteType)
		if err != nil {
			deleteType.HardDelete = false
		}

		err = handler.BookUseCase.Delete(r.Context(), id, deleteType.HardDelete)
		if err != nil {
			render.Status(r, errors_handling.GetStatusCode(err))
			render.JSON(w, r, errors_handling.ResponseError{Error: err.Error()})
			return
		}
		render.JSON(w, r, success.Accepted)
	}
}

type bookResource struct {
	BookID        int64       `json:"book_id" deepcopier:"field:book_id"`
	Title         string      `json:"title" deepcopier:"field:title"`
	PublishedDate time.Time   `json:"published_date" deepcopier:"field:force"`
	ImageURL      null.String `json:"image_url" deepcopier:"field:image_url"`
	Description   null.String `json:"description" deepcopier:"field:description"`
}

func resourceReflection(book *model.Book) (bookResource, error) {
	var resource bookResource

	err := copier.Copy(&resource, &book)
	if err != nil {
		return resource, err
	}

	return resource, nil
}

func resourcesReflection(books []*model.Book) (interface{}, error) {
	var resource bookResource

	rt := reflect.TypeOf(books)
	if rt.Kind() == reflect.Slice {
		var resources []bookResource
		for _, book := range books {
			res, _ := resourceReflection(book)
			resources = append(resources, res)
		}
		return resources, nil
	}

	err := copier.Copy(&resource, books)
	if err != nil {
		return resource, err
	}

	return resource, nil
}
