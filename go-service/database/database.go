package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"my-property/go-service/models"
)

var DB *gorm.DB

func Connect() {
	var err error
	dsn := os.Getenv("DATABASE_URL")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connection successful.")
}

func InitDB() {
	Connect()
	// Auto-migrate the schema
	DB.AutoMigrate(
		&models.Property{}, 
		&models.PropertyImage{}, 
		&models.Comment{},
		&models.Building{},
		&models.Like{},
		// Chat models
		&models.ChatRoom{},
		&models.ChatParticipant{},
		&models.ChatMessage{},
		&models.ChatAttachment{},
		&models.ChatReaction{},
		&models.ChatFolder{},
		&models.ChatNotification{},
		&models.ChatModerationLog{},
		// Financial models
		&models.SaleTransaction{},
		&models.LeaseContract{},
		&models.LeasePayment{},
		&models.Offer{},
		&models.OfferUse{},
		&models.TransactionDocument{},
		&models.LeaseDocument{},
		&models.FinancialReport{},
	)
	log.Println("Database migrated")
}

// Migration for multilingual support
func MigratePropertyTranslations(db *gorm.DB) error {
	return db.Exec(`ALTER TABLE properties 
	ALTER COLUMN title TYPE JSONB USING title::jsonb,
	ALTER COLUMN description TYPE JSONB USING description::jsonb;`).Error
} 