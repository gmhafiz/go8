package app

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"

	"eight/config"
	"eight/database"
)

type Server struct {
	config *config.Conf
	db     *gorm.DB
}

func NewApp(c *config.Conf) *Server {
	db, err := database.NewDatabase(c)
	if err != nil {
		log.Fatalln(err)
	}
	return &Server{
		config: c,
		db: db,
	}
}

func (s *Server) GetConfig() *config.Conf {
	return  s.config
}