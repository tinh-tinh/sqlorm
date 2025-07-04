package sqlorm_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/sqlorm/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestRetryConnect(t *testing.T) {
	require.Panics(t, func() {
		dsn := "host=localhost user=postgres password=postgres dbname=xoxinh port=5432 sslmode=disable TimeZone=Asia/Shanghai"
		conn := sqlorm.NewConnect(sqlorm.Config{
			Dialect: postgres.Open(dsn),
			Retry: &sqlorm.RetryOptions{
				MaxRetries: 3,
				Delay:      time.Second,
			},
		})
		require.NotNil(t, conn)
	})
}

func TestOnInit(t *testing.T) {
	require.NotPanics(t, func() {
		createDatabaseForTest("test")
	})
	dsn := "host=localhost user=postgres password=postgres dbname=test port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	conn := sqlorm.NewConnect(sqlorm.Config{
		Dialect: postgres.Open(dsn),
		OnInit: func(db *gorm.DB) {
			require.NotNil(t, db)
		},
	})

	require.NotNil(t, conn)
}
