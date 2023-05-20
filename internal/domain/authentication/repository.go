package authentication

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/gmhafiz/scs/v2"

	"github.com/gmhafiz/go8/ent/gen"
	"github.com/gmhafiz/go8/ent/gen/session"
	"github.com/gmhafiz/go8/ent/gen/user"
)

type repo struct {
	ent     *gen.Client
	db      *sql.DB
	session *scs.SessionManager
}

var (
	ErrEmailNotAvailable = errors.New("email is not available")
	ErrNotLoggedIn       = errors.New("you are not logged in yet")
)

type Repo interface {
	Register(ctx context.Context, firstName, lastName, email, hashedPassword string) error
	Login(ctx context.Context, req LoginRequest) (*gen.User, bool, error)
	Logout(ctx context.Context, userID uint64) (bool, error)
	Csrf(ctx context.Context) (string, error)
}

func (r *repo) Register(ctx context.Context, firstName, lastName, email, hashedPassword string) error {
	err := r.ent.User.Create().
		SetFirstName(firstName).
		SetLastName(lastName).
		SetEmail(email).
		SetPassword(hashedPassword).
		Exec(ctx)
	if err != nil {
		if gen.IsConstraintError(err) {
			return ErrEmailNotAvailable
		}
		return err
	}

	return nil
}

func (r *repo) Login(ctx context.Context, req LoginRequest) (*gen.User, bool, error) {
	u, err := r.ent.User.Query().Where(user.EmailEqualFold(req.Email)).First(ctx)
	if err != nil {
		return nil, false, err
	}

	match, err := argon2id.ComparePasswordAndHash(req.Password, u.Password)
	if err != nil {
		return nil, false, errors.New("wrong password is provided")
	}

	return u, match, nil
}

func (r *repo) Logout(ctx context.Context, userID uint64) (bool, error) {
	var found bool
	rows := r.db.QueryRowContext(ctx, `
			   SELECT CASE
			   WHEN EXISTS(SELECT *
						   FROM sessions
						   WHERE sessions.user_id = $1)
				   THEN true
			   ELSE false
			   END
	;
	`, userID)
	err := rows.Scan(&found)
	if err != nil {
		return false, err
	}

	if !found {
		return false, ErrNotLoggedIn
	}

	_, err = r.ent.Session.Delete().Where(session.UserIDEQ(userID)).Exec(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *repo) Csrf(ctx context.Context) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}

	err = r.session.CtxStore.CommitCtx(ctx, token, []byte("csrf_token"), time.Now().Add(r.session.Lifetime))
	if err != nil {
		return "", err
	}

	return token, nil
}

func NewRepo(ent *gen.Client, db *sql.DB, manager *scs.SessionManager) *repo {
	return &repo{
		ent:     ent,
		db:      db,
		session: manager,
	}
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
