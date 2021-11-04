package database

import (
	"testing"

	_ "github.com/jackc/pgx/stdlib"
)

//go:generate mockgen -package mock -source ./postgres.go -destination=../../mock/mock_postgres.go

func TestAuthorRepository_Create(t *testing.T) {}

func TestRepository_Find(t *testing.T) {}

func TestRepository_List(t *testing.T) {}

func TestRepository_Update(t *testing.T) {}

func TestRepository_Delete(t *testing.T) {}

func TestRepository_Search(t *testing.T) {}
