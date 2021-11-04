package database

import (
	"context"
	"fmt"

	"github.com/friendsofgo/errors"
	"github.com/jinzhu/now"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/gmhafiz/go8/internal/domain/author"
	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/models"
	"github.com/gmhafiz/go8/internal/utility/respond"
)

type repository struct {
	db *sqlx.DB
}

type Repository interface {
	Create(ctx context.Context, r *models.Author) (int64, error)
	CreateRead(ctx context.Context, r author.CreateRequest) (*author.CreateResponse, error)
	List(ctx context.Context, f *author.Filter) ([]*models.Author, int64, error)
	Read(ctx context.Context, authorID uint64) (*models.Author, error)
	Update(ctx context.Context, author *models.Author) (*models.Author, error)
	Delete(ctx context.Context, authorID int64) error
	ReadWithBooks(ctx context.Context, id uint64) (*author.WithBooks, error)
}

type Searcher interface {
	Search(ctx context.Context, f *author.Filter) ([]*models.Author, int64, error)
}

func New(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) ReadWithBooks(ctx context.Context, id uint64) (*author.WithBooks, error) {
	a, err := r.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	var books []book.Res
	rows, err := r.db.QueryContext(ctx, `SELECT b.book_id, b.title, b.published_date, b.image_url, b.description
FROM books b
         INNER JOIN book_authors ba on ba.book_id = b.book_id
         INNER JOIN authors a on a.id = ba.author_id
where a.id = $1`, a.ID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var b models.Book
		err = rows.Scan(&b.BookID, &b.Title, &b.PublishedDate, &b.ImageURL, &b.Description)
		if err != nil {
			return nil, fmt.Errorf("error retrieving book")
		}
		books = append(books, book.Res{
			BookID:        b.BookID,
			Title:         b.Title,
			PublishedDate: b.PublishedDate,
			ImageURL:      b.ImageURL,
			Description: null.String{
				String: b.Description,
				Valid:  true,
			},
		})
	}

	return &author.WithBooks{
		Author: a,
		Books:  books,
	}, nil
}

func (r *repository) Create(ctx context.Context, author *models.Author) (int64, error) {
	err := author.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return 0, errors.Wrapf(err, "author.repository.Create")
	}

	return author.ID, nil
}

func (r *repository) CreateRead(ctx context.Context, authorReq author.CreateRequest) (*author.CreateResponse, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		// This is a neat trick that rolls back transaction if any logic
		// before doing a tx.Commit() below fails.
		// This is to ensure we roll back an early return happens before
		// this transaction is committed.
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	newAuthor := &models.Author{
		FirstName: authorReq.FirstName,
		MiddleName: null.String{
			String: authorReq.MiddleName,
			Valid:  true,
		},
		LastName: authorReq.LastName,
	}
	err = newAuthor.Insert(ctx, tx, boil.Infer())
	if err != nil {
		return nil, err
	}
	var newBooks []author.Book
	for _, b := range authorReq.Books {
		pDate, _ := now.Parse(b.PublishedDate)
		newBook := &models.Book{
			BookID:        b.BookID,
			Title:         b.Title,
			PublishedDate: pDate,
			ImageURL: null.String{
				String: b.ImageURL,
				Valid:  true,
			},
			Description: b.Description,
		}
		err := newBook.Insert(ctx, tx, boil.Infer())
		if err != nil {
			return nil, err
		}
		newBooks = append(newBooks, author.Book{
			BookID:        newBook.BookID,
			Title:         newBook.Title,
			PublishedDate: newBook.PublishedDate.String(),
			ImageURL:      newBook.ImageURL.String,
			Description:   newBook.Description,
		})

		ba := &models.BookAuthor{
			BookID:   newBook.BookID,
			AuthorID: newAuthor.ID,
		}
		err = ba.Insert(ctx, tx, boil.Infer())
		if err != nil {
			return nil, err
		}
	}

	// any return before this line will not be committed and will be roll back
	// run by the `defer()` above. `defer()` run just before ths function exits.
	err = tx.Commit()
	if err != nil {
		return nil, errors.Wrapf(err, "error creating an author")
	}

	return &author.CreateResponse{
		ID:         newAuthor.ID,
		FirstName:  newAuthor.FirstName,
		MiddleName: newAuthor.MiddleName.String,
		LastName:   newAuthor.LastName,
		Books:      newBooks,
	}, nil
}

func (r *repository) List(ctx context.Context, f *author.Filter) ([]*models.Author, int64, error) {
	var mods []qm.QueryMod

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if f.Base.Limit != 0 && !f.Base.DisablePaging {
		mods = append(mods, qm.Limit(int(f.Base.Limit)))
	}
	if f.Base.Offset != 0 && !f.Base.DisablePaging {
		mods = append(mods, qm.Offset(f.Base.Offset))
	}
	query := models.Authors(mods...)
	total, err := query.Count(ctx, tx)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "error counting Authors")
	}

	mods = append(mods, qm.OrderBy(models.AuthorColumns.UpdatedAt))
	// todo: load many-2-many books relationship

	all, err := models.Authors(mods...).All(ctx, tx)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "error retrieving Author list")
	}

	err = tx.Commit()
	if err != nil {
		return nil, 0, errors.Wrapf(err, "error committing Author list")
	}

	return all, total, nil
}

func (r *repository) Read(ctx context.Context, authorID uint64) (*models.Author, error) {
	return models.FindAuthor(ctx, r.db, int64(authorID))
}

func (r *repository) Update(ctx context.Context, author *models.Author) (*models.Author, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	a, err := models.FindAuthor(ctx, tx, author.ID)
	if err != nil {
		return nil, err
	}

	a.FirstName = author.FirstName
	a.MiddleName = author.MiddleName
	a.LastName = author.LastName

	_, err = a.Update(ctx, tx, boil.Infer())
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.Wrapf(err, "error committing Author list")
	}

	return a, nil
}

func (r *repository) Delete(ctx context.Context, authorID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	found, err := models.FindAuthor(ctx, tx, authorID)
	if err != nil {
		return respond.ErrNoRecord
	}

	_, err = found.Delete(ctx, tx, false)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrapf(err, "error deleting Author")
	}
	return nil
}
