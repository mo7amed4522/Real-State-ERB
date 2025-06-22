package models

import (
	"time"

	"gorm.io/gorm"
)

// Property represents a real estate property.
type Property struct {
	gorm.Model
	Title             string `gorm:"type:jsonb"`
	Description       string `gorm:"type:jsonb"`
	Type              string
	Status            string
	Price             float64
	Currency          string
	Bedrooms          int
	Bathrooms         int
	Area              float64
	Furnished         bool
	Amenities         string // Comma-separated list
	Images            string // Comma-separated list of URLs
	IsFeatured        bool
	FavoritesCount    int `gorm:"default:0"`
	Views             int `gorm:"default:0"`
	LandlordId        string
	ListedBy          string
	AvailabilityDate  *time.Time
	IsVerified        bool
	OwnershipType     string
	RentalPeriod      string
	DepositRequired   bool
	DepositAmount     float64
	CommissionAmount  float64
	GoogleMapsLink    string
	Neighborhood      string
	NearbyLandmarks   string
	FloorNumber       int
	BuildingName      string
	YearBuilt         int
	ParkingSpaces     int
	Balcony           bool
	Elevator          bool
	MaintenanceFee    float64
	FloorPlan         string
	VideoTourUrl      string
	VirtualTour360Url string
	FloorPlanUrl      string
	Documents         string
	InquiryCount      int `gorm:"default:0"`
	LastViewedAt      *time.Time
	BoostedUntil      *time.Time
	Rating            float64
	Tags              string
	Comments          []Comment `gorm:"foreignKey:PropertyID"`
}
