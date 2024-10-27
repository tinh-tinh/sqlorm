package sqlorm

import (
	"reflect"

	"gorm.io/gorm"
)

type Repository[M any] struct {
	DB *gorm.DB
}

func (repo *Repository[M]) Create(val interface{}) (*M, error) {
	input := MapOne[M](val)
	result := repo.DB.Create(input)
	if result.Error != nil {
		return nil, result.Error
	}
	return input, nil
}

func (repo *Repository[M]) BatchCreate(val interface{}) ([]*M, error) {
	input := MapMany[M](val)
	result := repo.DB.Create(input)
	if result.Error != nil {
		return nil, result.Error
	}
	return input, nil
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
