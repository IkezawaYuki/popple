package repository

import "gorm.io/gorm"

type PostFilter struct {
	ID               *int
	InstagramMediaID *string
	CustomerID       *int
	Limit            *int
	Offset           *int
}

func (p *PostFilter) GenerateMods(db *gorm.DB) *gorm.DB {
	if p.ID != nil {
		db = db.Where("id = ?", *p.ID)
	}
	if p.InstagramMediaID != nil {
		db = db.Where("instagram_media_id = ?", *p.InstagramMediaID)
	}
	if p.CustomerID != nil {
		db = db.Where("customer_id = ?", *p.CustomerID)
	}
	if p.Limit != nil {
		db = db.Limit(*p.Limit)
		if p.Offset != nil {
			db = db.Offset(*p.Offset)
		}
	}
	return db
}
