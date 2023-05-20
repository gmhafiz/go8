package csrf

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"

	"github.com/cespare/xxhash/v2"
)

// ValidToken Checks if CSRF token is valid
func ValidToken(ctx context.Context, db *sql.DB, token string) bool {
	hash, err := sum(token)
	if err != nil {
		return false
	}

	var exists bool
	row := db.QueryRowContext(ctx, `
			SELECT EXISTS(
				SELECT token FROM sessions 
				WHERE token = $1 
				  AND current_timestamp < expiry
			) `, hash)
	err = row.Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

// ValidAndDeleteToken deletes the token from the store if and only if token is valid.
// Useful for one-time token use.
func ValidAndDeleteToken(ctx context.Context, db *sql.DB, token string) error {
	hash, err := sum(token)
	if err != nil {
		return nil
	}

	res, err := db.ExecContext(ctx, `
		DELETE FROM sessions WHERE token = $1 AND current_timestamp < expiry
	`, hash)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return errors.New("token not found")
	}

	if rowsAffected != 1 {
		return errors.New("no csrf token was found")
	}
	return nil
}

func sum(token string) (string, error) {
	h := xxhash.New()
	_, err := h.Write([]byte(token))
	if err != nil {

		return "", err
	}
	sum := h.Sum(nil)
	str := hex.EncodeToString(sum)

	return str, err
}
