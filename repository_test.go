package sqlorm

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Test_Create(t *testing.T) {
	require.NotPanics(t, func() {
		createDatabaseForTest("test")
	})
	dsn := "host=localhost user=postgres password=postgres dbname=test port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.Nil(t, err)

	type Todo struct {
		Model `gorm:"embedded"`
		Name  string `gorm:"type:varchar(255);not null"`
	}
	err = db.AutoMigrate(&Todo{})
	require.Nil(t, err)

	repo := Repository[Todo]{DB: db}

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
	require.NotPanics(t, func() {
		createDatabaseForTest("test")
	})
	dsn := "host=localhost user=postgres password=postgres dbname=test port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.Nil(t, err)

	type Todo struct {
		Model `gorm:"embedded"`
		Name  string `gorm:"type:varchar(255);not null"`
	}
	err = db.AutoMigrate(&Todo{})
	require.Nil(t, err)

	repo := Repository[Todo]{DB: db}

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

func Test_Map(t *testing.T) {
	type Todo struct {
		Model `gorm:"embedded"`
		Name  string `gorm:"type:varchar(255);not null"`
	}

	type CreateTodo struct {
		Name string
		Haha string
		Hihi string
	}
	data := MapOne[Todo](&CreateTodo{Name: "haha", Haha: "haha", Hihi: "hihi"})
	require.Equal(t, "haha", data.Name)

	data2 := MapOne[Todo](map[string]interface{}{
		"Name": "hihi",
		"Age":  6,
	})
	require.Equal(t, "hihi", data2.Name)

	data3 := MapOne[Todo](nil)
	require.Empty(t, data3.Name)

	data4 := MapOne[Todo]("abc")
	require.Empty(t, data4.Name)

	input := []*CreateTodo{
		{Name: "abc", Haha: "haha", Hihi: "hihi"},
		{Name: "def", Haha: "haha", Hihi: "hihi"},
		{Name: "ghi", Haha: "haha", Hihi: "hihi"},
		{Name: "jkl", Haha: "haha", Hihi: "hihi"},
	}
	lists := MapMany[Todo](input)
	require.Equal(t, 4, len(lists))
	require.Equal(t, "abc", lists[0].Name)
	require.Equal(t, "def", lists[1].Name)
	require.Equal(t, "ghi", lists[2].Name)
	require.Equal(t, "jkl", lists[3].Name)

	input2 := []map[string]interface{}{
		{"Name": "abc", "Haha": "haha", "Hihi": "hihi"},
		{"Name": "def", "Haha": "haha", "Hihi": "hihi"},
		{"Name": "ghi", "Haha": "haha", "Hihi": "hihi"},
		{"name": "jkl", "Haha": "haha", "Hihi": "hihi"},
	}

	lists2 := MapMany[Todo](input2)
	require.Equal(t, 4, len(lists2))
	require.Equal(t, "abc", lists2[0].Name)
	require.Equal(t, "def", lists2[1].Name)
	require.Equal(t, "ghi", lists2[2].Name)
	require.Equal(t, "jkl", lists2[3].Name)

	list3 := MapMany[Todo](map[string]interface{}{
		"Name": "abc",
		"Age":  6,
	})
	require.Equal(t, 0, len(list3))

	list4 := MapMany[Todo](nil)
	require.Equal(t, 0, len(list4))
}
