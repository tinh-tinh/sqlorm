package sqlorm

import (
	"gorm.io/gorm"
)

func NewConnect(config Config) *gorm.DB {
	conn, err := gorm.Open(config.Dialect, config.Options...)
	if err != nil {
		panic(err)
	}
	conn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if config.Sync {
		err = conn.AutoMigrate(config.Models...)
		if err != nil {
			panic(err)
		}
	}
	return conn
}
