package repository

import "gorm.io/gorm"

type AdminFilter struct {
	ID    *int
	Email *string
}

func (f *AdminFilter) GenerateMods(db *gorm.DB) *gorm.DB {
	if f.ID != nil {
		db = db.Where("id = ?", *f.ID)
	}
	if f.Email != nil {
		db = db.Where("email = ?", *f.Email)
	}
	return db
}
