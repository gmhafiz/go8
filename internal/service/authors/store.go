package authors

import (
	"context"
	"database/sql"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"eight/internal/models"
)

type store interface {
	All(context.Context) (models.AuthorSlice, error)
	GetAuthor(context.Context, int64) (*models.Author, error)
	CreateAuthor(ctx context.Context, authorID *models.Author) (*models.Author, error)
}

type authorStore struct {
	qbuilder squirrel.StatementBuilderType
	pqdriver *pgxpool.Pool
	db       *sql.DB
}

func newStore(pqdriver *pgxpool.Pool, db *sql.DB) (*authorStore, error) {
	return &authorStore{
		qbuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		pqdriver: pqdriver,
		db:       db,
	}, nil
}

func (as *authorStore) All(ctx context.Context) (models.AuthorSlice, error) {
	authorSlice := models.AuthorSlice{}
	boil.DebugMode = true

	authorSlice, err := models.Authors().All(ctx, as.db)
	if err != nil {
		log.Println(err)
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
	boil.DebugMode = true
	var author *models.Author

	foundAuthor, err := models.Authors(models.AuthorWhere.AuthorID.EQ(authorID)).One(ctx, as.db)
	log.Println(foundAuthor)
	if err != nil {
		return author, err
	}

	return foundAuthor, nil
}
