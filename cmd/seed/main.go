package main

import (
	"fmt"
	"github.com/gmhafiz/go8/config"
	"github.com/gmhafiz/go8/database"
	db "github.com/gmhafiz/go8/third_party/database"
)

func main() {
	cfg := config.New()
	store := db.NewSqlx(cfg.Database)

	seeder := database.Seeder(store.DB)
	seeder.SeedUsers()
	fmt.Println("seeding completed.")
}
