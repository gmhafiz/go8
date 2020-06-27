package app

import (
	"github.com/mattn/go-sqlite3"
)

type Server struct {
	db     *sqlite3.SQLiteConn
}

func NewApp() *Server {
	return &Server{
		db: nil,
	}
}