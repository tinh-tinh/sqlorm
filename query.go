package sqlorm

import "gorm.io/gorm"

type FindOneOptions struct {
	Select      []string
	Order       []string
	WithDeleted bool
}

type FindOptions struct {
	Select      []string
	Order       []string
	WithDeleted bool
	Limit       int
	Offset      int
}

func (repo *Repository[M]) FindAll(where Query, options ...FindOptions) ([]M, error) {
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
	if opt.WithDeleted {
		tx = tx.Unscoped()
	}

	if IsQueryBuilder(where) {
		queryFnc, ok := where.(func(qb *QueryBuilder))
		if ok {
			qb := &QueryBuilder{qb: tx}
			queryFnc(qb)
			tx = qb.qb
		}
	} else {
		tx = tx.Where(where)
	}

	result := tx.Find(&model)
	if result.Error != nil {
		return nil, result.Error
	}
	return model, nil
}

func (repo *Repository[M]) FindOne(where Query, options ...FindOneOptions) (*M, error) {
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
	if opt.WithDeleted {
		tx = tx.Unscoped()
	}

	if IsQueryBuilder(where) {
		queryFnc, ok := where.(func(qb *QueryBuilder))
		if ok {
			qb := &QueryBuilder{qb: tx}
			queryFnc(qb)
			tx = qb.qb
		}
	} else {
		tx = tx.Where(where)
	}

	result := tx.First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &model, nil
}

func (repo *Repository[M]) FindByID(id string, options ...FindOneOptions) (*M, error) {
	return repo.FindOne(map[string]interface{}{"id": id}, options...)
}

func (repo *Repository[M]) Count(where interface{}, args ...interface{}) (int64, error) {
	var count int64
	var model M
	result := repo.DB.Model(&model).Where(where, args...).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

func (repo *Repository[M]) Exist(where Query, options ...FindOneOptions) (bool, error) {
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
	if opt.WithDeleted {
		tx = tx.Unscoped()
	}

	if IsQueryBuilder(where) {
		queryFnc, ok := where.(func(qb *QueryBuilder))
		if ok {
			qb := &QueryBuilder{qb: tx}
			queryFnc(qb)
			tx = qb.qb
		}
	} else {
		tx = tx.Where(where)
	}

	result := tx.First(&model)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}
