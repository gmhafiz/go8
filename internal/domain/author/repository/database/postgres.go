package database

import (
	"context"
	"fmt"
	stime "time"

	"github.com/friendsofgo/errors"

	"github.com/gmhafiz/go8/ent/gen"
	entAuthor "github.com/gmhafiz/go8/ent/gen/author"
	"github.com/gmhafiz/go8/ent/gen/predicate"
	"github.com/gmhafiz/go8/internal/domain/author"
	"github.com/gmhafiz/go8/internal/utility/time"
)

type repository struct {
	ent *gen.Client
}

type Repository interface {
	Create(ctx context.Context, r author.CreateRequest) (*gen.Author, error)
	List(ctx context.Context, f *author.Filter) ([]*gen.Author, int64, error)
	Read(ctx context.Context, id uint64) (*gen.Author, error)
	Update(ctx context.Context, toAuthor *author.Update) (*gen.Author, error)
	Delete(ctx context.Context, authorID int64) error
}

type Searcher interface {
	Search(ctx context.Context, f *author.Filter) ([]*gen.Author, int64, error)
}

func New(ent *gen.Client) *repository {
	return &repository{
		ent: ent,
	}
}

func (r *repository) Create(ctx context.Context, author author.CreateRequest) (*gen.Author, error) {
	bulk := make([]*gen.BookCreate, len(author.Books))
	for i, b := range author.Books {
		bulk[i] = r.ent.Book.Create().
			SetTitle(b.Title).
			SetDescription(b.Description).
			SetPublishedDate(time.Parse(b.PublishedDate))
	}
	books, err := r.ent.Book.CreateBulk(bulk...).Save(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "author.repository.Create bulk books")
	}

	created, err := r.ent.Author.Create().
		SetFirstName(author.FirstName).
		SetNillableMiddleName(&author.MiddleName).
		SetLastName(author.LastName).
		AddBooks(books...).
		Save(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "author.repository.Create")
	}

	created.Edges.Books = books

	return created, nil
}

func (r *repository) List(ctx context.Context, f *author.Filter) ([]*gen.Author, int64, error) {
	// filter by first and last names, if exists
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

	// sort by field
	var orderFunc []gen.OrderFunc
	for col, ord := range f.Base.Sort {
		if ord == "ASC" {
			orderFunc = append(orderFunc, gen.Asc(col))
		} else {
			orderFunc = append(orderFunc, gen.Desc(col))
		}
	}

	total, err := r.ent.Author.Query().Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	authors, err := r.ent.Author.Query().
		WithBooks().
		Where(predicateUser...).
		Where(entAuthor.DeletedAtIsNil()).
		Limit(int(f.Base.Limit)).
		Offset(f.Base.Offset).
		Order(orderFunc...).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return authors, int64(total), err
}

func (r *repository) Read(ctx context.Context, id uint64) (*gen.Author, error) {
	found, err := r.ent.Author.Query().
		WithBooks().
		Where(entAuthor.ID(uint(id))).
		First(ctx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving book: %w", err)
	}

	return found, err
}

func (r *repository) Update(ctx context.Context, author *author.Update) (*gen.Author, error) {
	updated, err := r.ent.Author.UpdateOne(&gen.Author{
		ID:         uint(author.ID),
		FirstName:  author.FirstName,
		MiddleName: author.MiddleName,
		LastName:   author.LastName,
	}).Save(ctx)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (r *repository) Delete(ctx context.Context, authorID int64) error {
	_, err := r.ent.Debug().Author.UpdateOneID(uint(authorID)).
		SetDeletedAt(stime.Now()).
		Save(ctx)
	return err

	//return r.ent.Author.DeleteOneID(uint(authorID)).Exec(ctx)
}
