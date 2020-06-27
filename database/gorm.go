package database

import (
	"fmt"

	"github.com/jinzhu/gorm"

	"eight/config"
)

func NewDatabase(conf *config.Conf) (*gorm.DB, error) {
	var connStr string
	if conf.Testing {
		connStr = fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=%v",
			conf.Db.TestHost,
			conf.Db.TestPort,
			conf.Db.TestUsername,
			conf.Db.TestDbName,
			conf.Db.TestPassword,
			conf.Db.TestSslMode,
		)
	} else {
		connStr = fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=%v",
			conf.Db.Host,
			conf.Db.Port,
			conf.Db.Username,
			conf.Db.DbName,
			conf.Db.Password,
			conf.Db.SslMode,
		)
	}

	return gorm.Open("postgres", connStr)
}
