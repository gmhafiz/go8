package authors

import (
	"context"
	"database/sql"
	"eight/internal/middleware"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"log"

	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"eight/internal/models"
)

type store interface {
	All(context.Context) (models.AuthorSlice, error)
	GetAuthor(context.Context, int64) (*models.Author, error)
	CreateAuthor(ctx context.Context, authorID *models.Author) (*models.Author, error)
}

type authorStore struct {
	db     *sql.DB
	logger zerolog.Logger
}

func newStore(db *sql.DB, logger zerolog.Logger) (*authorStore, error) {
	return &authorStore{
		db:     db,
		logger: logger,
	}, nil
}

func (as *authorStore) All(ctx context.Context) (models.AuthorSlice, error) {
	page := ctx.Value("pagination").(middleware.Pagination).Page
	size := ctx.Value("pagination").(middleware.Pagination).Size

	authorSlice := models.AuthorSlice{}

	var err error
	if page != 0 && size != 0 {
		authorSlice, err = models.Authors(qm.OrderBy(`created_at DESC`), qm.Limit(size),
			qm.Offset(page-1)).All(ctx, as.db)
	} else {
		authorSlice, err = models.Authors().All(ctx, as.db)
	}

	if err != nil {
		as.logger.Error().Msg(err.Error())
		return nil, err

	}
	return authorSlice, nil
}

func (as *authorStore) CreateAuthor(ctx context.Context, author *models.Author) (*models.Author, error) {
	err := author.Insert(ctx, as.db, boil.Infer())
	if err != nil {
		return author, err
	}
	return author, nil
}

func (as *authorStore) GetAuthor(ctx context.Context, authorID int64) (*models.Author, error) {
	var author *models.Author

	foundAuthor, err := models.Authors(models.AuthorWhere.AuthorID.EQ(authorID)).One(ctx, as.db)
	log.Println(foundAuthor)
	if err != nil {
		return author, err
	}

	return foundAuthor, nil
}
