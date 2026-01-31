package main

import (
	"fmt"

	"codeberg.org/gmhafiz/go8/config"
	"codeberg.org/gmhafiz/go8/database"

	db "codeberg.org/gmhafiz/go8/third_party/database"
)

func main() {
	cfg := config.New()
	store := db.NewSqlx(cfg.Database)

	seeder := database.Seeder(store.DB)
	seeder.SeedUsers()
	fmt.Println("seeding completed.")
}
