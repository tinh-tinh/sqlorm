package tenancy_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/sqlorm/v2"
	"github.com/tinh-tinh/tenancy"
	"github.com/tinh-tinh/tinhtinh/v2/common"
	"github.com/tinh-tinh/tinhtinh/v2/core"
	"gorm.io/gorm"
)

func Test_Tenant(t *testing.T) {
	type User struct {
		gorm.Model
		Name  string `gorm:"type:varchar(255);not null"`
		Email string `gorm:"type:varchar(255);not null"`
	}

	userController := func(module core.Module) core.Controller {
		ctrl := module.NewController("users")
		ctrl.Post("", func(ctx core.Ctx) error {
			repo := tenancy.InjectRepository[User](module, ctx)
			result, err := repo.Create(&User{Name: "John", Email: "john@gmail.com"})
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
		userMod := module.New(core.NewModuleOptions{
			Imports:     []core.Modules{tenancy.ForFeature(sqlorm.NewRepo(User{}))},
			Controllers: []core.Controllers{userController},
		})
		return userMod
	}

	appModule := func() core.Module {
		appModule := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{
				tenancy.ForRoot(tenancy.Options{
					Connect: tenancy.ConnectOptions{
						Host:     "localhost",
						Port:     5432,
						User:     "postgres",
						Password: "postgres",
					},
					GetTenantID: func(r *http.Request) string {
						return r.Header.Get("x-tenant-id")
					},
					Models: []interface{}{User{}},
				}),
				userModule,
			},
		})
		return appModule
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("/app")
	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	req, err := http.NewRequest("POST", testServer.URL+"/app/users", nil)
	require.Nil(t, err)

	req.Header.Set("x-tenant-id", "kaka")

	resp, err := testClient.Do(req)
	require.Nil(t, err)

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	req, err = http.NewRequest("POST", testServer.URL+"/app/users", nil)
	require.Nil(t, err)

	req.Header.Set("x-tenant-id", "roro")

	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func Test_Sync(t *testing.T) {
	type User struct {
		gorm.Model
		Name  string `gorm:"type:varchar(255);not null"`
		Email string `gorm:"type:varchar(255);not null"`
	}

	userController := func(module core.Module) core.Controller {
		ctrl := module.NewController("users")
		ctrl.Post("", func(ctx core.Ctx) error {
			repo := tenancy.InjectRepository[User](module, ctx)
			result, err := repo.Create(&User{Name: "John", Email: "john@gmail.com"})
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
		userMod := module.New(core.NewModuleOptions{
			Imports:     []core.Modules{tenancy.ForFeature(sqlorm.NewRepo(User{}))},
			Controllers: []core.Controllers{userController},
		})
		return userMod
	}

	appModule := func() core.Module {
		appModule := core.NewModule(core.NewModuleOptions{
			Imports: []core.Modules{
				tenancy.ForRoot(tenancy.Options{
					Factory: func(module core.Module) tenancy.ConnectOptions {
						return tenancy.ConnectOptions{
							Host:     "localhost",
							Port:     5432,
							User:     "postgres",
							Password: "postgres",
						}
					},
					GetTenantID: func(r *http.Request) string {
						return r.Header.Get("x-tenant-id")
					},
					Models: []interface{}{User{}},
					Sync:   true,
				}),
				userModule,
			},
		})
		return appModule
	}

	app := core.CreateFactory(appModule)
	app.SetGlobalPrefix("/app")
	testServer := httptest.NewServer(app.PrepareBeforeListen())
	defer testServer.Close()

	testClient := testServer.Client()

	req, err := http.NewRequest("POST", testServer.URL+"/app/users", nil)
	require.Nil(t, err)

	req.Header.Set("x-tenant-id", "abc")

	resp, err := testClient.Do(req)
	require.Nil(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	req, err = http.NewRequest("POST", testServer.URL+"/app/users", nil)
	require.Nil(t, err)

	req.Header.Set("x-tenant-id", "xyz")

	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = testClient.Post(testServer.URL+"/app/users", "application/json", nil)
	require.Nil(t, err)
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func Test_NIl(t *testing.T) {
	type User struct {
		gorm.Model
		Name  string `gorm:"type:varchar(255);not null"`
		Email string `gorm:"type:varchar(255);not null"`
	}

	require.Nil(t, tenancy.InjectRepository[User](core.NewModule(core.NewModuleOptions{}), nil))
}
