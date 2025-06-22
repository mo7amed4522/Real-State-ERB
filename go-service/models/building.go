package models

import (
	"time"
)

type Building struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	Address      string    `json:"address"`
	City         string    `json:"city"`
	Region       string    `json:"region"`
	Price        float64   `json:"price"`
	Status       string    `json:"status"`
	SoldAt       *time.Time `json:"sold_at"`
	CompanyID    uint      `json:"company_id"`
	DeveloperID  uint      `json:"developer_id"`
	TotalLikes   int       `json:"total_likes"`
	TotalComments int      `json:"total_comments"`
	TotalViews   int       `json:"total_views"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
} 