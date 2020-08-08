package authors

import (
	"context"
	"database/sql"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"eight/internal/models"
)

type store interface {
	All(context.Context) (models.AuthorSlice, error)
	GetAuthor(context.Context, int64) (*models.Author, error)
	CreateAuthor(ctx context.Context, authorID *models.Author) (*models.Author, error)
}

type authorStore struct {
	db    *sql.DB
	cache *redis.Client
}

func newStore(db *sql.DB, rdb *redis.Client) (*authorStore, error) {
	return &authorStore{
		db:    db,
		cache: rdb,
	}, nil
}

func (as *authorStore) All(ctx context.Context) (models.AuthorSlice, error) {
	authorSlice := models.AuthorSlice{}

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
	var author *models.Author

	foundAuthor, err := models.Authors(models.AuthorWhere.AuthorID.EQ(authorID)).One(ctx, as.db)
	log.Println(foundAuthor)
	if err != nil {
		return author, err
	}

	return foundAuthor, nil
}
