package sqlorm_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/sqlorm/v2"
)

func Test_Map(t *testing.T) {
	type Todo struct {
		sqlorm.Model `gorm:"embedded"`
		Name         string `gorm:"type:varchar(255);not null"`
	}

	type CreateTodo struct {
		Name string
		Haha string
		Hihi string
	}
	data := sqlorm.MapOne[Todo](&CreateTodo{Name: "haha", Haha: "haha", Hihi: "hihi"})
	require.Equal(t, "haha", data.Name)

	data2 := sqlorm.MapOne[Todo](map[string]interface{}{
		"Name": "hihi",
		"Age":  6,
	})
	require.Equal(t, "hihi", data2.Name)

	data3 := sqlorm.MapOne[Todo](nil)
	require.Empty(t, data3.Name)

	data4 := sqlorm.MapOne[Todo]("abc")
	require.Empty(t, data4.Name)

	input := []*CreateTodo{
		{Name: "abc", Haha: "haha", Hihi: "hihi"},
		{Name: "def", Haha: "haha", Hihi: "hihi"},
		{Name: "ghi", Haha: "haha", Hihi: "hihi"},
		{Name: "jkl", Haha: "haha", Hihi: "hihi"},
	}
	lists := sqlorm.MapMany[Todo](input)
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

	lists2 := sqlorm.MapMany[Todo](input2)
	require.Equal(t, 4, len(lists2))
	require.Equal(t, "abc", lists2[0].Name)
	require.Equal(t, "def", lists2[1].Name)
	require.Equal(t, "ghi", lists2[2].Name)
	require.Equal(t, "jkl", lists2[3].Name)

	list3 := sqlorm.MapMany[Todo](map[string]interface{}{
		"Name": "abc",
		"Age":  6,
	})
	require.Equal(t, 0, len(list3))

	list4 := sqlorm.MapMany[Todo](nil)
	require.Equal(t, 0, len(list4))
}

type Abc struct {
	sqlorm.Model `gorm:"embedded"`
	Name         string `gorm:"type:varchar(255);not null"`
}

func (Abc) RepositoryName() string {
	return "service"
}

func TestNameRepo(t *testing.T) {
	repo := sqlorm.NewRepo(Abc{})
	require.Equal(t, "service", repo.GetName())
}
