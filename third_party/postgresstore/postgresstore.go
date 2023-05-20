// Package postgresstore is a custom postgres implementation of github.com/alexedwards/scs library.
// It is nearly identical to it except:
//
//  1. It saves a uint64 data along with the session data for the purpose of user session invalidation.
//  2. Tokens are hashed before being saved into the database.
//
// The schema is identical to scs library but with added `user_id` foreign key column:
//
//	CREATE TABLE IF NOT EXISTS sessions
//	(
//	    token   TEXT PRIMARY KEY,
//	    user_id BIGINT      NOT NULL CONSTRAINT session_user_fk REFERENCES users ON DELETE CASCADE ,
//	    data    BYTEA       NOT NULL,
//	    expiry  TIMESTAMPTZ NOT NULL
//	);
//
// If number of records in `expiry` column is large, can consider indexing it using BRIN index
//
//	CREATE INDEX sessions_expiry_idx ON sessions USING brin (expiry);
package postgresstore

import (
	"context"
	"database/sql"
	"encoding/hex"
	"log"
	"time"

	"github.com/cespare/xxhash/v2"

	"github.com/gmhafiz/go8/internal/middleware"
)

// PostgresStore represents the session store.
type PostgresStore struct {
	db          *sql.DB
	stopCleanup chan bool
}

func (p *PostgresStore) Delete(token string) (err error) {
	panic("missing context arg")
}

func (p *PostgresStore) Find(token string) (b []byte, found bool, err error) {
	panic("missing context arg")
}

func (p *PostgresStore) Commit(token string, b []byte, expiry time.Time) (err error) {
	panic("missing context arg")
}

// FindCtx returns the data for a given session token from the PostgresStore instance.
// If the session token is not found or is expired, the returned exists flag will
// be set to false.
func (p *PostgresStore) FindCtx(ctx context.Context, token string) (b []byte, exists bool, err error) {
	hash, err := sum(token)
	if err != nil {
		return nil, false, err
	}

	row := p.db.QueryRowContext(ctx, `
		SELECT data FROM sessions 
            WHERE token = $1 
              AND current_timestamp < expiry 
            ORDER BY expiry desc`, hash)
	err = row.Scan(&b)
	if err == sql.ErrNoRows {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	}
	return b, true, nil
}

// CommitCtx adds a session token and data to the PostgresStore instance with the
// given expiry time. If the session token already exists, then the data and expiry
// time are updated. Hashed token is stored into database. User ID is retrieved from request
// context since modifying method signature will no longer implements scs's Store interface.
func (p *PostgresStore) CommitCtx(ctx context.Context, token string, b []byte, expiry time.Time) error {
	var userID any
	userID, ok := ctx.Value(middleware.KeyID).(uint64)
	if !ok {
		userID = nil
	}

	hash, err := sum(token)
	if err != nil {
		return err
	}

	_, err = p.db.ExecContext(ctx, `
		INSERT INTO sessions (token, user_id, data, expiry) 
		VALUES ($1, $2, $3, $4) 
		ON CONFLICT (token) 
			DO UPDATE 
			SET data = EXCLUDED.data, 
				expiry = EXCLUDED.expiry
				`, hash, userID, b, expiry)
	if err != nil {
		return err
	}
	return nil
}

// DeleteCtx removes a session token and corresponding data from the PostgresStore
// instance.
func (p *PostgresStore) DeleteCtx(ctx context.Context, token string) error {
	hash, err := sum(token)
	if err != nil {
		return err
	}

	_, err = p.db.ExecContext(ctx, "DELETE FROM sessions WHERE token = $1", hash)
	return err
}

// AllCtx returns a map containing the token and data for all active (i.e.
// not expired) sessions in the PostgresStore instance.
func (p *PostgresStore) AllCtx(ctx context.Context) (map[string][]byte, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT token, data FROM sessions WHERE current_timestamp < expiry")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := make(map[string][]byte)

	for rows.Next() {
		var (
			token string
			data  []byte
		)

		err = rows.Scan(&token, &data)
		if err != nil {
			return nil, err
		}

		sessions[token] = data
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

// New returns a new PostgresStore instance, with a background cleanup goroutine
// that runs every 5 minutes to remove expired session data.
func New(db *sql.DB) *PostgresStore {
	return NewWithCleanupInterval(db, 5*time.Minute)
}

// NewWithCleanupInterval returns a new PostgresStore instance. The cleanupInterval
// parameter controls how frequently expired session data is removed by the
// background cleanup goroutine. Setting it to 0 prevents the cleanup goroutine
// from running (i.e. expired sessions will not be removed).
func NewWithCleanupInterval(db *sql.DB, cleanupInterval time.Duration) *PostgresStore {
	p := &PostgresStore{db: db}
	if cleanupInterval > 0 {
		go p.startCleanup(cleanupInterval)
	}
	return p
}

func (p *PostgresStore) startCleanup(interval time.Duration) {
	p.stopCleanup = make(chan bool)
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			err := p.deleteExpired()
			if err != nil {
				log.Println(err)
			}
		case <-p.stopCleanup:
			ticker.Stop()
			return
		}
	}
}

// StopCleanup terminates the background cleanup goroutine for the PostgresStore
// instance. It's rare to terminate this; generally PostgresStore instances and
// their cleanup goroutines are intended to be long-lived and run for the lifetime
// of your application.
//
// There may be occasions though when your use of the PostgresStore is transient.
// An example is creating a new PostgresStore instance in a test function. In this
// scenario, the cleanup goroutine (which will run forever) will prevent the
// PostgresStore object from being garbage collected even after the test function
// has finished. You can prevent this by manually calling StopCleanup.
func (p *PostgresStore) StopCleanup() {
	if p.stopCleanup != nil {
		p.stopCleanup <- true
	}
}

func (p *PostgresStore) deleteExpired() error {
	_, err := p.db.Exec("DELETE FROM sessions WHERE expiry < current_timestamp")
	return err
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
