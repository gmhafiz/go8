package authors

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"go8ddd/internal/middleware"
	"go8ddd/internal/model"
)

type AuthorRepository interface {
	All(ctx context.Context) (model.AuthorSlice, error)
}

type authorRepo struct {
	log zerolog.Logger
	db  *sql.DB
}

func NewRepository(log zerolog.Logger, db *sql.DB) AuthorRepository {
	return &authorRepo{
		log: log,
		db:  db,
	}
}

func (r authorRepo) All(ctx context.Context) (model.AuthorSlice, error) {
	page := ctx.Value("pagination").(middleware.Pagination).Page
	size := ctx.Value("pagination").(middleware.Pagination).Size

	var err error
	var authors []*model.Author

	if page == 0 && size == 0 {
		authors, err = model.Authors().All(ctx, r.db)
	} else {
		authors, err = model.Authors(
			qm.OrderBy(`created_at DESC`),
			qm.Limit(size),
			qm.Offset(page-1)).
			All(ctx, r.db)
	}

	if err != nil {
		return nil, err
	}
	return authors, nil
}
