package sqlorm_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tinh-tinh/sqlorm/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Test_QueryBuilder(t *testing.T) {
	require.NotPanics(t, func() {
		createDatabaseForTest("test")
	})
	dsn := "host=localhost user=postgres password=postgres dbname=test port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	require.Nil(t, err)
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	type Documents struct {
		sqlorm.Model `gorm:"embedded"`
		Name         string `gorm:"type:varchar(255);not null;default:'Unnamed'"`
		Status       string `gorm:"type:varchar(50);default:'inactive'"`
		Priority     int    `gorm:"type:int;default:0"`
	}
	err = db.AutoMigrate(&Documents{})
	require.Nil(t, err)

	repo := sqlorm.Repository[Documents]{DB: db}

	count, err := repo.Count(nil)
	require.Nil(t, err)

	if count == 0 {
		_, err = repo.BatchCreate([]*Documents{
			{Name: "test", Status: "active", Priority: 1},
			{Name: "test2", Status: "active", Priority: 2},
			{Name: "test3", Status: "active", Priority: 3},
		}, 5)
		require.Nil(t, err)
	}
	// Equal case
	docs, err := repo.FindAll(func(qb *sqlorm.QueryBuilder) {
		qb.Equal("Status", "active")
	})
	require.Nil(t, err)
	require.Equal(t, 3, len(docs))

	// NotEqual
	docs, err = repo.FindAll(func(qb *sqlorm.QueryBuilder) {
		qb.NotEqual("Name", "test")
	})
	require.Nil(t, err)
	require.Equal(t, 2, len(docs))
	require.Equal(t, "test2", docs[0].Name)
	require.Equal(t, "test3", docs[1].Name)

	// Not case
	docs, err = repo.FindAll(func(qb *sqlorm.QueryBuilder) {
		qb.Not("Status", "active")
	})
	require.Nil(t, err)
	require.Equal(t, 0, len(docs))

	// Or case
	docs, err = repo.FindAll(func(qb *sqlorm.QueryBuilder) {
		qb.Or("Name", "test").Or("Priority", 3)
	})
	require.Nil(t, err)
	require.Equal(t, 2, len(docs))

	// In case
	docs, err = repo.FindAll(func(qb *sqlorm.QueryBuilder) {
		qb.In("Priority", 1, 2)
	})
	require.Nil(t, err)
	require.Equal(t, 2, len(docs))
	require.Equal(t, "test", docs[0].Name)
	require.Equal(t, "test2", docs[1].Name)

	// Not in case
	doc, err := repo.FindOne(func(qb *sqlorm.QueryBuilder) {
		qb.NotIn("Priority", 1, 2)
	})
	require.Nil(t, err)
	require.Equal(t, "test3", doc.Name)

	// More than case
	// More than case
	docs, err = repo.FindAll(func(qb *sqlorm.QueryBuilder) {
		qb.MoreThan("Priority", 1)
	})
	require.Nil(t, err)
	require.Equal(t, 2, len(docs))
	require.Equal(t, "test2", docs[0].Name)
	require.Equal(t, "test3", docs[1].Name)

	// More than or equal case
	docs, err = repo.FindAll(func(qb *sqlorm.QueryBuilder) {
		qb.MoreThanOrEqual("Priority", 2)
	})
	require.Nil(t, err)
	require.Equal(t, 2, len(docs))
	require.Equal(t, "test2", docs[0].Name)
	require.Equal(t, "test3", docs[1].Name)

	// Less than case
	docs, err = repo.FindAll(func(qb *sqlorm.QueryBuilder) {
		qb.LessThan("Priority", 3)
	})
	require.Nil(t, err)
	require.Equal(t, 2, len(docs))
	require.Equal(t, "test", docs[0].Name)
	require.Equal(t, "test2", docs[1].Name)

	// Less than or equal case
	docs, err = repo.FindAll(func(qb *sqlorm.QueryBuilder) {
		qb.LessThanOrEqual("Priority", 2)
	})
	require.Nil(t, err)
	require.Equal(t, 2, len(docs))
	require.Equal(t, "test", docs[0].Name)
	require.Equal(t, "test2", docs[1].Name)

	// Like
	docs, err = repo.FindAll(func(qb *sqlorm.QueryBuilder) {
		qb.Like("Name", "%test%")
	})
	require.Nil(t, err)
	require.Equal(t, 3, len(docs))
	require.Equal(t, "test", docs[0].Name)
	require.Equal(t, "test2", docs[1].Name)
	require.Equal(t, "test3", docs[2].Name)

	// Ilike
	docs, err = repo.FindAll(func(qb *sqlorm.QueryBuilder) {
		qb.ILike("Name", "%TEST%")
	})
	require.Nil(t, err)
	require.Equal(t, 3, len(docs))
	require.Equal(t, "test", docs[0].Name)
	require.Equal(t, "test2", docs[1].Name)
	require.Equal(t, "test3", docs[2].Name)

	// Between case
	docs, err = repo.FindAll(func(qb *sqlorm.QueryBuilder) {
		qb.Between("Priority", 1, 2)
	})
	require.Nil(t, err)
	require.Equal(t, 2, len(docs))
	require.Equal(t, "test", docs[0].Name)
	require.Equal(t, "test2", docs[1].Name)

	// Is null case
	exist, err := repo.Exist(func(qb *sqlorm.QueryBuilder) {
		qb.IsNull("Name")
	})
	require.Nil(t, err)
	require.Equal(t, false, exist)

	// Raw
	docs, err = repo.FindAll(func(qb *sqlorm.QueryBuilder) {
		qb.Raw("SELECT * FROM documents WHERE Name = ?", "test")
	})
	require.Nil(t, err)
	require.Equal(t, 1, len(docs))
	require.Equal(t, "test", docs[0].Name)
}

func Test_IsValidColumn(t *testing.T) {
	require.NotPanics(t, func() {
		createDatabaseForTest("test_valid_column")
	})
	dsn := "host=localhost user=postgres password=postgres dbname=test_valid_column port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	require.Nil(t, err)
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	type TestEntity struct {
		sqlorm.Model `gorm:"embedded"`
		Name         string `gorm:"type:varchar(255)"`
		Value        int    `gorm:"type:int"`
	}
	err = db.AutoMigrate(&TestEntity{})
	require.Nil(t, err)

	repo := sqlorm.Repository[TestEntity]{DB: db}

	// Create test data
	count, err := repo.Count(nil)
	require.Nil(t, err)

	if count == 0 {
		_, err = repo.Create(&TestEntity{Name: "valid", Value: 1})
		require.Nil(t, err)
	}

	// Define invalid column names to test
	invalidColumns := []string{
		"name; DROP TABLE users;--",
		"column' OR '1'='1",
		"col=1",
		"column!",
		"col@name",
		"col#name",
		"col$name",
		"col%name",
		"col^name",
		"col&name",
		"col*name",
		"column()",
		"[column]",
		"<column>",
		"col/name",
		"col\\name",
		"col|name",
		"col`name",
		"col~name",
		"col+name",
		"col name",
		"col-name",
		"col.name",
		"col:name",
		"col;name",
		"col'name",
		"col\"name",
		"col{name}",
	}

	// Test Equal with invalid columns
	t.Run("Equal_InvalidColumn", func(t *testing.T) {
		for _, invalidCol := range invalidColumns {
			docs, err := repo.FindAll(func(qb *sqlorm.QueryBuilder) {
				qb.Equal(invalidCol, "test")
			})
			require.Nil(t, err)
			require.GreaterOrEqual(t, len(docs), 1, "Invalid column %q should not filter results", invalidCol)
		}
	})

	// Test Not with invalid columns
	t.Run("Not_InvalidColumn", func(t *testing.T) {
		for _, invalidCol := range invalidColumns {
			docs, err := repo.FindAll(func(qb *sqlorm.QueryBuilder) {
				qb.Not(invalidCol, "test")
			})
			require.Nil(t, err)
			require.GreaterOrEqual(t, len(docs), 1, "Invalid column %q should not filter results", invalidCol)
		}
	})

	// Test Or with invalid columns
	t.Run("Or_InvalidColumn", func(t *testing.T) {
		for _, invalidCol := range invalidColumns {
			docs, err := repo.FindAll(func(qb *sqlorm.QueryBuilder) {
				qb.Or(invalidCol, "test")
			})
			require.Nil(t, err)
			require.GreaterOrEqual(t, len(docs), 1, "Invalid column %q should not filter results", invalidCol)
		}
	})

	// Test In with invalid columns
	t.Run("In_InvalidColumn", func(t *testing.T) {
		for _, invalidCol := range invalidColumns {
			docs, err := repo.FindAll(func(qb *sqlorm.QueryBuilder) {
				qb.In(invalidCol, "test", "test2")
			})
			require.Nil(t, err)
			require.GreaterOrEqual(t, len(docs), 1, "Invalid column %q should not filter results", invalidCol)
		}
	})

	// Test MoreThan with invalid columns
	t.Run("MoreThan_InvalidColumn", func(t *testing.T) {
		for _, invalidCol := range invalidColumns {
			docs, err := repo.FindAll(func(qb *sqlorm.QueryBuilder) {
				qb.MoreThan(invalidCol, 0)
			})
			require.Nil(t, err)
			require.GreaterOrEqual(t, len(docs), 1, "Invalid column %q should not filter results", invalidCol)
		}
	})

	// Test MoreThanOrEqual with invalid columns
	t.Run("MoreThanOrEqual_InvalidColumn", func(t *testing.T) {
		for _, invalidCol := range invalidColumns {
			docs, err := repo.FindAll(func(qb *sqlorm.QueryBuilder) {
				qb.MoreThanOrEqual(invalidCol, 0)
			})
			require.Nil(t, err)
			require.GreaterOrEqual(t, len(docs), 1, "Invalid column %q should not filter results", invalidCol)
		}
	})

	// Test LessThan with invalid columns
	t.Run("LessThan_InvalidColumn", func(t *testing.T) {
		for _, invalidCol := range invalidColumns {
			docs, err := repo.FindAll(func(qb *sqlorm.QueryBuilder) {
				qb.LessThan(invalidCol, 100)
			})
			require.Nil(t, err)
			require.GreaterOrEqual(t, len(docs), 1, "Invalid column %q should not filter results", invalidCol)
		}
	})

	// Test LessThanOrEqual with invalid columns
	t.Run("LessThanOrEqual_InvalidColumn", func(t *testing.T) {
		for _, invalidCol := range invalidColumns {
			docs, err := repo.FindAll(func(qb *sqlorm.QueryBuilder) {
				qb.LessThanOrEqual(invalidCol, 100)
			})
			require.Nil(t, err)
			require.GreaterOrEqual(t, len(docs), 1, "Invalid column %q should not filter results", invalidCol)
		}
	})

	// Test Like with invalid columns
	t.Run("Like_InvalidColumn", func(t *testing.T) {
		for _, invalidCol := range invalidColumns {
			docs, err := repo.FindAll(func(qb *sqlorm.QueryBuilder) {
				qb.Like(invalidCol, "%test%")
			})
			require.Nil(t, err)
			require.GreaterOrEqual(t, len(docs), 1, "Invalid column %q should not filter results", invalidCol)
		}
	})

	// Test ILike with invalid columns
	t.Run("ILike_InvalidColumn", func(t *testing.T) {
		for _, invalidCol := range invalidColumns {
			docs, err := repo.FindAll(func(qb *sqlorm.QueryBuilder) {
				qb.ILike(invalidCol, "%TEST%")
			})
			require.Nil(t, err)
			require.GreaterOrEqual(t, len(docs), 1, "Invalid column %q should not filter results", invalidCol)
		}
	})

	// Test Between with invalid columns
	t.Run("Between_InvalidColumn", func(t *testing.T) {
		for _, invalidCol := range invalidColumns {
			docs, err := repo.FindAll(func(qb *sqlorm.QueryBuilder) {
				qb.Between(invalidCol, 0, 100)
			})
			require.Nil(t, err)
			require.GreaterOrEqual(t, len(docs), 1, "Invalid column %q should not filter results", invalidCol)
		}
	})

	// Test NotEqual with invalid columns
	t.Run("NotEqual_InvalidColumn", func(t *testing.T) {
		for _, invalidCol := range invalidColumns {
			docs, err := repo.FindAll(func(qb *sqlorm.QueryBuilder) {
				qb.NotEqual(invalidCol, "nonexistent")
			})
			require.Nil(t, err)
			require.GreaterOrEqual(t, len(docs), 1, "Invalid column %q should not filter results", invalidCol)
		}
	})

	// Test NotIn with invalid columns
	t.Run("NotIn_InvalidColumn", func(t *testing.T) {
		for _, invalidCol := range invalidColumns {
			docs, err := repo.FindAll(func(qb *sqlorm.QueryBuilder) {
				qb.NotIn(invalidCol, "nonexistent1", "nonexistent2")
			})
			require.Nil(t, err)
			require.GreaterOrEqual(t, len(docs), 1, "Invalid column %q should not filter results", invalidCol)
		}
	})

	// Test IsNull with invalid columns
	t.Run("IsNull_InvalidColumn", func(t *testing.T) {
		for _, invalidCol := range invalidColumns {
			docs, err := repo.FindAll(func(qb *sqlorm.QueryBuilder) {
				qb.IsNull(invalidCol)
			})
			require.Nil(t, err)
			require.GreaterOrEqual(t, len(docs), 1, "Invalid column %q should not filter results", invalidCol)
		}
	})

	// Test valid column names still work
	t.Run("ValidColumn_Equal", func(t *testing.T) {
		docs, err := repo.FindAll(func(qb *sqlorm.QueryBuilder) {
			qb.Equal("Name", "valid")
		})
		require.Nil(t, err)
		require.Equal(t, 1, len(docs))
		require.Equal(t, "valid", docs[0].Name)
	})

	t.Run("ValidColumn_MoreThan", func(t *testing.T) {
		docs, err := repo.FindAll(func(qb *sqlorm.QueryBuilder) {
			qb.MoreThan("Value", 0)
		})
		require.Nil(t, err)
		require.Equal(t, 1, len(docs))
	})
}
