package sqlorm

import "gorm.io/gorm"

type Where interface {
	string | map[string]interface{} | []interface{}
}

type FindOneOptions struct {
	Select []string
	Order  []string
}

type FindOptions struct {
	Select []string
	Order  []string
	Limit  int
	Offset int
}

func (repo *Repository[M]) FindAll(where interface{}, options ...FindOptions) ([]M, error) {
	var model []M
	tx := repo.DB

	var opt FindOptions
	if len(options) > 0 {
		opt = options[0]
	}
	if len(opt.Select) > 0 && opt.Select != nil {
		tx = tx.Select(opt.Select)
	}
	if len(opt.Order) > 0 && opt.Order != nil {
		for _, order := range opt.Order {
			tx = tx.Order(order)
		}
	}
	if opt.Limit > 0 {
		tx = tx.Limit(opt.Limit)
	}
	if opt.Offset > 0 {
		tx = tx.Offset(opt.Offset)
	}
	result := tx.Where(where).Find(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return model, nil
}

func (repo *Repository[M]) FindOne(where interface{}, options ...FindOneOptions) (*M, error) {
	var model M
	tx := repo.DB

	var opt FindOneOptions
	if len(options) > 0 {
		opt = options[0]
	}
	if opt.Select != nil {
		tx = tx.Select(opt.Select)
	}
	if opt.Order != nil {
		for _, order := range opt.Order {
			tx = tx.Order(order)
		}
	}
	result := repo.DB.Where(where).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &model, nil
}

func (repo *Repository[M]) FindByID(id string, options FindOneOptions) (*M, error) {
	return repo.FindOne(map[string]interface{}{"id": id}, options)
}

func (repo *Repository[M]) Count(where interface{}) (int64, error) {
	var count int64
	var model M
	result := repo.DB.Model(&model).Where(where).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (repo *Repository[M]) Exist(where interface{}, options ...FindOneOptions) (bool, error) {
	var model M
	tx := repo.DB

	var opt FindOneOptions
	if len(options) > 0 {
		opt = options[0]
	}
	if opt.Select != nil {
		tx = tx.Select(opt.Select)
	}
	if opt.Order != nil {
		for _, order := range opt.Order {
			tx = tx.Order(order)
		}
	}
	result := tx.Where(where).First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}
