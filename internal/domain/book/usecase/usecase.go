package usecase

import (
	"context"

	"github.com/jinzhu/now"
	"github.com/volatiletech/null/v8"

	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/model"
)

type BookUseCase struct {
	bookRepo book.Repository
}

func (u *BookUseCase) Update(ctx context.Context, book *model.Book) (*model.Book, error) {
	return u.bookRepo.Update(ctx, book)
}

func NewBookUseCase(bookRepo book.Repository) *BookUseCase {
	return &BookUseCase{
		bookRepo: bookRepo,
	}
}

func (u *BookUseCase) Create(ctx context.Context, title, description, imageURL, publishedDate string) (*model.Book, error) {
	bk := &model.Book{
		Title:         title,
		PublishedDate: now.MustParse(publishedDate),
		ImageURL: null.String{
			String: imageURL,
			Valid:  true,
		},
		Description: null.String{
			String: description,
			Valid:  true,
		},
	}

	bookID, err := u.bookRepo.Create(ctx, bk)
	if err != nil {
		return nil, err
	}
	bookFound, err := u.bookRepo.Find(context.Background(), bookID)
	if err != nil {
		return nil, err
	}
	return bookFound, err
}

func (u *BookUseCase) All(ctx context.Context) ([]*model.Book, error) {
	return u.bookRepo.All(ctx)
}

func (u *BookUseCase) Find(ctx context.Context, bookID int64) (*model.Book, error) {
	return u.bookRepo.Find(ctx, bookID)
}

func (u *BookUseCase) Delete(ctx context.Context, bookID int64) error {
	return u.bookRepo.Delete(ctx, bookID)
}
