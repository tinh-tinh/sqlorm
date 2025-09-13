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

func (repo *Repository[M]) BatchCreate(val interface{}, size int) ([]*M, error) {
	input := MapMany[M](val)
	result := repo.DB.CreateInBatches(input, size)
	if result.Error != nil {
		return nil, result.Error
	}
	return input, nil
}

func (repo *Repository[M]) UpdateOne(where interface{}, val interface{}) (*M, error) {
	var record M
	input := MapOne[M](val)
	result := repo.DB.Model(&record).Where(where).Updates(input)
	if result.Error != nil {
		return nil, result.Error
	}
	return input, nil
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

func (repo *Repository[M]) DeleteOne(where interface{}, isForceDelete ...bool) error {
	withDeleted := false
	if len(isForceDelete) > 0 && isForceDelete[0] {
		withDeleted = true
	}

	record, err := repo.FindOne(where, FindOneOptions{
		WithDeleted: withDeleted,
	})
	if err != nil {
		return err
	}
	if record == nil {
		return gorm.ErrRecordNotFound
	}
	tx := repo.DB
	if withDeleted {
		tx = tx.Unscoped()
	}
	result := tx.Delete(record)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *Repository[M]) DeleteByID(id any, isForceDelete ...bool) error {
	return repo.DeleteOne(map[string]any{"id": id}, isForceDelete...)
}

func (repo *Repository[M]) DeleteMany(where interface{}, isForceDelete ...bool) error {
	var model M
	tx := repo.DB
	if where != nil {
		tx = tx.Where(where)
	} else {
		tx = tx.Where("1 = 1")
	}

	if len(isForceDelete) > 0 && isForceDelete[0] {
		tx = tx.Unscoped()
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
