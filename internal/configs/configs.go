package configs

import (
	"eight/pkg/redis"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"time"

	"github.com/jinzhu/now"

	"eight/internal/datastore"
	"eight/internal/server/http"
)

type Configs struct {
	Api      http.Config        `yaml:"Api"`
	Database datastore.Database `yaml:"Database"`
	Cache    redis.Config       `yaml:"Redis"`
	Time     now.Config         `yaml:"Time"`
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

func (cfg *Configs) CacheStore() (*redis.Config, error) {
	return &redis.Config{
		Host:     cfg.Cache.Host,
		Port:     cfg.Cache.Port,
		Name:     cfg.Cache.Name,
		Username: cfg.Cache.Username,
		Password: cfg.Cache.Password,
	}, nil
}

func NewService(file string) (*Configs, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	cfg := &Configs{}
	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}

	timeConfig, err := timeConfigInit()
	if err != nil {
		return nil, err
	}
	cfg.Time = *timeConfig

	return cfg, nil
}

//func NewService(mode string) (*Configs, error) {
//	wd, _ := os.Getwd()
//	fileName := path.Join(wd, "config", mode) + ".yml"
//	content, err := ioutil.ReadFile(fileName)
//	if err != nil {
//		log.Panic(err)
//	}
//
//	cfg := &Configs{}
//	err = yaml.Unmarshal(content, &cfg)
//	if err != nil {
//		log.Panic(err)
//	}
//
//	timeConfig, err := timeConfigInit()
//	if err != nil {
//		log.Panic(err)
//	}
//	cfg.Time = *timeConfig
//
//	boil.DebugMode = true
//
//	return cfg, nil
//}

func timeConfigInit() (*now.Config, error) {
	location, err := time.LoadLocation("Australia/Sydney")
	if err != nil {
		return nil, err
	}
	myConfig := &now.Config{
		WeekStartDay: time.Sunday,
		TimeLocation: location,
		TimeFormats:  []string{"2006-01-02 15:04:05.999999999"},
	}

	return myConfig, nil
}
