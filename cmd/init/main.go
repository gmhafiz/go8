package main

// Simplifies initial setup before starting the kit
// 1. Syncs dependencies
// 2. Copies `example` files to its config files
// 3. Fills required environment variables
// 4. Downloads and installs various tools to work with this kit

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/ini.v1"

	"go8ddd/configs"
)

const (
	GolangMigrateVersion = "v4.11.0"
	GolangCIVersion      = "v1.31.0"
)

func main() {
	var dbConfigs configs.Database

	go syncsGoMod()

	go installTools()

	copyENV()

	copySqlBoiler()

	cfg := getDbDetails(&dbConfigs)

	fillIntoENV(cfg)

	fillIntoSQLBoiler(cfg)

	exportEnv(cfg)
}

func exportEnv(cfg *configs.Database) {
	postgresqlUrl := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s", cfg.Driver, cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name, cfg.SslMode)
	export := fmt.Sprintf("export POSTGRESQL_URL=%s", postgresqlUrl)
	_ = exec.Command(export)
}

func installTools() {
	_ = exec.Command("curl", "-L",
		"https://github.com/golang-migrate/migrate/releases/download/"+GolangMigrateVersion+"/migrate.linux-amd64.tar.gz", "|", "tar", "xvz")
	_ = exec.Command("mv migrate.linux-amd64 ~/.local/bin/migrate")
	_ = exec.Command("source ~/.bashrc")
	_ = exec.Command("GO111MODULE=off go get -u -t github.com/volatiletech/sqlboiler")
	_ = exec.Command("GO111MODULE=off go get github.com/volatiletech/sqlboiler/drivers/sqlboiler-psql")
	_ = exec.Command("GO111MODULE=off go get github.com/volatiletech/sqlboiler/drivers/sqlboiler-mysql")
	_ = exec.Command("GO111MODULE=off go get github.com/volatiletech/sqlboiler/drivers/sqlboiler-mssql")
	_ = exec.Command("go get github.com/stretchr/testify")
	_ = exec.Command("go get -u github.com/swaggo/swag/cmd/swag")
	_ = exec.Command("go get gopkg.in/ini.v1")
	_ = exec.Command("go get github.com/go-redis/redis/v8")
	_ = exec.Command("curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin " + GolangCIVersion)
}

func syncsGoMod() {
	_ = exec.Command("go", "mod", "tidy")
}

func fillIntoENV(cfg *configs.Database) {
	file, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal(err)
	}

	file["DB_DRIVER"] = cfg.Driver
	file["DB_HOST"] = cfg.Host
	file["DB_PORT"] = cfg.Port
	file["DB_NAME"] = cfg.Name
	file["DB_USER"] = cfg.User
	file["DB_PASS"] = cfg.Pass
	file["DB_SSL_MODE"] = cfg.SslMode

	err = godotenv.Write(file, ".env")
	if err != nil {
		log.Fatal(err)
	}
}

func fillIntoSQLBoiler(dbConfigs *configs.Database) {
	cfg, err := ini.Load("sqlboiler.toml")
	if err != nil {
		log.Fatal(err)
	}

	var section string
	if cfg.Section("").Key("app_mode").String() == "postgres" {
		section = "psql"
	}

	cfg.Section(section).Key("dbname").SetValue(dbConfigs.Name)
	cfg.Section(section).Key("host").SetValue(dbConfigs.Host)
	cfg.Section(section).Key("port").SetValue(dbConfigs.Port)
	cfg.Section(section).Key("user").SetValue(dbConfigs.User)
	cfg.Section(section).Key("pass").SetValue(dbConfigs.Pass)
	cfg.Section(section).Key("sslmode").SetValue(dbConfigs.SslMode)

	err = cfg.SaveTo("sqlboiler.toml")
	if err != nil {
		log.Fatal(err)
	}
}

func getDbDetails(cfg *configs.Database) *configs.Database {
	cfg.Driver = getInput("Enter database driver (postgres/mysql) (default: postgres): ",
		"postgres")
	cfg.Host = getInput("Enter database host (default: localhost): ", "localhost")
	cfg.Port = getInput("Enter database port (default: 5432): ", "5432")
	cfg.Name = getInput("Enter database name: ", "")
	cfg.User = getInput("Enter database username: ", "")
	cfg.Pass = getInput("Enter database password: ", "")
	cfg.SslMode = getInput("Database SSL MODE (disable/enable) (default: disable): ", "disable")

	return cfg
}

func getInput(prompt, defaultValue string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Println("error reading input")
		os.Exit(1)
	}

	var userInput string
	if input == "\n" {
		if defaultValue == "" {
			fmt.Println(" cannot be empty. Please try again")
			userInput = getInput(prompt, "")
		} else {
			userInput = defaultValue
		}
	} else {
		userInput = strings.Trim(input, "\n")
	}

	return userInput
}

func copySqlBoiler() {
	if fileExists("sqlboiler.toml") {
		fmt.Println("sqlboiler.toml file already exists. Run 'task configs' to update new values")
	} else {
		_, err := copyFile("sqlboiler.toml.example", "sqlboiler.toml")
		if err != nil {
			log.Println("error copying to sqlboiler.toml")
			os.Exit(1)
		}
	}
}

func copyENV() {
	if fileExists(".env") {
		fmt.Println(".env file already exists. Run 'task configs' to update new values")
	} else {
		_, err := copyFile("env.example", ".env")
		if err != nil {
			log.Println("error copying to .env")
			os.Exit(1)
		}
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
