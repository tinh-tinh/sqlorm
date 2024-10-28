package sqlorm

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Test_FindMany(t *testing.T) {
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
	result, err := repo.FindAll(map[string]interface{}{"name": "haha"}, FindOptions{
		Order:  []string{"name desc"},
		Select: []string{"id", "name"},
		Limit:  1,
		Offset: 2,
	})
	require.Nil(t, err)
	if len(result) > 0 {
		require.Equal(t, "haha", result[0].Name)
	}

	result1, err := repo.FindOne(map[string]interface{}{"name": "haha"}, FindOneOptions{
		Order: []string{"name desc"},
	})
	require.Nil(t, err)
	if result1 != nil {
		require.Equal(t, "haha", result1.Name)
	}

	if result1 != nil {
		result2, err := repo.FindByID(result1.ID.String(), FindOneOptions{
			Select: []string{"name"},
		})
		require.Nil(t, err)
		if result2 != nil {
			require.Equal(t, "haha", result2.Name)
		}
	}

	result3, err := repo.FindOne(map[string]interface{}{"name": "hjbhjgbhjvghjvgh"}, FindOneOptions{
		Select: []string{"name"},
	})
	require.Nil(t, err)
	require.Nil(t, result3)
}
