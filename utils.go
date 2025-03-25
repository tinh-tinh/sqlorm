package sqlorm

import (
	"database/sql"
	"time"
)

func ToSqlNullString(input *string) sql.NullString {
	if input != nil {
		return sql.NullString{String: *input, Valid: true}
	}
	return sql.NullString{Valid: false}
}

func ToSqlNullInt64(input *int64) sql.NullInt64 {
	if input != nil {
		return sql.NullInt64{Int64: *input, Valid: true}
	}
	return sql.NullInt64{Valid: false}
}

func ToSqlNullInt32(input *int32) sql.NullInt32 {
	if input != nil {
		return sql.NullInt32{Int32: *input, Valid: true}
	}
	return sql.NullInt32{Valid: false}
}

func ToSqlNullInt16(input *int16) sql.NullInt16 {
	if input != nil {
		return sql.NullInt16{Int16: *input, Valid: true}
	}
	return sql.NullInt16{Valid: false}
}

func ToSqlNullFloat64(input *float64) sql.NullFloat64 {
	if input != nil {
		return sql.NullFloat64{Float64: *input, Valid: true}
	}
	return sql.NullFloat64{Valid: false}
}

func ToSqlNullBool(input *bool) sql.NullBool {
	if input != nil {
		return sql.NullBool{Bool: *input, Valid: true}
	}
	return sql.NullBool{Valid: false}
}

func ToSqlNullTime(input *time.Time) sql.NullTime {
	if input != nil {
		return sql.NullTime{Time: *input, Valid: true}
	}
	return sql.NullTime{Valid: false}
}

func ToSqlNullByte(input *[]byte) sql.NullByte {
	if input != nil && len(*input) > 0 {
		return sql.NullByte{Byte: (*input)[0], Valid: true}
	}
	return sql.NullByte{Valid: false}
}
