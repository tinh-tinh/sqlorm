package sqlorm

import (
	"reflect"
	"strings"

	"gorm.io/gorm"
)

type Query interface {
	func(qb *QueryBuilder) | interface{}
}

func IsQueryBuilder(q Query) bool {
	if q == nil {
		return false
	}
	typeQ := reflect.TypeOf(q)
	return typeQ.Kind() == reflect.Func
}

type QueryBuilder struct {
	qb *gorm.DB
}

func (q *QueryBuilder) Equal(column string, value interface{}) *QueryBuilder {
	query := column + " = ?"
	q.qb = q.qb.Where(query, value)
	return q
}

func (q *QueryBuilder) Not(column string, args ...interface{}) *QueryBuilder {
	query := column + " = ?"
	q.qb = q.qb.Not(query, args...)
	return q
}

func (q *QueryBuilder) Or(column string, args ...interface{}) *QueryBuilder {
	query := column + " = ?"
	q.qb = q.qb.Or(query, args...)
	return q
}

func (q *QueryBuilder) In(column string, values ...interface{}) *QueryBuilder {
	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = "?"
	}
	query := column + " IN (" + strings.Join(placeholders, ", ") + ")"
	q.qb = q.qb.Where(query, values...)
	return q
}

func (q *QueryBuilder) MoreThan(column string, value interface{}) *QueryBuilder {
	query := column + " > ?"
	q.qb = q.qb.Where(query, value)
	return q
}

func (q *QueryBuilder) MoreThanOrEqual(column string, value interface{}) *QueryBuilder {
	query := column + " >= ?"
	q.qb = q.qb.Where(query, value)
	return q
}

func (q *QueryBuilder) LessThan(column string, value interface{}) *QueryBuilder {
	query := column + " < ?"
	q.qb = q.qb.Where(query, value)
	return q
}

func (q *QueryBuilder) LessThanOrEqual(column string, value interface{}) *QueryBuilder {
	query := column + " <= ?"
	q.qb = q.qb.Where(query, value)
	return q
}

func (q *QueryBuilder) Like(column string, value interface{}) *QueryBuilder {
	query := column + " LIKE ?"
	q.qb = q.qb.Where(query, value)
	return q
}

func (q *QueryBuilder) ILike(column string, value interface{}) *QueryBuilder {
	query := column + " ILIKE ?"
	q.qb = q.qb.Where(query, value)
	return q
}

func (q *QueryBuilder) Between(column string, start interface{}, end interface{}) *QueryBuilder {
	query := column + " BETWEEN ? AND ?"
	q.qb = q.qb.Where(query, start, end)
	return q
}

func (q *QueryBuilder) NotEqual(column string, value interface{}) *QueryBuilder {
	query := column + " <> ?"
	q.qb = q.qb.Where(query, value)
	return q
}

func (q *QueryBuilder) NotIn(column string, values ...interface{}) *QueryBuilder {
	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = "?"
	}
	query := column + " NOT IN (" + strings.Join(placeholders, ", ") + ")"
	q.qb = q.qb.Where(query, values...)
	return q
}

func (q *QueryBuilder) IsNull(column string) *QueryBuilder {
	query := column + " IS NULL"
	q.qb = q.qb.Where(query)
	return q
}

func (q *QueryBuilder) Raw(sql string, values ...interface{}) *QueryBuilder {
	q.qb = q.qb.Raw(sql, values...)
	return q
}
