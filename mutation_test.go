package sqlorm_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/sqlorm/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Test_Create(t *testing.T) {
	db := prepareBeforeTest(t)

	type Todo struct {
		sqlorm.Model `gorm:"embedded"`
		Name         string `gorm:"type:varchar(255);not null"`
	}
	err := db.AutoMigrate(&Todo{})
	require.Nil(t, err)

	repo := sqlorm.Repository[Todo]{DB: db}

	require.NotPanics(t, func() {
		type CreateTodo struct {
			Name string
			Haha string
			Hihi string
		}
		result, err := repo.Create(&CreateTodo{Name: "haha", Haha: "haha", Hihi: "hihi"})
		require.Nil(t, err)
		require.Equal(t, "haha", result.Name)
	})
}

func Test_BatchCreate(t *testing.T) {
	db := prepareBeforeTest(t)

	type Todo struct {
		sqlorm.Model `gorm:"embedded"`
		Name         string `gorm:"type:varchar(255);not null"`
	}
	err := db.AutoMigrate(&Todo{})
	require.Nil(t, err)

	repo := sqlorm.Repository[Todo]{DB: db}

	require.NotPanics(t, func() {
		type CreateTodo struct {
			Name string
			Haha string
			Hihi string
		}
		result, err := repo.BatchCreate([]*CreateTodo{
			{Name: "abc", Haha: "haha", Hihi: "hihi"},
			{Name: "def", Haha: "haha", Hihi: "hihi"},
			{Name: "ghi", Haha: "haha", Hihi: "hihi"},
			{Name: "jkl", Haha: "haha", Hihi: "hihi"},
		})
		require.Nil(t, err)
		require.Len(t, result, 4)
		require.Equal(t, "abc", result[0].Name)
		require.Equal(t, "def", result[1].Name)
		require.Equal(t, "ghi", result[2].Name)
		require.Equal(t, "jkl", result[3].Name)
	})
}

func Test_Update(t *testing.T) {
	db := prepareBeforeTest(t)

	type Todo struct {
		sqlorm.Model `gorm:"embedded"`
		Name         string `gorm:"type:varchar(255);not null"`
	}
	err := db.AutoMigrate(&Todo{})
	require.Nil(t, err)

	repo := sqlorm.Repository[Todo]{DB: db}
	require.NotPanics(t, func() {
		type CreateTodo struct {
			Name string
			Haha string
			Hihi string
		}
		result, err := repo.Create(&CreateTodo{Name: "Babadook", Haha: "haha", Hihi: "hihi"})
		require.Nil(t, err)
		require.Equal(t, "Babadook", result.Name)

		type UpdateTodo struct {
			Name string
			Haha string
			Hihi string
		}
		result, err = repo.UpdateOne(map[string]interface{}{"id": result.ID.String()}, &UpdateTodo{Name: "haha", Haha: "haha", Hihi: "hihi"})
		require.Nil(t, err)
		require.Equal(t, "haha", result.Name)
	})
}

func Test_UpdateMany(t *testing.T) {
	db := prepareBeforeTest(t)

	type Todo struct {
		sqlorm.Model `gorm:"embedded"`
		Name         string `gorm:"type:varchar(255);not null"`
	}
	err := db.AutoMigrate(&Todo{})
	require.Nil(t, err)

	repo := sqlorm.Repository[Todo]{DB: db}
	require.NotPanics(t, func() {
		err := repo.UpdateMany(map[string]interface{}{"name": "haha"}, map[string]interface{}{"name": "lulu"})
		require.Nil(t, err)

		err = repo.UpdateMany(nil, map[string]interface{}{"name": "mahula"})
		require.Nil(t, err)
	})
}

func Test_Delete(t *testing.T) {
	db := prepareBeforeTest(t)

	type Todo struct {
		sqlorm.Model `gorm:"embedded"`
		Name         string `gorm:"type:varchar(255);not null"`
	}
	err := db.AutoMigrate(&Todo{})
	require.Nil(t, err)

	repo := sqlorm.Repository[Todo]{DB: db}
	require.NotPanics(t, func() {
		type CreateTodo struct {
			Name string
			Haha string
			Hihi string
		}
		result, err := repo.Create(&CreateTodo{Name: "Babadook", Haha: "haha", Hihi: "hihi"})
		require.Nil(t, err)
		require.Equal(t, "Babadook", result.Name)

		err = repo.DeleteOne(map[string]interface{}{"name": "Babadook"})
		require.Nil(t, err)

		err = repo.DeleteOne(map[string]interface{}{"name": "Luxembuar"})
		require.NotNil(t, err)
	})
}

func Test_DeleteMany(t *testing.T) {
	db := prepareBeforeTest(t)

	type Todo struct {
		sqlorm.Model `gorm:"embedded"`
		Name         string `gorm:"type:varchar(255);not null"`
	}
	err := db.AutoMigrate(&Todo{})
	require.Nil(t, err)

	repo := sqlorm.Repository[Todo]{DB: db}
	require.NotPanics(t, func() {
		err := repo.DeleteMany(map[string]interface{}{"name": "lulu"})
		require.Nil(t, err)

		err = repo.DeleteMany(nil)
		require.Nil(t, err)
	})
}

func prepareBeforeTest(t *testing.T) *gorm.DB {
	require.NotPanics(t, func() {
		createDatabaseForTest("test")
	})
	dsn := "host=localhost user=postgres password=postgres dbname=test port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.Nil(t, err)

	return db
}
