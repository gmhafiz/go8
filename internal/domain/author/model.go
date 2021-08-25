package author

import "github.com/gmhafiz/go8/internal/models"

type AuthorB struct {
	*models.Author
	Books models.BookSlice `json:"books"`
}
