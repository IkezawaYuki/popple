package model

import (
	"time"
)

type Post struct {
	ID               int       `gorm:"column:id;primaryKey"`
	CustomerID       int       `gorm:"column:customer_id"`
	InstagramMediaID string    `gorm:"column:instagram_media_id"`
	InstagramLink    string    `gorm:"column:instagram_link"`
	WordpressMediaID string    `gorm:"column:wordpress_media_id"`
	WordpressLink    string    `gorm:"column:wordpress_link"`
	CreatedAt        time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt        time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}
