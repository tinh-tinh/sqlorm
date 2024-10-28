package sqlorm

import "gorm.io/gorm"

type Where interface {
	string | map[string]interface{}
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

func (repo *Repository[M]) FindAll(where interface{}, options FindOptions) ([]M, error) {
	var model []M
	tx := repo.DB
	if len(options.Select) > 0 && options.Select != nil {
		tx = tx.Select(options.Select)
	}
	if len(options.Order) > 0 && options.Order != nil {
		for _, order := range options.Order {
			tx = tx.Order(order)
		}
	}
	if options.Limit > 0 {
		tx = tx.Limit(options.Limit)
	}
	if options.Offset > 0 {
		tx = tx.Offset(options.Offset)
	}
	result := tx.Where(where).Find(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return model, nil
}

func (repo *Repository[M]) FindOne(where interface{}, options FindOneOptions) (*M, error) {
	var model M
	tx := repo.DB
	if options.Select != nil {
		tx = tx.Select(options.Select)
	}
	if options.Order != nil {
		for _, order := range options.Order {
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
