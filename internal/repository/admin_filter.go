package repository

import "gorm.io/gorm"

type AdminFilter struct {
	ID          *int
	Email       *string
	PartialName *string
	Limit       *int
	Offset      *int
}

func (f *AdminFilter) GenerateMods(db *gorm.DB) *gorm.DB {
	if f.ID != nil {
		db = db.Where("id = ?", *f.ID)
	}
	if f.PartialName != nil {
		db = db.Where("name like '%?%'", *f.PartialName)
	}
	if f.Email != nil {
		db = db.Where("email = ?", *f.Email)
	}
	if f.Limit != nil {
		db = db.Limit(*f.Limit)
		if f.Offset != nil {
			db = db.Offset(*f.Offset)
		}
	}
	return db
}
