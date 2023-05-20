package postgresstore

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/gmhafiz/go8/config"
	"github.com/gmhafiz/go8/internal/middleware"
)

func TestMain(m *testing.M) {
	getwd, err := os.Getwd()
	if err != nil {
		return
	}
	if err != nil {
		log.Fatalln(err)
	}
	if strings.Contains(getwd, "/third_party/postgresstore") {
		err := os.Chdir("../../")
		if err != nil {
			log.Fatalln(err)
		}
	}

	cfg := config.New()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Pass,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.TestName,
	)

	err = os.Setenv("SCS_POSTGRES_TEST_DSN", dsn)
	if err != nil {
		log.Fatalln(err)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	_, err = db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS users
(
    id          bigint generated always as identity primary key ,
    first_name  text,
    middle_name text,
    last_name   text,
    email       text unique,
    password    text,
    verified_at timestamptz
);`)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS sessions
(
    token   TEXT PRIMARY KEY,
    user_id BIGINT      NOT NULL CONSTRAINT session_user_fk REFERENCES users ON DELETE CASCADE ,
    data    BYTEA       NOT NULL,
    expiry  TIMESTAMPTZ NOT NULL
);`)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = db.ExecContext(ctx, `
INSERT INTO users (email, password, verified_at) 
VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING 
`, "admin@example.com", "test", time.Now(),
	)
	if err != nil {
		log.Fatalln(err)
	}

	code := m.Run()
	os.Exit(code)
}

func TestFind(t *testing.T) {
	dsn := os.Getenv("SCS_POSTGRES_TEST_DSN")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("TRUNCATE TABLE sessions")
	if err != nil {
		t.Fatal(err)
	}

	hash, err := sum("session_token")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
INSERT INTO sessions
VALUES($1, $2, $3, current_timestamp + interval '1 minute')`, hash, 1, "encoded_data")
	if err != nil {
		t.Fatal(err)
	}

	p := NewWithCleanupInterval(db, 0)

	ctx := context.Background()

	b, found, err := p.FindCtx(ctx, "session_token")
	if err != nil {
		t.Fatal(err)
	}
	if found != true {
		t.Fatalf("got %v: expected %v", found, true)
	}
	if bytes.Equal(b, []byte("encoded_data")) == false {
		t.Fatalf("got %v: expected %v", b, []byte("encoded_data"))
	}
}

func TestFindMissing(t *testing.T) {
	dsn := os.Getenv("SCS_POSTGRES_TEST_DSN")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("TRUNCATE TABLE sessions")
	if err != nil {
		t.Fatal(err)
	}

	p := NewWithCleanupInterval(db, 0)

	ctx := context.Background()

	_, found, err := p.FindCtx(ctx, "missing_session_token")
	if err != nil {
		t.Fatalf("got %v: expected %v", err, nil)
	}
	if found != false {
		t.Fatalf("got %v: expected %v", found, false)
	}
}

func TestSaveNew(t *testing.T) {
	dsn := os.Getenv("SCS_POSTGRES_TEST_DSN")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("TRUNCATE TABLE sessions")
	if err != nil {
		t.Fatal(err)
	}

	p := NewWithCleanupInterval(db, 0)

	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.KeyID, uint64(1))

	err = p.CommitCtx(ctx, "session_token", []byte("encoded_data"), time.Now().Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}

	hash, err := sum("session_token")
	if err != nil {
		t.Fatal(err)
	}

	row := db.QueryRow("SELECT data FROM sessions WHERE token = $1", hash)
	var data []byte
	err = row.Scan(&data)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(data, []byte("encoded_data")) == false {
		t.Fatalf("got %v: expected %v", data, []byte("encoded_data"))
	}
}

func TestSaveUpdated(t *testing.T) {
	dsn := os.Getenv("SCS_POSTGRES_TEST_DSN")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("TRUNCATE TABLE sessions")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO sessions VALUES('session_token', 1, 'encoded_data', current_timestamp + interval '1 minute')")
	if err != nil {
		t.Fatal(err)
	}

	p := NewWithCleanupInterval(db, 0)

	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.KeyID, uint64(1))

	err = p.CommitCtx(ctx, "session_token", []byte("new_encoded_data"), time.Now().Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}

	hash, err := sum("session_token")
	if err != nil {
		t.Fatal(err)
	}

	row := db.QueryRow("SELECT data FROM sessions WHERE token = $1", hash)
	var data []byte
	err = row.Scan(&data)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(data, []byte("new_encoded_data")) == false {
		t.Fatalf("got %v: expected %v", data, []byte("new_encoded_data"))
	}
}

func TestExpiry(t *testing.T) {
	dsn := os.Getenv("SCS_POSTGRES_TEST_DSN")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("TRUNCATE TABLE sessions")
	if err != nil {
		t.Fatal(err)
	}

	p := NewWithCleanupInterval(db, 0)

	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.KeyID, uint64(1))

	err = p.CommitCtx(ctx, "session_token", []byte("encoded_data"), time.Now().Add(100*time.Millisecond))
	if err != nil {
		t.Fatal(err)
	}

	_, found, _ := p.FindCtx(ctx, "session_token")
	if found != true {
		t.Fatalf("got %v: expected %v", found, true)
	}

	time.Sleep(100 * time.Millisecond)
	_, found, _ = p.FindCtx(ctx, "session_token")
	if found != false {
		t.Fatalf("got %v: expected %v", found, false)
	}
}

func TestDelete(t *testing.T) {
	dsn := os.Getenv("SCS_POSTGRES_TEST_DSN")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("TRUNCATE TABLE sessions")
	if err != nil {
		t.Fatal(err)
	}

	hash, err := sum("session_token")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO sessions VALUES($1, 1, 'encoded_data', current_timestamp + interval '1 minute')", hash)
	if err != nil {
		t.Fatal(err)
	}

	p := NewWithCleanupInterval(db, 0)

	ctx := context.Background()

	err = p.DeleteCtx(ctx, "session_token")
	if err != nil {
		t.Fatal(err)
	}

	row := db.QueryRow("SELECT COUNT(*) FROM sessions WHERE token = 'session_token'")
	var count int
	err = row.Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("got %d: expected %d", count, 0)
	}
}

func TestCleanup(t *testing.T) {
	dsn := os.Getenv("SCS_POSTGRES_TEST_DSN")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("TRUNCATE TABLE sessions")
	if err != nil {
		t.Fatal(err)
	}

	p := NewWithCleanupInterval(db, 200*time.Millisecond)
	defer p.StopCleanup()

	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.KeyID, uint64(1))

	hash, err := sum("session_token")
	if err != nil {
		t.Fatal(err)
	}

	err = p.CommitCtx(ctx, "session_token", []byte("encoded_data"), time.Now().Add(100*time.Millisecond))
	if err != nil {
		t.Fatal(err)
	}

	row := db.QueryRow("SELECT COUNT(*) FROM sessions WHERE token = $1", hash)
	var count int
	err = row.Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("got %d: expected %d", count, 1)
	}

	time.Sleep(300 * time.Millisecond)
	row = db.QueryRow("SELECT COUNT(*) FROM sessions WHERE token = $1", hash)
	err = row.Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatalf("got %d: expected %d", count, 0)
	}
}

func TestStopNilCleanup(t *testing.T) {
	dsn := os.Getenv("SCS_POSTGRES_TEST_DSN")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		t.Fatal(err)
	}

	p := NewWithCleanupInterval(db, 0)
	time.Sleep(100 * time.Millisecond)
	// A send to a nil channel will block forever
	p.StopCleanup()
}
