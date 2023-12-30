package repository

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"

	"github.com/gmhafiz/go8/ent/gen"
	entAuthor "github.com/gmhafiz/go8/ent/gen/author"
	"github.com/gmhafiz/go8/ent/gen/predicate"
	"github.com/gmhafiz/go8/internal/domain/author"
)

func NewSearch(db *gen.Client) *repository {
	return &repository{ent: db}
}

// Search using the same store. May use other store e.g. elasticsearch/bleve as
// the search repository.
func (r *repository) Search(ctx context.Context, f *author.Filter) ([]*author.Schema, int, error) {
	tracer := otel.Tracer("")
	ctx, span := tracer.Start(ctx, "AuthorSearch")
	defer span.End()

	var predicateUser []predicate.Author
	if f.FirstName != "" {
		predicateUser = append(predicateUser, entAuthor.FirstNameContainsFold(f.FirstName))
	}
	if f.MiddleName != "" {
		predicateUser = append(predicateUser, entAuthor.MiddleNameContainsFold(f.MiddleName))
	}
	if f.LastName != "" {
		predicateUser = append(predicateUser, entAuthor.LastNameContainsFold(f.LastName))
	}

	total, err := r.ent.Author.Query().
		Where(entAuthor.DeletedAtIsNil()).
		Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving Author list: %w", err)
	}

	// The section where the search is done
	//
	// ILIKE is to search for case-insensitive in postgres
	//
	// Speed optimization is possible by creating an index on lowered names
	// Reference: https://www.postgresql.org/docs/current/indexes-expressional.html
	// CREATE INDEX on authors (LOWER(first_name));
	// CREATE INDEX on authors (LOWER(middle_name));
	// CREATE INDEX on authors (LOWER(last_name));
	//
	// Second alternative is to use ~*
	// mods = append(mods, qm.Or(models.AuthorColumns.MiddleName+" ~* ?", f.Name))
	//
	// postgres has a builtin full-text search: https://www.postgresql.org/docs/current/textsearch.html
	//
	// Also, may use term frequency-inverted index search (tf-idf) like
	// elasticsearch or bleve.
	authors, err := r.ent.Author.Query().
		WithBooks().
		Where(predicateUser...).
		Where(entAuthor.DeletedAtIsNil()).
		Limit(f.Base.Limit).
		Offset(f.Base.Offset).
		Order(authorOrder(f.Base.Sort)...).
		All(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("error retrieving Author list: %w", err)
	}

	var resp []*author.Schema

	for _, a := range authors {
		resp = append(resp, &author.Schema{
			ID:         a.ID,
			FirstName:  a.FirstName,
			MiddleName: a.MiddleName,
			LastName:   a.LastName,
			CreatedAt:  a.CreatedAt,
			UpdatedAt:  a.UpdatedAt,
			DeletedAt:  a.DeletedAt,
			Books:      nil,
		})
	}

	return resp, total, nil
}
