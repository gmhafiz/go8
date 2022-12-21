package message

import "errors"

var (
	ErrBadRequest    = errors.New("error bad request")
	ErrInternalError = errors.New("error internal")

	ErrFormingResponse = errors.New("error forming response")

	ErrNoRecord = errors.New("no record found")

	ErrFetchingBook = errors.New("error fetching books")
)
