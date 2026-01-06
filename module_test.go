package sqlorm_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/sqlorm/v2"
	"github.com/tinh-tinh/tinhtinh/v2/common"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Test_Module(t *testing.T) {
	require.NotPanics(t, func() {
		createDatabaseForTest("test")
	})
	dsn := "host=localhost user=postgres password=postgres dbname=test port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	type User struct {
		gorm.Model
		Name  string `gorm:"type:varchar(255);not null"`
		Email string `gorm:"type:varchar(255);not null"`
	}

	userController := func(module core.Module) core.Controller {
		ctrl := module.NewController("users")
		repo := sqlorm.InjectRepository[User](module)

		ctrl.Post("", func(ctx core.Ctx) error {
			result, err := repo.Create(&User{Name: "John", Email: "john@gmail.com"})
			if err != nil {
				return common.InternalServerException(ctx.Res(), err.Error())
			}
			return ctx.JSON(core.Map{
				"data": result,
			})
		})

		ctrl.Get("", func(ctx core.Ctx) error {
			result, err := repo.FindAll(nil, sqlorm.FindOptions{})
			if err != nil {
				return common.InternalServerException(ctx.Res(), err.Error())
			}
			return ctx.JSON(core.Map{
				"data": result,
			})
		})

		return ctrl
	}

	userModule := func(module core.Module) core.Module {
		mod := module.New(core.NewModuleOptions{
			Imports:     []core.Modules{sqlorm.ForFeature(sqlorm.NewRepo(User{}))},
			Controllers: []core.Controllers{userController},
		})

		return mod
	}

	appModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{
				sqlorm.ForRoot(sqlorm.Config{
					Dialect: postgres.Open(dsn),
					Models:  []interface{}{&User{}},
				}),
				userModule,
			},
		})

		return module
	}

	connect := sqlorm.Inject(appModule())
	require.NotNil(t, connect)

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	resp, err := testClient.Post(testServer.URL+"/api/users", "application/json", nil)
	require.Nil(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/users")
	require.Nil(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func Test_ModuleFactory(t *testing.T) {
	require.NotPanics(t, func() {
		createDatabaseForTest("test2")
	})
	dsn := "host=localhost user=postgres password=postgres dbname=test2 port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	type User struct {
		gorm.Model
		Name  string `gorm:"type:varchar(255);not null"`
		Email string `gorm:"type:varchar(255);not null"`
	}

	userController := func(module core.Module) core.Controller {
		ctrl := module.NewController("users")
		repo := sqlorm.InjectRepository[User](module)

		ctrl.Post("", func(ctx core.Ctx) error {
			result, err := repo.Create(&User{Name: "John", Email: "john@gmail.com"})
			if err != nil {
				return common.InternalServerException(ctx.Res(), err.Error())
			}
			return ctx.JSON(core.Map{
				"data": result,
			})
		})

		ctrl.Get("", func(ctx core.Ctx) error {
			result, err := repo.FindAll(nil, sqlorm.FindOptions{})
			if err != nil {
				return common.InternalServerException(ctx.Res(), err.Error())
			}
			return ctx.JSON(core.Map{
				"data": result,
			})
		})

		return ctrl
	}

	userModule := func(module core.Module) core.Module {
		mod := module.New(core.NewModuleOptions{
			Imports:     []core.Modules{sqlorm.ForFeature(sqlorm.NewRepo(User{}))},
			Controllers: []core.Controllers{userController},
		})

		return mod
	}

	appModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{
				sqlorm.ForRootFactory(func(module core.RefProvider) sqlorm.Config {
					return sqlorm.Config{
						Dialect: postgres.Open(dsn),
						Models:  []interface{}{User{}},
					}
				}),
				userModule,
			},
		})

		return module
	}

	connect := sqlorm.Inject(appModule())
	require.NotNil(t, connect)

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("/api")

	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	resp, err := testClient.Post(testServer.URL+"/api/users", "application/json", nil)
	require.Nil(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	resp, err = testClient.Get(testServer.URL + "/api/users")
	require.Nil(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func Test_ModuleSync(t *testing.T) {
	require.NotPanics(t, func() {
		createDatabaseForTest("test")
	})
	dsn := "host=localhost user=postgres password=postgres dbname=test port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	type User struct {
		gorm.Model
		Name  string `gorm:"type:varchar(255);not null"`
		Email string `gorm:"type:varchar(255);not null"`
	}

	appModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{
				sqlorm.ForRoot(sqlorm.Config{
					Dialect: postgres.Open(dsn),
					Models:  []interface{}{&User{}},
				}),
			},
		})

		return module
	}

	connect := sqlorm.Inject(appModule())
	require.NotNil(t, connect)
}

func Test_Panic(t *testing.T) {
	require.Panics(t, func() {
		appModule := func() core.Module {
			module := core.NewModule(core.NewModuleOptions{
				Imports: []core.Modules{
					sqlorm.ForRoot(sqlorm.Config{
						Dialect: postgres.Open(""),
					}),
				},
			})

			return module
		}

		connect := sqlorm.Inject(appModule())
		require.NotNil(t, connect)

		app := core.CreateFactory(appModule)
		app.SetGlobalPrefix("/api")

		testServer := httptest.NewServer(app.PrepareBeforeListen())
		defer testServer.Close()
	})
}

func Test_Nil(t *testing.T) {

	type User struct {
		gorm.Model
		Name  string `gorm:"type:varchar(255);not null"`
		Email string `gorm:"type:varchar(255);not null"`
	}

	require.Nil(t, sqlorm.InjectRepository[User](core.NewModule(core.NewModuleOptions{})))
}

func Test_RefName(t *testing.T) {
	require.NotPanics(t, func() {
		createDatabaseForTest("test")
	})
	dsn := "host=localhost user=postgres password=postgres dbname=test port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	appModule := func() core.Module {
		module := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{
				sqlorm.ForRoot(sqlorm.Config{
					Dialect: postgres.Open(dsn),
					Models:  []interface{}{&Abc{}},
				}),
				sqlorm.ForFeature(sqlorm.NewRepo(Abc{})),
			},
		})

		return module
	}
	abcRepo := sqlorm.InjectRepository[Abc](appModule())
	require.NotNil(t, abcRepo)
}

func createDatabaseForTest(dbName string) {
	connStr := "host=localhost user=postgres password=postgres port=5432 dbname=postgres sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		panic(err)
	}

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
