package repository

import "gorm.io/gorm"

type CustomerFilter struct {
	ID              *int
	Email           *string
	PartialName     *string
	FacebookToken   *string
	IsFacebookToken *bool
	Limit           *int
	Offset          *int
}

func (c *CustomerFilter) GenerateMods(db *gorm.DB) *gorm.DB {
	if c.ID != nil {
		db = db.Where("id = ?", *c.ID)
	}
	if c.Email != nil {
		db = db.Where("email = ?", *c.Email)
	}
	if c.FacebookToken != nil {
		db = db.Where("facebook_token = ?", *c.FacebookToken)
	}
	if c.PartialName != nil {
		db = db.Where("partial_name like '%?%'", *c.PartialName)
	}
	if c.IsFacebookToken != nil {
		if *c.IsFacebookToken {
			db = db.Where("is_facebook_token is not null")
		} else {
			db = db.Where("is_facebook_token is null")
		}
	}
	if c.Limit != nil {
		db = db.Limit(*c.Limit)
		if c.Offset != nil {
			db = db.Offset(*c.Offset)
		}
	}
	return db
}
