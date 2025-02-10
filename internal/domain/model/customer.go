package model

import (
	"time"
)

type Customer struct {
	ID             int        `gorm:"column:id;primaryKey"`
	Name           string     `gorm:"column:name"`
	Email          string     `gorm:"column:email"`
	Password       string     `gorm:"column:password"`
	WordpressURL   string     `gorm:"column:wordpress_url"`
	FacebookToken  *string    `gorm:"column:facebook_token"`
	StartDate      *time.Time `gorm:"column:start_date"`
	InstagramID    *string    `gorm:"column:instagram_id"`
	InstagramName  *string    `gorm:"column:instagram_name"`
	DeleteHashFlag int        `gorm:"column:delete_hash_flag"`
	CreatedAt      time.Time  `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}
