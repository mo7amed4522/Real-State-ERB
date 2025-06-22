package models

import (
	"gorm.io/gorm"
)

// Comment represents a comment on a property.
type Comment struct {
	gorm.Model
	PropertyID uint
	Content    string
	UserID     string
	LikesCount int `gorm:"default:0"`
	ParentID   *uint
	Replies    []Comment `gorm:"foreignKey:ParentID"`
}
