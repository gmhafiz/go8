package main

import (
	"eight/app"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pressly/goose"

	"eight/config"
)

const dialect = "postgres"

var (
	flags = flag.NewFlagSet("migrate", flag.ExitOnError)
	dir   = flags.String("dir", "./database/migrations", "directory with migrations files")
)

func main() {
	flags.Usage = usage
	_ = flags.Parse(os.Args[1:])

	args := flags.Args()
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		flags.Usage()
		return
	}

	command := args[0]
	switch command {
	case "create":
		if err := goose.Run("create", nil, *dir, args[1:]...); err != nil {
			log.Fatalf("migrate run: %v", err)
		}
		return
	case "fix":
		if err := goose.Run("fix", nil, *dir); err != nil {
			log.Fatalf("migrate run: %v", err)
		}
		return
	}

	appConfig := config.AppConfig()
	application := app.NewApp(appConfig)
	defer application.GetDB().Close()

	if err := goose.SetDialect(dialect); err != nil {
		log.Fatal(err)
	}
	cwd, err := os.Getwd()
	log.Println(cwd)
	if err != nil {
		log.Fatalln(err)
	}

	err = application.GetDB().DB().Ping()
	application.GetDB()
	if err != nil {
		log.Fatalln(err)
	}

	if err := goose.Run(command, application.GetDB().DB(), *dir, args[1:]...); err != nil {
		log.Fatalf("migrate run: %v", err)
	}

	appConfig.Testing = true
	application = app.NewApp(appConfig)
	defer application.GetDB().Close()
	if err := goose.Run(command, application.GetDB().DB(), *dir, args[1:]...); err != nil {
		log.Fatalf("migrate run: %v", err)
	}
}

func usage() {
	fmt.Println(usagePrefix)
	flags.PrintDefaults()
	fmt.Println(usageCommands)
}

var (
	usagePrefix = `Usage: migrate [OPTIONS] COMMAND
Examples:
    migrate status
Options:
`

	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migrations
    reset                Roll back all migrations
    status               Dump the migrations status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migrations file with the current timestamp
    fix                  Apply sequential ordering to migrations
`
)
