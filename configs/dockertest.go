package configs

import "os"

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

func DockerTestCfg() *DockerTest {
	return &DockerTest{
		Driver:  os.Getenv("DOCKERTEST_DRIVER"),
		Dialect: os.Getenv("DOCKERTEST_DIALECT"),
		Host:    os.Getenv("DOCKERTEST_HOST"),
		Port:    os.Getenv("DOCKERTEST_PORT"),
		User:    os.Getenv("DOCKERTEST_USER"),
		Name:    os.Getenv("DOCKERTEST_NAME"),
		Pass:    os.Getenv("DOCKERTEST_PASS"),
		SslMode: os.Getenv("DOCKERTEST_SSL_MODE"),
	}
}
