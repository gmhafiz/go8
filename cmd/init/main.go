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

	"github.com/gmhafiz/go8/configs"
)

const (
	GolangMigrateVersion = "v4.11.0"
	GolangCIVersion      = "v1.31.0"
)

var (
	HomeDir string
)

func init() {
	HomeDir = os.Getenv("HOME")
}

func main() {
	var dbConfigs configs.Database

	go syncsGoMod()

	installTools()

	copyENV()

	cfg := getDbDetails(&dbConfigs)

	fillIntoENV(cfg)

	exportEnv(cfg)
}

func exportEnv(cfg *configs.Database) {
	postgresqlUrl := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s", cfg.Driver, cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name, cfg.SslMode)
	export := fmt.Sprintf("export POSTGRESQL_URL=%s", postgresqlUrl)
	_ = exec.Command(export)
}

func installTools() {
	installGolangMigrate()
	installTestify()
	installSwag()
	installGolangCILint()
}

func installGolangCILint() {
	cmd := fmt.Sprintf("curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin %s", GolangCIVersion)
	_, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		log.Fatalln("error installing")
	}
}

func installSwag() {
	if binaryExists("swag") {
		return
	}

	goGet("github.com/swaggo/swag/cmd/swag")
}

func installTestify() {
	if binaryExists("testify") {
		return
	}

	goGet("github.com/stretchr/testify")
}

func goGet(path string) {
	log.Printf("installing %s\n", path)
	cmd := exec.Command("go", "get", "-u", "-t", path)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "GO111MODULE=off")
	err := cmd.Run()
	if err != nil {
		log.Fatalln("error downloading")
	}
}

func installGolangMigrate() {
	if binaryExists("migrate") {
		return
	}

	log.Println("installing golang-migrate")
	cmd := exec.Command("wget", "https://github.com/golang-migrate/migrate/releases/download/"+GolangMigrateVersion+"/migrate.linux-amd64.tar.gz")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("golang-migrate download failed with %s \n", err)
	}

	cmd = exec.Command("tar", "xvzf", "migrate.linux-amd64.tar.gz")
	err = cmd.Run()
	if err != nil {
		log.Fatalf("extracting golang-migrate failed with %s \n", err)
	}

	cmd = exec.Command("mkdir", "-p", fmt.Sprintf("%s/.local/bin", HomeDir))
	err = cmd.Run()
	if err != nil {
		log.Fatalf("error creating folder with %s \n", err)
	}

	cmd = exec.Command("mv", "migrate.linux-amd64", fmt.Sprintf("%s/.local/bin/migrate", HomeDir))
	err = cmd.Run()
	if err != nil {
		log.Fatalf("error moving binary with %s \n", err)
	}

	return
}

func binaryExists(binaryName string) bool {
	cmd := exec.Command("which", binaryName)
	stdout, _ := cmd.CombinedOutput()
	if len(stdout) != 0 {
		return true
	}
	return false
}

func syncsGoMod() {
	cmd := exec.Command("go", "mod", "tidy")
	err := cmd.Run()
	if err != nil {
		log.Fatalln("error executing go mod tidy")
	}
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
	cfg.Driver = getInput("Enter database driver (postgres/mysql) (default: postgres): ", "postgres")
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
