package tenancy_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/sqlorm"
	"github.com/tinh-tinh/sqlorm/tenancy"
	"github.com/tinh-tinh/tinhtinh/common"
	"github.com/tinh-tinh/tinhtinh/core"
)

func Test_Tenant(t *testing.T) {
	type User struct {
		sqlorm.Model `gorm:"embedded"`
		Name         string `gorm:"type:varchar(255);not null"`
		Email        string `gorm:"type:varchar(255);not null"`
	}

	userController := func(module *core.DynamicModule) *core.DynamicController {
		ctrl := module.NewController("users")
		ctrl.Post("", func(ctx core.Ctx) error {
			repo := tenancy.InjectRepository[User](module)
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

	userModule := func(module *core.DynamicModule) *core.DynamicModule {
		userMod := module.New(core.NewModuleOptions{
			Imports:     []core.Module{tenancy.ForFeature[User]()},
			Controllers: []core.Controller{userController},
		})
		return userMod
	}

	appModule := func() *core.DynamicModule {
		appModule := core.NewModule(core.NewModuleOptions{
			Imports: []core.Module{
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
					Models: []interface{}{&User{}},
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

	req.Header.Set("x-tenant-id", "anc")

	resp, err := testClient.Do(req)
	require.Nil(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	req, err = http.NewRequest("POST", testServer.URL+"/app/users", nil)
	require.Nil(t, err)

	req.Header.Set("x-tenant-id", "xyz")

	resp, err = testClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
