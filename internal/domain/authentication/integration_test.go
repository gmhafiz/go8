package authentication

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	entsql "entgo.io/ent/dialect/sql"
	"github.com/alexedwards/argon2id"
	"github.com/gmhafiz/go8/internal/utility/csrf"
	"github.com/gmhafiz/scs/v2"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"

	"github.com/gmhafiz/go8/database"
	"github.com/gmhafiz/go8/ent/gen"
	"github.com/gmhafiz/go8/internal/middleware"
	"github.com/gmhafiz/go8/third_party/postgresstore"
)

const (
	DBDriver = "postgres"

	sessionName = "session"
)

var (
	migrator *database.Migrate
)

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "17",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseURL := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)
	log.Println(databaseURL)

	_ = resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	var db *sql.DB

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err = sql.Open(DBDriver, databaseURL)
		if err != nil {
			log.Println(err)
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	migrator = database.Migrator(db, database.WithDSN(databaseURL))

	// Performing a migration this way means all tests in this package shares
	// the same db schema across all unit test.
	// If isolation is needed, then do away with using `testing.M`. Do a
	// migration for each test handler instead.
	migrator.Up()

	// Seed with super admin suer
	hashedPassword, err := argon2id.CreateHash("highEntropyPassword", argon2id.DefaultParams)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = migrator.DB.ExecContext(context.Background(), `
		INSERT INTO users (email, password) VALUES ($1, $2)
		ON CONFLICT (email) DO NOTHING 
		`, "admin@gmhafiz.com", hashedPassword)
	if err != nil {
		log.Fatalln(err)
	}

	// We can access database with m.hostAndPort or m.databaseURL
	// port changes everytime a new docker instance is run
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestHandler_RegisterIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	type args struct {
		*RegisterRequest
	}
	type want struct {
		error
		status int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "simple",
			args: args{
				RegisterRequest: &RegisterRequest{
					Email:    "email@example.com",
					Password: "highEntropyPassword",
				},
			},
			want: want{
				error:  nil,
				status: http.StatusCreated,
			},
		},
		{
			name: "email already registered",
			args: args{
				RegisterRequest: &RegisterRequest{
					Email:    "email@example.com",
					Password: "highEntropyPassword",
				},
			},
			want: want{
				error:  ErrEmailNotAvailable,
				status: http.StatusBadRequest,
			},
		},
		{
			name: "no email is supplied",
			args: args{
				RegisterRequest: &RegisterRequest{
					Email:    "",
					Password: "highEntropyPassword",
				},
			},
			want: want{
				error:  ErrEmailRequired,
				status: http.StatusBadRequest,
			},
		},
		{
			name: "no password is supplied",
			args: args{
				RegisterRequest: &RegisterRequest{
					Email:    "email@example.com",
					Password: "",
				},
			},
			want: want{
				error:  ErrPasswordLength,
				status: http.StatusBadRequest,
			},
		},
	}

	client := dbClient()
	session := newSession(migrator.DB, 1*time.Hour)
	repo := NewRepo(client, migrator.DB, session)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			var err error
			if tt.args.RegisterRequest != nil {
				err = json.NewEncoder(&buf).Encode(tt.args.RegisterRequest)
			}
			assert.Nil(t, err)

			rr := httptest.NewRequest(http.MethodPost, "/api/v1/register", &buf)
			ww := httptest.NewRecorder()

			router := chi.NewRouter()
			router.Use(middleware.LoadAndSave(session))

			RegisterHTTPEndPoints(router, session, repo)

			router.ServeHTTP(ww, rr)

			assert.Equal(t, tt.want.status, ww.Code)

			b, err := io.ReadAll(ww.Body)
			assert.Nil(t, err)

			if len(b) > 0 {
				errStruct := struct {
					Message string `json:"message"`
				}{
					Message: string(b),
				}

				err = json.Unmarshal(b, &errStruct)
				assert.Nil(t, err)

				assert.Equal(t, tt.want.error.Error(), errStruct.Message)
			}
		})
	}
}

func TestHandler_LoginIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	type args struct {
		*LoginRequest
	}
	type want struct {
		error
		status int
		token  struct{ Token string }
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "simple",
			args: args{
				LoginRequest: &LoginRequest{
					Email:    "email@example.com",
					Password: "highEntropyPassword",
				},
			},
			want: want{
				error:  nil,
				status: http.StatusOK,
				token: struct {
					Token string
				}{
					Token: "",
				},
			},
		},
		{
			name: "not registered",
			args: args{
				LoginRequest: &LoginRequest{
					Email:    "email@example.XXX", // non-existent email
					Password: "highEntropyPassword",
				},
			},
			want: want{
				error:  nil,
				status: http.StatusUnauthorized,
				token: struct {
					Token string
				}{
					Token: "",
				},
			},
		},
	}

	client := dbClient()
	session := newSession(migrator.DB, 1*time.Hour)
	repo := NewRepo(client, migrator.DB, session)

	hashedPassword, err := argon2id.CreateHash("highEntropyPassword", argon2id.DefaultParams)
	assert.Nil(t, err)

	_, err = repo.db.ExecContext(context.Background(), `
		INSERT INTO users (email, password) VALUES ($1, $2)
		ON CONFLICT (email) DO NOTHING 
		`, "email@example.com", hashedPassword)
	assert.Nil(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			var err error
			if tt.args.LoginRequest != nil {
				err = json.NewEncoder(&buf).Encode(tt.args.LoginRequest)
			}
			assert.Nil(t, err)

			rr := httptest.NewRequest(http.MethodPost, "/api/v1/login", &buf)
			ww := httptest.NewRecorder()

			router := chi.NewRouter()
			router.Use(middleware.LoadAndSave(session))

			RegisterHTTPEndPoints(router, session, repo)

			router.ServeHTTP(ww, rr)

			assert.Equal(t, tt.want.status, ww.Code)

			assert.NotNil(t, ww.Header().Get("Set-Cookie"))
		})
	}
}

func TestHandler_ProtectedIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	type args struct {
		*LoginRequest
	}
	type want struct {
		error
		status int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "normal",
			args: args{
				LoginRequest: &LoginRequest{
					Email:    "email@example.com",
					Password: "highEntropyPassword",
				},
			},
			want: want{
				error:  nil,
				status: http.StatusOK,
			},
		},
	}

	client := dbClient()
	session := newSession(migrator.DB, 1*time.Hour)
	repo := NewRepo(client, migrator.DB, session)

	hashedPassword, err := argon2id.CreateHash("highEntropyPassword", argon2id.DefaultParams)
	assert.Nil(t, err)

	_, err = repo.db.ExecContext(context.Background(), `
		INSERT INTO users (email, password) VALUES ($1, $2)
		ON CONFLICT (email) DO NOTHING 
		`, "email@example.com", hashedPassword)
	assert.Nil(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			var err error
			if tt.args.LoginRequest != nil {
				err = json.NewEncoder(&buf).Encode(tt.args.LoginRequest)
			}
			assert.Nil(t, err)

			rr := httptest.NewRequest(http.MethodPost, "/api/v1/login", &buf)
			ww := httptest.NewRecorder()

			router := chi.NewRouter()
			router.Use(middleware.LoadAndSave(session))

			RegisterHTTPEndPoints(router, session, repo)

			router.ServeHTTP(ww, rr)

			assert.Equal(t, tt.want.status, ww.Code)

			assert.NotNil(t, ww.Header().Get("Set-Cookie"))
			token, err := extractToken(ww.Header().Get("Set-Cookie"))
			assert.Nil(t, err)

			rr = httptest.NewRequest(http.MethodGet, "/api/v1/restricted", nil)
			ww = httptest.NewRecorder()

			rr.AddCookie(&http.Cookie{
				Name:  sessionName,
				Value: token,
			})

			router = chi.NewRouter()
			router.Use(middleware.LoadAndSave(session))
			RegisterHTTPEndPoints(router, session, repo)
			router.ServeHTTP(ww, rr)

			assert.Equal(t, tt.want.status, ww.Code)
		})
	}
}

func TestHandler_MeIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	type args struct {
		*LoginRequest
	}
	type want struct {
		error
		status int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "normal",
			args: args{
				LoginRequest: &LoginRequest{
					Email:    "email@example.com",
					Password: "highEntropyPassword",
				},
			},
			want: want{
				error:  nil,
				status: http.StatusOK,
			},
		},
	}

	client := dbClient()
	session := newSession(migrator.DB, 1*time.Hour)
	repo := NewRepo(client, migrator.DB, session)

	hashedPassword, err := argon2id.CreateHash("highEntropyPassword", argon2id.DefaultParams)
	assert.Nil(t, err)

	_, err = repo.db.ExecContext(context.Background(), `
		INSERT INTO users (email, password) VALUES ($1, $2)
		ON CONFLICT (email) DO NOTHING 
		`, "email@example.com", hashedPassword)
	assert.Nil(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			var err error
			if tt.args.LoginRequest != nil {
				err = json.NewEncoder(&buf).Encode(tt.args.LoginRequest)
			}
			assert.Nil(t, err)

			rr := httptest.NewRequest(http.MethodPost, "/api/v1/login", &buf)
			ww := httptest.NewRecorder()

			router := chi.NewRouter()
			router.Use(middleware.LoadAndSave(session))

			RegisterHTTPEndPoints(router, session, repo)

			router.ServeHTTP(ww, rr)

			assert.Equal(t, tt.want.status, ww.Code)

			assert.NotNil(t, ww.Header().Get("Set-Cookie"))
			token, err := extractToken(ww.Header().Get("Set-Cookie"))
			assert.Nil(t, err)

			rr = httptest.NewRequest(http.MethodGet, "/api/v1/restricted/me", nil)
			ww = httptest.NewRecorder()

			rr.AddCookie(&http.Cookie{
				Name:  sessionName,
				Value: token,
			})

			router = chi.NewRouter()
			router.Use(middleware.LoadAndSave(session))
			RegisterHTTPEndPoints(router, session, repo)
			router.ServeHTTP(ww, rr)

			assert.Equal(t, tt.want.status, ww.Code)

			b, err := io.ReadAll(ww.Body)
			assert.Nil(t, err)

			if len(b) > 0 {
				type userID struct {
					UserID int `json:"user_id"`
				}
				var responseUserID userID
				err = json.Unmarshal(b, &responseUserID)
				assert.Nil(t, err)

				// There already is a super admin account created in the seed.
				// So this user ID is the next one which is 2
				assert.Equal(t, 2, responseUserID.UserID)
			}
		})
	}
}

func TestHandler_LogoutIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	type args struct {
		*LoginRequest
	}
	type want struct {
		error
		status int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "normal",
			args: args{
				LoginRequest: &LoginRequest{
					Email:    "email@example.com",
					Password: "highEntropyPassword",
				},
			},
			want: want{
				error:  nil,
				status: http.StatusOK,
			},
		},
	}

	client := dbClient()
	session := newSession(migrator.DB, 1*time.Hour)
	repo := NewRepo(client, migrator.DB, session)

	hashedPassword, err := argon2id.CreateHash("highEntropyPassword", argon2id.DefaultParams)
	assert.Nil(t, err)

	_, err = repo.db.ExecContext(context.Background(), `
		INSERT INTO users (email, password) VALUES ($1, $2)
		ON CONFLICT (email) DO NOTHING 
		`, "email@example.com", hashedPassword)
	assert.Nil(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			var err error
			if tt.args.LoginRequest != nil {
				err = json.NewEncoder(&buf).Encode(tt.args.LoginRequest)
			}
			assert.Nil(t, err)

			rr := httptest.NewRequest(http.MethodPost, "/api/v1/login", &buf)
			ww := httptest.NewRecorder()

			router := chi.NewRouter()
			router.Use(middleware.LoadAndSave(session))

			RegisterHTTPEndPoints(router, session, repo)
			router.ServeHTTP(ww, rr)

			assert.Equal(t, tt.want.status, ww.Code)

			assert.NotNil(t, ww.Header().Get("Set-Cookie"))
			token, err := extractToken(ww.Header().Get("Set-Cookie"))
			assert.Nil(t, err)

			rr = httptest.NewRequest(http.MethodGet, "/api/v1/restricted", nil)
			ww = httptest.NewRecorder()

			rr.AddCookie(&http.Cookie{
				Name:  sessionName,
				Value: token,
			})

			router.ServeHTTP(ww, rr)

			assert.Equal(t, tt.want.status, ww.Code)

			rr = httptest.NewRequest(http.MethodPost, "/api/v1/logout", nil)
			ww = httptest.NewRecorder()

			rr.AddCookie(&http.Cookie{
				Name:  sessionName,
				Value: token,
			})

			router.ServeHTTP(ww, rr)

			token, err = extractToken(ww.Header().Get("Set-Cookie"))
			assert.Nil(t, err)
			assert.Equal(t, token, "")
		})
	}
}

func TestHandler_Force_LogoutIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	type args struct {
		*LoginRequest
	}
	type want struct {
		error
		status int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "normal",
			args: args{
				LoginRequest: &LoginRequest{
					Email:    "email@example.com",
					Password: "highEntropyPassword",
				},
			},
			want: want{
				error:  nil,
				status: http.StatusOK,
			},
		},
	}

	client := dbClient()
	session := newSession(migrator.DB, 1*time.Hour)
	repo := NewRepo(client, migrator.DB, session)

	hashedPassword, err := argon2id.CreateHash("highEntropyPassword", argon2id.DefaultParams)
	assert.Nil(t, err)

	// Create normal user
	_, err = repo.db.ExecContext(context.Background(), `
		INSERT INTO users (email, password) VALUES ($1, $2)
		ON CONFLICT (email) DO NOTHING 
		`, "email@example.com", hashedPassword)
	assert.Nil(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			var err error
			if tt.args.LoginRequest != nil {
				err = json.NewEncoder(&buf).Encode(tt.args.LoginRequest)
			}
			assert.Nil(t, err)

			rr := httptest.NewRequest(http.MethodPost, "/api/v1/login", &buf)
			ww := httptest.NewRecorder()

			router := chi.NewRouter()
			router.Use(middleware.LoadAndSave(session))

			RegisterHTTPEndPoints(router, session, repo)
			router.ServeHTTP(ww, rr)

			assert.Equal(t, tt.want.status, ww.Code)

			assert.NotNil(t, ww.Header().Get("Set-Cookie"))
			token, err := extractToken(ww.Header().Get("Set-Cookie"))
			assert.Nil(t, err)

			rr = httptest.NewRequest(http.MethodGet, "/api/v1/restricted", nil)
			ww = httptest.NewRecorder()

			rr.AddCookie(&http.Cookie{
				Name:  sessionName,
				Value: token,
			})

			router.ServeHTTP(ww, rr)

			assert.Equal(t, tt.want.status, ww.Code)

			// Now we log in as super admin to log out admin@gmail.com
			admin := &LoginRequest{
				Email:    "admin@gmhafiz.com",
				Password: "highEntropyPassword",
			}
			err = json.NewEncoder(&buf).Encode(admin)
			assert.Nil(t, err)

			rr = httptest.NewRequest(http.MethodPost, "/api/v1/login", &buf)
			ww = httptest.NewRecorder()

			router.ServeHTTP(ww, rr)

			assert.NotNil(t, ww.Header().Get("Set-Cookie"))
			adminToken, err := extractToken(ww.Header().Get("Set-Cookie"))
			assert.Nil(t, err)

			// ID 2 is our normal user to be forced-log out
			rr = httptest.NewRequest(http.MethodPost, "/api/v1/restricted/logout/2", nil)
			ww = httptest.NewRecorder()
			rr.AddCookie(&http.Cookie{
				Name:  sessionName,
				Value: adminToken,
			})

			router.ServeHTTP(ww, rr)

			assert.Equal(t, http.StatusOK, ww.Code)

			// Check normal user ID 2 cannot access restricted route anymore
			rr = httptest.NewRequest(http.MethodGet, "/api/v1/restricted", nil)
			ww = httptest.NewRecorder()
			rr.AddCookie(&http.Cookie{
				Name:  sessionName,
				Value: token,
			})

			router.ServeHTTP(ww, rr)

			assert.Equal(t, http.StatusUnauthorized, ww.Code)
		})
	}
}

func TestHandler_Csrf_Valid_TokenIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	type args struct {
		*LoginRequest
	}
	type want struct {
		error
		status            int
		response          RespondCsrf
		csrfTokenValidity bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "normal",
			args: args{
				LoginRequest: &LoginRequest{
					Email:    "email@example.com",
					Password: "highEntropyPassword",
				},
			},
			want: want{
				error:             nil,
				status:            http.StatusOK,
				response:          RespondCsrf{CsrfToken: ""},
				csrfTokenValidity: true,
			},
		},
	}

	client := dbClient()
	session := newSession(migrator.DB, 100*time.Millisecond) // short expiry to test token expiry means test complete faster.
	repo := NewRepo(client, migrator.DB, session)

	hashedPassword, err := argon2id.CreateHash("highEntropyPassword", argon2id.DefaultParams)
	assert.Nil(t, err)

	_, err = repo.db.ExecContext(context.Background(), `
		INSERT INTO users (email, password) VALUES ($1, $2)
		ON CONFLICT (email) DO NOTHING 
		`, "email@example.com", hashedPassword)
	assert.Nil(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			var err error
			if tt.args.LoginRequest != nil {
				err = json.NewEncoder(&buf).Encode(tt.args.LoginRequest)
			}
			assert.Nil(t, err)

			rr := httptest.NewRequest(http.MethodPost, "/api/v1/login", &buf)
			ww := httptest.NewRecorder()

			router := chi.NewRouter()
			router.Use(middleware.LoadAndSave(session))

			RegisterHTTPEndPoints(router, session, repo)
			router.ServeHTTP(ww, rr)

			assert.Equal(t, tt.want.status, ww.Code)

			assert.NotNil(t, ww.Header().Get("Set-Cookie"))
			token, err := extractToken(ww.Header().Get("Set-Cookie"))
			assert.Nil(t, err)

			rr = httptest.NewRequest(http.MethodGet, "/api/v1/restricted/csrf", nil)
			ww = httptest.NewRecorder()

			rr.AddCookie(&http.Cookie{
				Name:  sessionName,
				Value: token,
			})

			router.ServeHTTP(ww, rr)

			assert.Equal(t, ww.Code, http.StatusOK)

			b, err := io.ReadAll(ww.Body)
			assert.Nil(t, err)

			var resp RespondCsrf
			err = json.Unmarshal(b, &resp)
			assert.Nil(t, err)

			assert.NotNil(t, resp.CsrfToken)

			validity := csrf.ValidToken(context.Background(), migrator.DB, resp.CsrfToken)
			assert.Equal(t, tt.want.csrfTokenValidity, validity)

			// csrf token does not get deleted yet
			validity = csrf.ValidToken(context.Background(), migrator.DB, resp.CsrfToken)
			assert.Equal(t, tt.want.csrfTokenValidity, validity)

			time.Sleep(101 * time.Millisecond)

			validity = csrf.ValidToken(context.Background(), migrator.DB, resp.CsrfToken)
			assert.Equal(t, false, validity)
		})
	}
}

func TestHandler_Csrf_Valid_And_Delete_TokenIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	type args struct {
		*LoginRequest
	}
	type want struct {
		error
		status            int
		response          RespondCsrf
		csrfTokenValidity bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "normal",
			args: args{
				LoginRequest: &LoginRequest{
					Email:    "email@example.com",
					Password: "highEntropyPassword",
				},
			},
			want: want{
				error:             nil,
				status:            http.StatusOK,
				response:          RespondCsrf{CsrfToken: ""},
				csrfTokenValidity: true,
			},
		},
	}

	client := dbClient()
	session := newSession(migrator.DB, 10*time.Minute)
	repo := NewRepo(client, migrator.DB, session)

	hashedPassword, err := argon2id.CreateHash("highEntropyPassword", argon2id.DefaultParams)
	assert.Nil(t, err)

	_, err = repo.db.ExecContext(context.Background(), `
		INSERT INTO users (email, password) VALUES ($1, $2)
		ON CONFLICT (email) DO NOTHING 
		`, "email@example.com", hashedPassword)
	assert.Nil(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			var err error
			if tt.args.LoginRequest != nil {
				err = json.NewEncoder(&buf).Encode(tt.args.LoginRequest)
			}
			assert.Nil(t, err)

			rr := httptest.NewRequest(http.MethodPost, "/api/v1/login", &buf)
			ww := httptest.NewRecorder()

			router := chi.NewRouter()
			router.Use(middleware.LoadAndSave(session))

			RegisterHTTPEndPoints(router, session, repo)
			router.ServeHTTP(ww, rr)

			assert.Equal(t, tt.want.status, ww.Code)

			assert.NotNil(t, ww.Header().Get("Set-Cookie"))
			token, err := extractToken(ww.Header().Get("Set-Cookie"))
			assert.Nil(t, err)

			rr = httptest.NewRequest(http.MethodGet, "/api/v1/restricted/csrf", nil)
			ww = httptest.NewRecorder()

			rr.AddCookie(&http.Cookie{
				Name:  sessionName,
				Value: token,
			})

			router.ServeHTTP(ww, rr)

			assert.Equal(t, ww.Code, http.StatusOK)

			b, err := io.ReadAll(ww.Body)
			assert.Nil(t, err)

			var resp RespondCsrf
			err = json.Unmarshal(b, &resp)
			assert.Nil(t, err)

			assert.NotNil(t, resp.CsrfToken)

			err = csrf.ValidAndDeleteToken(context.Background(), migrator.DB, resp.CsrfToken)
			assert.Nil(t, err)

			// at this point, the csrf token would have been deleted
			err = csrf.ValidAndDeleteToken(context.Background(), migrator.DB, resp.CsrfToken)
			assert.NotNil(t, err)
		})
	}
}

func extractToken(cookie string) (string, error) {
	parts := strings.Split(cookie, ";")
	if len(parts) == 0 {
		return "", errors.New("invalid cookie")
	}

	for _, part := range parts {
		keyVal := strings.Split(part, "=")
		if len(keyVal) != 2 {
			return "", errors.New("invalid cookie")
		}
		if keyVal[0] == sessionName {
			return keyVal[1], nil
		}
	}

	return "", errors.New("invalid cookie")
}

func dbClient() *gen.Client {
	drv := entsql.OpenDB(DBDriver, migrator.DB)
	return gen.NewClient(gen.Driver(drv))
}

func newSession(db *sql.DB, duration time.Duration) *scs.SessionManager {
	manager := scs.New()
	manager.Store = postgresstore.New(db)
	manager.CtxStore = postgresstore.New(db)
	manager.Lifetime = duration
	manager.Cookie.Name = sessionName
	manager.Cookie.HttpOnly = false
	manager.Cookie.Path = "/"
	manager.Cookie.SameSite = http.SameSiteLaxMode
	manager.Cookie.Secure = false

	return manager
}
