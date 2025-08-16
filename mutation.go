package sqlorm

import "gorm.io/gorm"

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

func (repo *Repository[M]) UpdateOne(where interface{}, val interface{}) (*M, error) {
	record, err := repo.FindOne(where, FindOneOptions{})
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, gorm.ErrRecordNotFound
	}
	input := MapOne[M](val)
	result := repo.DB.Model(record).Updates(input)
	if result.Error != nil {
		return nil, result.Error
	}
	return record, nil
}

func (repo *Repository[M]) UpdateByID(id any, val interface{}) (*M, error) {
	return repo.UpdateOne(map[string]any{"id": id}, val)
}

func (repo *Repository[M]) UpdateMany(where interface{}, val interface{}) error {
	var model M
	input := MapOne[M](val)
	tx := repo.DB.Model(&model)
	if where != nil {
		tx = tx.Where(where)
	} else {
		tx = tx.Where("1 = 1")
	}
	result := tx.Updates(input)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *Repository[M]) DeleteOne(where interface{}) error {
	record, err := repo.FindOne(where, FindOneOptions{})
	if err != nil {
		return err
	}
	if record == nil {
		return gorm.ErrRecordNotFound
	}
	result := repo.DB.Delete(record)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *Repository[M]) DeleteByID(id any) error {
	return repo.DeleteOne(map[string]any{"id": id})
}

func (repo *Repository[M]) DeleteMany(where interface{}) error {
	var model M
	tx := repo.DB
	if where != nil {
		tx = tx.Where(where)
	} else {
		tx = tx.Where("1 = 1")
	}
	result := tx.Delete(&model)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *Repository[M]) Increment(id any, field string, value int) error {
	record, err := repo.FindOne(map[string]interface{}{"id": id}, FindOneOptions{})

	if err != nil {
		return err
	}

	result := repo.DB.Model(record).Update(field, gorm.Expr(field+" + ?", value))
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *Repository[M]) Decrement(id any, field string, value int) error {
	record, err := repo.FindOne(map[string]interface{}{"id": id}, FindOneOptions{})
	if err != nil {
		return err
	}

	result := repo.DB.Model(record).Update(field, gorm.Expr(field+" - ?", value))
	if result.Error != nil {
		return result.Error
	}
	return nil
}
