package sqlorm

import (
	"reflect"
	"strings"

	"github.com/tinh-tinh/tinhtinh/v2/common"
	"gorm.io/gorm"
)

type RepoCommon interface {
	SetDB(db *gorm.DB)
	GetName() string
}

func NewRepo[M any](model M) *Repository[M] {
	return &Repository[M]{}
}

type Repository[M any] struct {
	DB *gorm.DB
}

func (r *Repository[M]) GetName() string {
	var model M
	return common.GetStructName(model)
}

func (r *Repository[M]) SetDB(db *gorm.DB) {
	r.DB = db
}

func MapOne[M any](data interface{}) *M {
	var model M
	if data == nil {
		return &model
	}
	ct := reflect.ValueOf(&model).Elem()
	for i := 0; i < ct.NumField(); i++ {
		key := ct.Type().Field(i).Name

		if reflect.TypeOf(data).Kind() == reflect.Pointer {
			ctData := reflect.ValueOf(data).Elem()
			if ctData.FieldByName(key).IsValid() {
				value := ctData.FieldByName(key).Interface()
				if value != nil {
					ct.Field(i).Set(reflect.ValueOf(value))
				}
			}
		} else if reflect.TypeOf(data).Kind() == reflect.Map {
			mapper, ok := data.(map[string]interface{})
			if ok {
				val := mapper[key]
				if val == nil {
					val = mapper[strings.ToLower(key)]
				}
				if val != nil {
					ct.Field(i).Set(reflect.ValueOf(val))
				}
			}
		} else {
			continue
		}

	}

	return &model
}

func MapMany[M any](data interface{}) []*M {
	var models []*M
	if data == nil {
		return models
	}
	if reflect.TypeOf(data).Kind() != reflect.Slice {
		return models
	}
	arrVal := reflect.ValueOf(data)
	if arrVal.IsValid() {
		for i := 0; i < arrVal.Len(); i++ {
			models = append(models, MapOne[M](arrVal.Index(i).Interface()))
		}
	}
	return models
}
