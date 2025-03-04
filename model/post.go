package model

import (
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID         uuid.UUID `json:"id" gorm:"primary_key;type:char(36)"`
	UserId     uint      `json:"user_id" gorm:"not null"`
	CategoryId uint      `json:"category_id" gorm:"not null"`
	Category   *Category
	Title      string `json:"title" gorm:"type:varchar(50);not null"`
	HeadImgUrl string `json:"head_img_url"`
	Content    string `json:"content" gorm:"type:text;not null"`
	CreatedAt  Time   `json:"created_at" gorm:"type: timestamp"`
	UpdatedAt  Time   `json:"updated_at" gorm:"type: timestamp"`
}

func (post *Post) BeforeCreate(tx *gorm.DB) error {
	post.ID = uuid.NewV4()
	return nil
}
