package configs

import "os"

type Database struct {
	Driver  string
	Host    string
	Port    string
	Name    string
	User    string
	Pass    string
	SslMode string
}

type DockerTest struct {
	Driver  string
	Dialect string
	Host    string
	Port    string
	Name    string
	User    string
	Pass    string
	SslMode string
}

func DataStore() *Database {
	return &Database{
		Driver:  os.Getenv("DB_DRIVER"),
		Host:    os.Getenv("DB_HOST"),
		Port:    os.Getenv("DB_PORT"),
		Name:    os.Getenv("DB_NAME"),
		User:    os.Getenv("DB_USER"),
		Pass:    os.Getenv("DB_PASS"),
		SslMode: os.Getenv("DB_SSL_MODE"),
	}
}

func DockerTestCfg() *DockerTest {
	return &DockerTest{
		Driver:  os.Getenv("DOCKERTEST_DRIVER"),
		Dialect: os.Getenv("DOCKERTEST_DIALECT"),
		Host:    os.Getenv("DOCKERTEST_HOST"),
		Port:    os.Getenv("DOCKERTEST_PORT"),
		User:    os.Getenv("DOCKERTEST_USER"),
		Pass:    os.Getenv("DOCKERTEST_PASS"),
		SslMode: os.Getenv("DOCKERTEST_SSL_MODE"),
	}
}
