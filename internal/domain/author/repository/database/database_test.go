package database

import (
	"testing"

	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

//go:generate mockgen -package mock -source ../../repository.go -destination=../../mock/mock_repository.go

func TestAuthorRepository_Create(t *testing.T) {}

func TestRepository_Find(t *testing.T) {}

func TestRepository_List(t *testing.T) {}

func TestRepository_Update(t *testing.T) {}

func TestRepository_Delete(t *testing.T) {}

func TestRepository_Search(t *testing.T) {}
