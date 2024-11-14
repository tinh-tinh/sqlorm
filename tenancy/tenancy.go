package tenancy

import (
	"fmt"
	"net/http"

	"github.com/tinh-tinh/sqlorm"
	"github.com/tinh-tinh/tinhtinh/common"
	"github.com/tinh-tinh/tinhtinh/core"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ConnectOptions struct {
	Host     string
	Port     int
	User     string
	Password string
}

const (
	CONNECT_MAPPER  core.Provide = "CONNECT_MAPPER"
	CONNECT_TENANCY core.Provide = "CONNECT_TENANCY"
)

type Options struct {
	Connect     ConnectOptions
	Factory     func(module *core.DynamicModule) ConnectOptions
	GetTenantID func(r *http.Request) string
	Models      []interface{}
}

type ConnectMapper map[string]*gorm.DB

func ForRoot(opt Options) core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
		var connectOpt ConnectOptions
		if opt.Factory != nil {
			connectOpt = opt.Factory(module)
		} else {
			connectOpt = opt.Connect
		}

		tenantModule := module.New(core.NewModuleOptions{})

		tenantModule.NewProvider(core.ProviderOptions{
			Name:  CONNECT_MAPPER,
			Value: make(ConnectMapper),
		})
		tenantModule.Export(CONNECT_MAPPER)

		tenantModule.NewProvider(core.ProviderOptions{
			Scope: core.Request,
			Name:  CONNECT_TENANCY,
			Factory: func(param ...interface{}) interface{} {
				req := param[0].(*http.Request)
				tenantID := opt.GetTenantID(req)
				if tenantID == "" {
					return nil
				}
				mapper, ok := param[1].(ConnectMapper)
				if !ok {
					return nil
				}
				if mapper[tenantID] == nil {
					CreateDabaseIfNotExist(tenantID, connectOpt)
					dns := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai", connectOpt.Host, connectOpt.Port, connectOpt.User, connectOpt.Password, tenantID)
					conn, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
					if err != nil {
						return nil
					}
					conn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
					fmt.Println("connected to database")

					err = conn.AutoMigrate(opt.Models...)
					if err != nil {
						panic(err)
					}
					fmt.Println("Migrated successful")

					mapper[tenantID] = conn
				}

				return mapper[tenantID]
			},
			Inject: []core.Provide{core.REQUEST, CONNECT_MAPPER},
		})
		tenantModule.Export(CONNECT_TENANCY)

		return tenantModule
	}
}

func ForFeature[M any]() core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
		modelModule := module.New(core.NewModuleOptions{})

		var model M
		modelName := core.Provide(common.GetStructName(model))
		modelModule.NewProvider(core.ProviderOptions{
			Scope: core.Request,
			Name:  modelName,
			Factory: func(param ...interface{}) interface{} {
				connect := param[0].(*gorm.DB)
				repo := sqlorm.Repository[M]{DB: connect}
				return &repo
			},
			Inject: []core.Provide{CONNECT_TENANCY},
		})
		modelModule.Export(modelName)

		return modelModule
	}
}

func InjectRepository[M any](module *core.DynamicModule) *sqlorm.Repository[M] {
	var model M
	modelName := core.Provide(common.GetStructName(model))
	data, ok := module.Ref(modelName).(*sqlorm.Repository[M])
	if data == nil || !ok {
		return nil
	}
	return data
}

func CreateDabaseIfNotExist(dbName string, opt ConnectOptions) error {
	conStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable TimeZone=Asia/Shanghai", opt.Host, opt.Port, opt.User, opt.Password)
	db, err := gorm.Open(postgres.Open(conStr), &gorm.Config{})
	if err != nil {
		return err
	}

	// check if db exists
	stmt := fmt.Sprintf("SELECT * FROM pg_database WHERE datname = '%s';", dbName)
	rs := db.Raw(stmt)
	if rs.Error != nil {
		return rs.Error
	}

	// if not create it
	var rec = make(map[string]interface{})
	if rs.Find(rec); len(rec) == 0 {
		stmt := fmt.Sprintf("CREATE DATABASE %s;", dbName)
		if rs := db.Exec(stmt); rs.Error != nil {
			return rs.Error
		}

		// close db connection
		sql, err := db.DB()
		defer func() {
			_ = sql.Close()
		}()
		if err != nil {
			return err
		}
	}
	return nil
}
