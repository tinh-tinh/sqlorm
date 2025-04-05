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
		repo.BatchCreate([]*Documents{
			{Name: "test", Status: "active", Priority: 1},
			{Name: "test2", Status: "active", Priority: 2},
			{Name: "test3", Status: "active", Priority: 3},
		})
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
