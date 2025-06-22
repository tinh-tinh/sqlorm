package sqlorm

import (
	"fmt"
	"time"

	"github.com/tinh-tinh/tinhtinh/v2/common/color"
	"gorm.io/gorm"
)

func NewConnect(config Config) *gorm.DB {
	conn, err := gorm.Open(config.Dialect, config.Options...)
	if err != nil {
		if config.Retry != nil && config.Retry.MaxRetries > 0 {
			fmt.Printf("%s %s %s %s\n",
				color.Green("[SQLORM]"),
				color.White("Failed to connect to database:"),
				color.Red(err.Error()),
				color.Yellow(fmt.Sprintf("Retrying attempt remain %d", config.Retry.MaxRetries)),
			)
			time.Sleep(config.Retry.Delay)
			config.Retry.MaxRetries--
			return NewConnect(config)
		}
		panic(err)
	}
	if config.OnInit != nil {
		config.OnInit(conn)
	}
	if config.Sync {
		err = conn.AutoMigrate(config.Models...)
		if err != nil {
			panic(err)
		}
	}
	return conn
}
