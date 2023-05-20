package database

import (
	"database/sql"
	"embed"
	"log"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type Migrate struct {
	DB  *sql.DB
	dsn string
}

type Options func(opts *Migrate) error

func Migrator(opts ...Options) *Migrate {

	m := &Migrate{}
	goose.SetBaseFS(embedMigrations)

	for _, opt := range opts {
		err := opt(m)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return m
}

func (m *Migrate) Up() {
	if err := goose.Up(m.DB, "migrations"); err != nil {
		panic(err)
	}
}

func (m *Migrate) Down() {
	if err := goose.Down(m.DB, "migrations"); err != nil {
		panic(err)
	}
}

func WithDSN(dsn string) func(opts *Migrate) error {
	return func(opts *Migrate) error {
		opts.dsn = dsn
		return nil
	}
}

func WithDB(db *sql.DB) func(opts *Migrate) error {
	return func(opts *Migrate) error {
		opts.DB = db
		return nil
	}
}
