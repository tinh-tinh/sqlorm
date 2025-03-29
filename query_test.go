package sqlorm_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/sqlorm/v2"
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
		sqlorm.Model `gorm:"embedded"`
		Name         string `gorm:"type:varchar(255);not null"`
	}
	err = db.AutoMigrate(&Todo{})
	require.Nil(t, err)

	repo := sqlorm.Repository[Todo]{DB: db}
	result, err := repo.FindAll(map[string]interface{}{"name": "haha"}, sqlorm.FindOptions{
		Order:  []string{"name desc"},
		Select: []string{"id", "name"},
		Limit:  1,
		Offset: 2,
	})
	require.Nil(t, err)
	if len(result) > 0 {
		require.Equal(t, "haha", result[0].Name)
	}

	result1, err := repo.FindOne(map[string]interface{}{"name": "haha"}, sqlorm.FindOneOptions{
		Order: []string{"name desc"},
	})
	require.Nil(t, err)
	if result1 != nil {
		require.Equal(t, "haha", result1.Name)
	}

	if result1 != nil {
		result2, err := repo.FindByID(result1.ID.String(), sqlorm.FindOneOptions{
			Select: []string{"name"},
		})
		require.Nil(t, err)
		if result2 != nil {
			require.Equal(t, "haha", result2.Name)
		}
	}

	result3, err := repo.FindOne(map[string]interface{}{"name": "hjbhjgbhjvghjvgh"}, sqlorm.FindOneOptions{
		Select: []string{"name"},
	})
	require.Nil(t, err)
	require.Nil(t, result3)
}

func Test_Count(t *testing.T) {
	require.NotPanics(t, func() {
		createDatabaseForTest("test")
	})
	dsn := "host=localhost user=postgres password=postgres dbname=test port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.Nil(t, err)

	type Count struct {
		sqlorm.Model `gorm:"embedded"`
		Name         string `gorm:"type:varchar(255);not null"`
		Status       string `gorm:"type:varchar(50)"`
		Priority     int    `gorm:"type:int"`
	}
	err = db.AutoMigrate(&Count{})
	require.Nil(t, err)

	repo := sqlorm.Repository[Count]{DB: db}

	// Create test data
	todos := []Count{
		{Name: "Test Todo 1", Status: "active", Priority: 1},
		{Name: "Test Todo 2", Status: "active", Priority: 2},
		{Name: "Test Todo 3", Status: "completed", Priority: 1},
		{Name: "Test Todo 4", Status: "completed", Priority: 3},
	}
	// Check if database is empty before creating test data
	existingCount, err := repo.Count(nil)
	require.Nil(t, err)
	if existingCount == 0 {
		// Only create test data if database is empty
		for _, todo := range todos {
			err := db.Create(&todo).Error
			require.Nil(t, err)
		}
	}

	// Test counting with multiple conditions
	count, err := repo.Count(map[string]interface{}{
		"status":   "active",
		"priority": 1,
	})
	require.Nil(t, err)
	require.Equal(t, int64(1), count)

	// Test counting with empty where condition
	count, err = repo.Count(nil)
	require.Nil(t, err)
	require.Equal(t, int64(4), count)

	// Test counting with non-existent conditions
	count, err = repo.Count(map[string]interface{}{
		"status":   "non-existent",
		"priority": 999,
	})
	require.Nil(t, err)
	require.Equal(t, int64(0), count)

	// Test counting after soft delete
	todoToDelete := todos[0]
	err = db.Where("priority = 1").Delete(&todoToDelete).Error
	require.Nil(t, err)

	// Should not count soft deleted records by default
	count, err = repo.Count(map[string]interface{}{
		"status": "active",
	})
	require.Nil(t, err)
	require.Equal(t, int64(1), count)

	// Reset test data by deleting all records
	err = db.Unscoped().Where("1 = 1").Delete(&Count{}).Error
	require.Nil(t, err)

	// Verify database is empty
	count, err = repo.Count(nil)
	require.Nil(t, err)
	require.Equal(t, int64(0), count)
}

func Test_Exist(t *testing.T) {
	require.NotPanics(t, func() {
		createDatabaseForTest("test")
	})
	dsn := "host=localhost user=postgres password=postgres dbname=test port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.Nil(t, err)

	type Exists struct {
		sqlorm.Model `gorm:"embedded"`
		Name         string `gorm:"type:varchar(255);not null"`
		Status       string `gorm:"type:varchar(50)"`
		Priority     int    `gorm:"type:int"`
	}
	err = db.AutoMigrate(&Exists{})
	require.Nil(t, err)

	repo := sqlorm.Repository[Exists]{DB: db}

	// Create test data
	todos := []Exists{
		{Name: "Test Todo 1", Status: "active", Priority: 1},
		{Name: "Test Todo 2", Status: "active", Priority: 2},
		{Name: "Test Todo 3", Status: "completed", Priority: 1},
	}
	for _, todo := range todos {
		err := db.Create(&todo).Error
		require.Nil(t, err)
	}

	// Test existing record with single condition
	exists, err := repo.Exist(map[string]interface{}{"name": "Test Todo 1"})
	require.Nil(t, err)
	require.True(t, exists)

	// Test existing record with multiple conditions
	exists, err = repo.Exist(map[string]interface{}{
		"name":     "Test Todo 1",
		"status":   "active",
		"priority": 1,
	})
	require.Nil(t, err)
	require.True(t, exists)

	// Test existing record with select option
	exists, err = repo.Exist(map[string]interface{}{"name": "Test Todo 1"}, sqlorm.FindOneOptions{
		Select: []string{"name", "status"},
	})
	require.Nil(t, err)
	require.True(t, exists)

	// Test non-existing record with multiple conditions
	exists, err = repo.Exist(map[string]interface{}{
		"name":   "Test Todo 1",
		"status": "completed",
	})
	require.Nil(t, err)
	require.False(t, exists)

	// Test with order option
	exists, err = repo.Exist(map[string]interface{}{"name": "Test Todo 1"}, sqlorm.FindOneOptions{
		Order: []string{"name desc"},
	})
	require.Nil(t, err)
	require.True(t, exists)

	// Reset test data by deleting all records
	err = db.Unscoped().Where("1 = 1").Delete(&Exists{}).Error
	require.Nil(t, err)

	// Verify database is empty
	count, err := repo.Count(nil)
	require.Nil(t, err)
	require.Equal(t, int64(0), count)
}

func Test_SoftDelete(t *testing.T) {
	require.NotPanics(t, func() {
		createDatabaseForTest("test")
	})
	dsn := "host=localhost user=postgres password=postgres dbname=test port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.Nil(t, err)

	type SoftDelete struct {
		sqlorm.Model `gorm:"embedded"`
		Name         string `gorm:"type:varchar(255);not null"`
	}

	repo := sqlorm.Repository[SoftDelete]{DB: db}
	// Create test records
	records := []SoftDelete{
		{Name: "Test Record 1"},
		{Name: "Test Record 2"},
		{Name: "Test Record 3"},
	}
	err = db.AutoMigrate(&SoftDelete{})
	require.Nil(t, err)

	// Check if database is empty before creating test data
	existingCount, err := repo.Count(nil)
	if err != nil || existingCount == 0 {
		// Only create test data if database is empty
		for _, record := range records {
			err := db.Create(&record).Error
			require.Nil(t, err)
		}
	}

	// Test soft delete
	err = repo.DeleteOne(map[string]interface{}{"name": "Test Record 1"})
	require.Nil(t, err)

	// Verify record is soft deleted (not found in normal query)
	exists, err := repo.Exist(map[string]interface{}{"name": "Test Record 1"})
	require.Nil(t, err)
	require.False(t, exists)

	// Verify record exists when unscoped
	var record SoftDelete
	err = db.Unscoped().Where("name = ?", "Test Record 1").First(&record).Error
	require.Nil(t, err)
	require.NotNil(t, record.DeletedAt)

	// Test find with unscoped option
	result, err := repo.FindOne(map[string]interface{}{"name": "Test Record 1"}, sqlorm.FindOneOptions{
		WithDeleted: true,
	})
	require.Nil(t, err)
	require.NotNil(t, result)

	// Reset test data by deleting all records
	err = db.Unscoped().Where("1 = 1").Delete(&SoftDelete{}).Error
	require.Nil(t, err)

	// Verify database is empty
	count, err := repo.Count(nil)
	require.Nil(t, err)
	require.Equal(t, int64(0), count)

}
