package sqlorm_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tinh-tinh/sqlorm/v2"
)

func TestToSqlString(t *testing.T) {
	// Test case với giá trị hợp lệ
	value := "test"
	result := sqlorm.ToSqlNullString(&value)
	assert.True(t, result.Valid)
	assert.Equal(t, value, result.String)

	// Test case với nil
	result = sqlorm.ToSqlNullString(nil)
	assert.False(t, result.Valid)
}

func TestToSqlInt16(t *testing.T) {
	value := int16(42)
	result := sqlorm.ToSqlNullInt16(&value)
	assert.True(t, result.Valid)
	assert.Equal(t, value, result.Int16)

	result = sqlorm.ToSqlNullInt16(nil)
	assert.False(t, result.Valid)
}

func TestToSqlInt64(t *testing.T) {
	value := int64(42)
	result := sqlorm.ToSqlNullInt64(&value)
	assert.True(t, result.Valid)
	assert.Equal(t, value, result.Int64)

	result = sqlorm.ToSqlNullInt64(nil)
	assert.False(t, result.Valid)
}

func TestToSqlInt32(t *testing.T) {
	value := int32(42)
	result := sqlorm.ToSqlNullInt32(&value)
	assert.True(t, result.Valid)
	assert.Equal(t, value, result.Int32)

	result = sqlorm.ToSqlNullInt32(nil)
	assert.False(t, result.Valid)
}

func TestToSqlFloat64(t *testing.T) {
	value := float64(3.14)
	result := sqlorm.ToSqlNullFloat64(&value)
	assert.True(t, result.Valid)
	assert.Equal(t, value, result.Float64)

	result = sqlorm.ToSqlNullFloat64(nil)
	assert.False(t, result.Valid)
}

func TestToSqlBool(t *testing.T) {
	value := true
	result := sqlorm.ToSqlNullBool(&value)
	assert.True(t, result.Valid)
	assert.Equal(t, value, result.Bool)

	result = sqlorm.ToSqlNullBool(nil)
	assert.False(t, result.Valid)
}

func TestToSqlTime(t *testing.T) {
	now := time.Now()
	result := sqlorm.ToSqlNullTime(&now)
	assert.True(t, result.Valid)
	assert.Equal(t, now, result.Time)

	result = sqlorm.ToSqlNullTime(nil)
	assert.False(t, result.Valid)
}

func TestToSqlByte(t *testing.T) {
	value := []byte{42}
	result := sqlorm.ToSqlNullByte(&value)
	assert.True(t, result.Valid)
	assert.Equal(t, value[0], result.Byte)

	// Test với slice rỗng
	empty := []byte{}
	result = sqlorm.ToSqlNullByte(&empty)
	assert.False(t, result.Valid)

	// Test với nil
	result = sqlorm.ToSqlNullByte(nil)
	assert.False(t, result.Valid)
}
