package configs

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"gopkg.in/yaml.v2"

	"eight/internal/platform/datastore"
	"eight/internal/server/http"
)

type Configs struct {
	Api      http.Config        `yaml:"Api"`
	Database datastore.Database `yaml:"Database"`
}

// HTTP returns the configuration required for HTTP package
func (cfg *Configs) HTTP() (*http.Config, error) {
	return &http.Config{
		Port:         cfg.Api.Port,
		ReadTimeout:  time.Second * 9995,
		WriteTimeout: time.Second * 9995,
		DialTimeout:  time.Second * 9993,
	}, nil
}

func (cfg *Configs) DataStore() (*datastore.Database, error) {
	return &datastore.Database{
		Driver:  cfg.Database.Driver,
		Host:    cfg.Database.Host,
		Port:    cfg.Database.Port,
		Name:    cfg.Database.Name,
		User:    cfg.Database.User,
		Pass:    cfg.Database.Pass,
		SslMode: cfg.Database.SslMode,
	}, nil
}

func NewService(mode string) (*Configs, error) {
	wd, _ := os.Getwd()
	fileName := path.Join(wd, "config", mode) + ".yml"
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Panic(err)
	}

	cfg := &Configs{}
	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		log.Panic(err)
	}

	return cfg, nil
}
