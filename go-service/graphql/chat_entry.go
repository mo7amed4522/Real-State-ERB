package graphql

import (
	"my-property/go-service/database"
	"my-property/go-service/services"
	"my-property/go-service/utils"
)

var ChatService *services.ChatService
var ChatResolverInstance *ChatResolver

func InitChatModule() {
	// Initialize Kafka and Encryption services
	kafkaService, err := utils.NewKafkaService([]string{"kafka:9092"})
	if err != nil {
		panic("Failed to initialize Kafka: " + err.Error())
	}
	encryptionService := utils.NewEncryptionService("your-secret-key-from-env")

	// Initialize ChatService
	ChatService = services.NewChatService(database.DB, kafkaService, encryptionService)
	ChatResolverInstance = NewChatResolver(ChatService)
} 