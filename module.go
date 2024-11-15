package sqlorm

import (
	"fmt"

	"github.com/tinh-tinh/tinhtinh/common"
	"github.com/tinh-tinh/tinhtinh/core"
	"gorm.io/gorm"
)

type Options struct {
	Dialect gorm.Dialector
	Factory func(module *core.DynamicModule) gorm.Dialector
	Models  []interface{}
}

const ConnectDB core.Provide = "ConnectDB"

func ForRoot(opt Options, configs ...gorm.Option) core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
		var dialector gorm.Dialector
		if opt.Factory != nil {
			dialector = opt.Factory(module)
		} else {
			dialector = opt.Dialect
		}
		conn, err := gorm.Open(dialector, configs...)
		if err != nil {
			panic(err)
		}
		conn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
		fmt.Println("connected to database")

		err = conn.AutoMigrate(opt.Models...)
		if err != nil {
			panic(err)
		}
		fmt.Println("Migrated successful")

		sqlModule := module.New(core.NewModuleOptions{})
		sqlModule.NewProvider(core.ProviderOptions{
			Name:  ConnectDB,
			Value: conn,
		})
		sqlModule.Export(ConnectDB)

		return sqlModule
	}
}

func Inject(module *core.DynamicModule) *gorm.DB {
	db, ok := module.Ref(ConnectDB).(*gorm.DB)
	if !ok {
		return nil
	}
	return db
}

func InjectRepository[M any](module *core.DynamicModule) *Repository[M] {
	var model M
	modelName := core.Provide(fmt.Sprintf("%sRepo", common.GetStructName(model)))
	data, ok := module.Ref(modelName).(*Repository[M])
	fmt.Println(data)
	if !ok {
		return nil
	}

	return data
}

func ForFeature(val ...RepoCommon) core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
		modelModule := module.New(core.NewModuleOptions{})

		for _, v := range val {
			name := GetRepoName(v.GetName())

			modelModule.NewProvider(core.ProviderOptions{
				Name: name,
				Factory: func(param ...interface{}) interface{} {
					connect := param[0].(*gorm.DB)
					if connect != nil {
						v.SetDB(connect)
					}
					return v
				},
				Inject: []core.Provide{ConnectDB},
			})
			modelModule.Export(name)
		}

		return modelModule
	}
}

func GetRepoName(name string) core.Provide {
	return core.Provide(fmt.Sprintf("%sRepo", name))
}
