package sqlorm

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/tinhtinh/common"
	"github.com/tinh-tinh/tinhtinh/core"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Test_Module(t *testing.T) {
	require.NotPanics(t, func() {
		createDatabaseForTest()
	})
	dsn := "host=localhost user=postgres password=postgres dbname=test port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	type User struct {
		Model `gorm:"embedded"`
		Name  string `gorm:"type:varchar(255);not null"`
		Email string `gorm:"type:varchar(255);not null"`
	}
	type UserService struct {
		DB *gorm.DB
	}
	const USER_SERVICE core.Provide = "user_service"

	userService := func(module *core.DynamicModule) *core.DynamicProvider {
		provider := module.NewProvider(core.ProviderOptions{
			Name: USER_SERVICE,
			Factory: func(param ...interface{}) interface{} {
				db := param[0].(*gorm.DB)

				return &UserService{DB: db}
			},
			Inject: []core.Provide{ConnectDB},
		})

		return provider
	}

	userController := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("users")
		service, ok := module.Ref(USER_SERVICE).(*UserService)
		ctrl.Post("", func(ctx core.Ctx) error {
			if !ok {
				return common.InternalServerException(ctx.Res(), "db error")
			}
			result := service.DB.Model(&User{}).Create(&User{Name: "John", Email: "john@gmail.com"})
			if result.Error != nil {
				return common.InternalServerException(ctx.Res(), result.Error.Error())
			}
			return ctx.JSON(core.Map{
				"data": "ok",
			})
		})

		ctrl.Get("", func(ctx core.Ctx) error {
			if !ok {
				return common.InternalServerException(ctx.Res(), "db error")
			}
			var users []User
			result := service.DB.Find(&users)
			if result.Error != nil {
				return common.InternalServerException(ctx.Res(), result.Error.Error())
			}
			return ctx.JSON(core.Map{
				"data": "ok",
			})
		})

		return ctrl
	}

	userModule := func(module *core.DynamicModule) *core.DynamicModule {
		mod := module.New(core.NewModuleOptions{
			Controllers: []core.Controller{userController},
			Providers:   []core.Provider{userService},
			Exports:     []core.Provider{userService},
		})

		return mod
	}

	appModule := func() *core.DynamicModule {
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Module{
				ForRoot(Options{
					Dialect: postgres.Open(dsn),
					Models:  []interface{}{&User{}},
				}),
				userModule,
			},
		})

		return module
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	resp, err := testClient.Post(testServer.URL+"/api/users", "application/json", nil)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/users")
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func createDatabaseForTest() {
	connStr := "host=localhost user=postgres password=postgres port=5432 dbname=postgres sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		panic(err)
	}

	dbName := "test"
	// check if db exists
	stmt := fmt.Sprintf("SELECT * FROM pg_database WHERE datname = '%s';", dbName)
	rs := db.Raw(stmt)
	if rs.Error != nil {
		panic(rs.Error)
	}

	// if not create it
	var rec = make(map[string]interface{})
	if rs.Find(rec); len(rec) == 0 {
		stmt := fmt.Sprintf("CREATE DATABASE %s;", dbName)
		if rs := db.Exec(stmt); rs.Error != nil {
			panic(rs.Error)
		}

		// close db connection
		sql, err := db.DB()
		defer func() {
			_ = sql.Close()
		}()
		if err != nil {
			panic(err)
		}
	}
}
