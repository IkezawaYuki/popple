package repository

import "gorm.io/gorm"

type CustomerFilter struct {
	ID              *int
	Email           *string
	FacebookToken   *string
	IsFacebookToken *bool
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
	if c.IsFacebookToken != nil {
		if *c.IsFacebookToken {
			db = db.Where("is_facebook_token != ''")
		} else {
			db = db.Where("is_facebook_token == ''")
		}
	}
	return db
}
