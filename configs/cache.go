package configs

import (
	"log"
	"os"
	"strconv"
)

type Cache struct {
	Host      string
	Port      string
	Name      int
	User      string
	Pass      string
	CacheTime int
}

func NewCache() *Cache {
	name, err := strconv.Atoi(os.Getenv("REDIS_NAME"))
	if err != nil {
		log.Fatal(err)
	}

	cacheTime, err := strconv.Atoi(os.Getenv("REDIS_CACHE_TIME"))
	if err != nil {
		log.Fatal(err)
	}

	return &Cache{
		Host:      os.Getenv("REDIS_HOST"),
		Port:      os.Getenv("REDIS_PORT"),
		Name:      name,
		User:      os.Getenv("REDIS_USER"),
		Pass:      os.Getenv("REDIS_PASS"),
		CacheTime: cacheTime,
	}
}
