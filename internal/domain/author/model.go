package author

import (
	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/models"
)

type WithBooks struct {
	*models.Author
	Books []book.Res `json:"books"`
}
