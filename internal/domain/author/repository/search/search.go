package search

import (
	"context"

	"github.com/friendsofgo/errors"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/gmhafiz/go8/ent/gen"
	"github.com/gmhafiz/go8/internal/domain/author"
	"github.com/gmhafiz/go8/internal/models"
)

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *repository {
	return &repository{db: db}
}

// Search using the same store. May use other store e.g. elasticsearch/bleve as
// the search repository.
func (r *repository) Search(ctx context.Context, f *author.Filter) ([]*gen.Author, int64, error) {
	var mods []qm.QueryMod

	if f.Base.Limit != 0 {
		mods = append(mods, qm.Limit(int(f.Base.Limit)))
	}
	if f.Base.Offset != 0 {
		mods = append(mods, qm.Offset(f.Base.Offset))
	}
	query := models.Authors(mods...)
	total, err := query.Count(ctx, r.db)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "error counting Authors")
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
	mods = append(mods, qm.Where(models.AuthorColumns.FirstName+" ILIKE ?", f.FirstName))
	mods = append(mods, qm.Or(models.AuthorColumns.MiddleName+" ILIKE ?", f.MiddleName))
	mods = append(mods, qm.Or(models.AuthorColumns.LastName+" ILIKE ?", f.LastName))

	mods = append(mods, qm.OrderBy(models.AuthorColumns.UpdatedAt))

	all, err := models.Authors(mods...).All(ctx, r.db)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "error retrieving Author list")
	}

	if err != nil {
		return nil, 0, errors.Wrapf(err, "error committing Author list")
	}

	var authors []*gen.Author
	for _, val := range all {
		a := &gen.Author{
			ID:         uint(val.ID),
			FirstName:  val.FirstName,
			MiddleName: val.MiddleName.String,
			LastName:   val.LastName,
			CreatedAt:  val.CreatedAt.Time,
			UpdatedAt:  val.UpdatedAt.Time,
			DeletedAt:  &val.DeletedAt.Time,
		}
		authors = append(authors, a)
	}

	return authors, total, nil
}
