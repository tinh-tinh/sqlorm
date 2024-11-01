package sqlorm

import (
	"fmt"

	"github.com/tinh-tinh/tinhtinh/core"
	"github.com/tinh-tinh/tinhtinh/utils"
	"gorm.io/gorm"
)

type Options struct {
	Dialect gorm.Dialector
	Factory func(module *core.DynamicModule) gorm.Dialector
	Models  []interface{}
}

const ConnectDB core.Provide = "ConnectDB"

func ForRoot(opt Options) core.Module {
	return func(module *core.DynamicModule) *core.DynamicModule {
		var dialector gorm.Dialector
		if opt.Factory != nil {
			dialector = opt.Factory(module)
		} else {
			dialector = opt.Dialect
		}
		conn, err := gorm.Open(dialector, &gorm.Config{})
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
	modelName := core.Provide(utils.GetNameStruct(model))
	data, ok := module.Ref(modelName).(*Repository[M])
	if data == nil || !ok {
		repo := Repository[M]{DB: Inject(module)}
		module.NewProvider(core.ProviderOptions{
			Name:  modelName,
			Value: &repo,
		})

		return &repo
	}

	return data
}
