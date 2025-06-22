package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"my-property/go-service/models"
	"my-property/go-service/utils"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

type ChatService struct {
	db                *gorm.DB
	kafkaService      *utils.KafkaService
	encryptionService *utils.EncryptionService
	uploadDir         string
	aiServiceURL      string
}

// Moderation types
type ModerationRequest struct {
	Content   string `json:"content,omitempty"`
	ImageURL  string `json:"image_url,omitempty"`
	ImagePath string `json:"image_path,omitempty"`
	UserID    uint   `json:"user_id"`
	UserType  string `json:"user_type"`
	RoomID    uint   `json:"room_id"`
}

type ModerationResponse struct {
	Allowed  bool   `json:"allowed"`
	Reason   string `json:"reason,omitempty"`
	Severity string `json:"severity,omitempty"` // "low", "medium", "high"
	Flagged  bool   `json:"flagged"`
}

func NewChatService(db *gorm.DB, kafkaService *utils.KafkaService, encryptionService *utils.EncryptionService) *ChatService {
	uploadDir := "./uploads/chat"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create upload directory: %v", err))
	}

	return &ChatService{
		db:                db,
		kafkaService:      kafkaService,
		encryptionService: encryptionService,
		uploadDir:         uploadDir,
		aiServiceURL:      "http://python-ai-service:8000",
	}
}

func (s *ChatService) KafkaService() *utils.KafkaService {
	return s.kafkaService
}

// Room Management
func (s *ChatService) CreateRoom(name, description, roomType string, createdBy uint, participantIDs []uint) (*models.ChatRoom, error) {
	room := &models.ChatRoom{
		Name:        name,
		Description: description,
		Type:        roomType,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tx := s.db.Begin()
	if err := tx.Create(room).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Add participants
	for _, userID := range participantIDs {
		participant := &models.ChatParticipant{
			RoomID:   room.ID,
			UserID:   userID,
			UserType: "user", // Default to user, can be extended
			Role:     "member",
			JoinedAt: time.Now(),
			IsActive: true,
		}
		if err := tx.Create(participant).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Add creator as admin
	creatorParticipant := &models.ChatParticipant{
		RoomID:   room.ID,
		UserID:   createdBy,
		UserType: "user",
		Role:     "admin",
		JoinedAt: time.Now(),
		IsActive: true,
	}
	if err := tx.Create(creatorParticipant).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return room, nil
}

func (s *ChatService) GetRoomsByUser(userID uint) ([]models.ChatRoom, error) {
	var rooms []models.ChatRoom
	err := s.db.Joins("JOIN chat_participants ON chat_rooms.id = chat_participants.room_id").
		Where("chat_participants.user_id = ? AND chat_participants.is_active = ?", userID, true).
		Preload("Participants").
		Preload("Messages", "is_deleted = ?", false).
		Find(&rooms).Error
	return rooms, err
}

func (s *ChatService) GetRoomByID(roomID uint) (*models.ChatRoom, error) {
	var room models.ChatRoom
	err := s.db.Preload("Participants").
		Preload("Messages", "is_deleted = ?", false).
		Preload("Messages.Attachments").
		Preload("Messages.Reactions").
		Preload("Messages.ReplyTo").
		Preload("Messages.Referenced").
		First(&room, roomID).Error
	return &room, err
}

// Message Management
func (s *ChatService) SendMessage(roomID, senderID uint, senderType, content, messageType string, replyToID, referenceID *uint) (*models.ChatMessage, error) {
	// AI Moderation for text content
	moderationResult, err := s.ModerateMessage(content, "", senderID, senderType, roomID)
	if err != nil {
		// Log moderation error but continue (fail open)
		fmt.Printf("Moderation error: %v\n", err)
		moderationResult = &ModerationResponse{
			Allowed: true,
			Reason:  "Moderation service unavailable",
		}
	}

	// Log moderation event
	s.LogModerationEvent(senderID, senderType, roomID, content, moderationResult)

	// Check if message is allowed
	if !moderationResult.Allowed {
		return nil, fmt.Errorf("message blocked by AI moderation: %s", moderationResult.Reason)
	}

	// Set moderation status
	moderationStatus := "approved"
	if moderationResult.Flagged {
		moderationStatus = "flagged"
	}

	message := &models.ChatMessage{
		RoomID:           roomID,
		SenderID:         senderID,
		SenderType:       senderType,
		Content:          content,
		MessageType:      messageType,
		ReplyToID:        replyToID,
		ReferenceID:      referenceID,
		IsEdited:         false,
		IsDeleted:        false,
		IsModerated:      true,
		ModerationStatus: moderationStatus,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.db.Create(message).Error; err != nil {
		return nil, err
	}

	// Publish to Kafka for real-time updates
	kafkaMessage := &utils.ChatMessage{
		ID:          message.ID,
		RoomID:      message.RoomID,
		SenderID:    message.SenderID,
		SenderType:  message.SenderType,
		Content:     message.Content,
		MessageType: message.MessageType,
		ReplyToID:   message.ReplyToID,
		ReferenceID: message.ReferenceID,
		CreatedAt:   message.CreatedAt,
	}

	if err := s.kafkaService.PublishRoomMessage(roomID, kafkaMessage); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Failed to publish message to Kafka: %v\n", err)
	}

	return message, nil
}

func (s *ChatService) GetMessages(roomID uint, limit, offset int) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	err := s.db.Where("room_id = ? AND is_deleted = ?", roomID, false).
		Preload("Attachments").
		Preload("Reactions").
		Preload("ReplyTo").
		Preload("Referenced").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, err
}

func (s *ChatService) EditMessage(messageID, senderID uint, newContent string) (*models.ChatMessage, error) {
	var message models.ChatMessage
	if err := s.db.First(&message, messageID).Error; err != nil {
		return nil, err
	}

	if message.SenderID != senderID {
		return nil, fmt.Errorf("unauthorized to edit this message")
	}

	message.Content = newContent
	message.IsEdited = true
	message.UpdatedAt = time.Now()

	if err := s.db.Save(&message).Error; err != nil {
		return nil, err
	}

	return &message, nil
}

func (s *ChatService) DeleteMessage(messageID, senderID uint) error {
	var message models.ChatMessage
	if err := s.db.First(&message, messageID).Error; err != nil {
		return err
	}

	if message.SenderID != senderID {
		return fmt.Errorf("unauthorized to delete this message")
	}

	message.IsDeleted = true
	message.UpdatedAt = time.Now()

	return s.db.Save(&message).Error
}

// File Upload Management
func (s *ChatService) UploadFile(messageID uint, fileHeader *multipart.FileHeader, file io.Reader) (*models.ChatAttachment, error) {
	// Create unique filename
	ext := filepath.Ext(fileHeader.Filename)
	filename := fmt.Sprintf("%d_%d%s", messageID, time.Now().Unix(), ext)
	filePath := filepath.Join(s.uploadDir, filename)

	// Save original file temporarily
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return nil, err
	}

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	// Copy file content
	if _, err = io.Copy(dst, file); err != nil {
		return nil, err
	}

	// Get message to extract user info for moderation
	var message models.ChatMessage
	if err := s.db.First(&message, messageID).Error; err != nil {
		return nil, err
	}

	// AI Moderation for file content
	moderationResult, err := s.ModerateMessage("", filePath, message.SenderID, message.SenderType, message.RoomID)
	if err != nil {
		// Log moderation error but continue (fail open)
		fmt.Printf("File moderation error: %v\n", err)
		moderationResult = &ModerationResponse{
			Allowed: true,
			Reason:  "File moderation service unavailable",
		}
	}

	// Log moderation event
	s.LogModerationEvent(message.SenderID, message.SenderType, message.RoomID, "File: "+fileHeader.Filename, moderationResult)

	// Check if file is allowed
	if !moderationResult.Allowed {
		// Remove the temporary file
		os.Remove(filePath)
		return nil, fmt.Errorf("file blocked by AI moderation: %s", moderationResult.Reason)
	}

	// Set moderation status
	moderationStatus := "approved"
	if moderationResult.Flagged {
		moderationStatus = "flagged"
	}

	// Encrypt the file
	encryptedPath := filePath + ".enc"
	encryptedKey, err := s.encryptionService.EncryptFile(filePath, encryptedPath)
	if err != nil {
		return nil, err
	}

	// Remove original file
	os.Remove(filePath)

	// Create attachment record
	attachment := &models.ChatAttachment{
		MessageID:        messageID,
		FileName:         fileHeader.Filename,
		FileType:         fileHeader.Header.Get("Content-Type"),
		FileSize:         fileHeader.Size,
		EncryptedPath:    encryptedPath,
		EncryptedKey:     encryptedKey,
		UploadedAt:       time.Now(),
		IsModerated:      true,
		ModerationStatus: moderationStatus,
	}

	if err := s.db.Create(attachment).Error; err != nil {
		return nil, err
	}

	return attachment, nil
}

func (s *ChatService) DownloadFile(attachmentID uint) ([]byte, string, error) {
	var attachment models.ChatAttachment
	if err := s.db.First(&attachment, attachmentID).Error; err != nil {
		return nil, "", err
	}

	// Create temporary file for decryption
	tempPath := filepath.Join(s.uploadDir, "temp_"+filepath.Base(attachment.FileName))
	defer os.Remove(tempPath)

	// Decrypt file
	if err := s.encryptionService.DecryptFile(attachment.EncryptedPath, tempPath, attachment.EncryptedKey); err != nil {
		return nil, "", err
	}

	// Read decrypted file
	data, err := os.ReadFile(tempPath)
	if err != nil {
		return nil, "", err
	}

	return data, attachment.FileName, nil
}

// Reaction Management
func (s *ChatService) AddReaction(messageID, userID uint, userType, emoji string) (*models.ChatReaction, error) {
	// Check if reaction already exists
	var existingReaction models.ChatReaction
	err := s.db.Where("message_id = ? AND user_id = ? AND user_type = ?", messageID, userID, userType).
		First(&existingReaction).Error

	if err == nil {
		// Update existing reaction
		existingReaction.Emoji = emoji
		existingReaction.CreatedAt = time.Now()
		if err := s.db.Save(&existingReaction).Error; err != nil {
			return nil, err
		}
		return &existingReaction, nil
	}

	// Create new reaction
	reaction := &models.ChatReaction{
		MessageID: messageID,
		UserID:    userID,
		UserType:  userType,
		Emoji:     emoji,
		CreatedAt: time.Now(),
	}

	if err := s.db.Create(reaction).Error; err != nil {
		return nil, err
	}

	return reaction, nil
}

func (s *ChatService) RemoveReaction(messageID, userID uint, userType string) error {
	return s.db.Where("message_id = ? AND user_id = ? AND user_type = ?", messageID, userID, userType).
		Delete(&models.ChatReaction{}).Error
}

// Folder Management
func (s *ChatService) CreateFolder(roomID uint, name, description string, parentID *uint, createdBy uint) (*models.ChatFolder, error) {
	folder := &models.ChatFolder{
		RoomID:      roomID,
		Name:        name,
		Description: description,
		ParentID:    parentID,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.Create(folder).Error; err != nil {
		return nil, err
	}

	return folder, nil
}

func (s *ChatService) GetFolders(roomID uint) ([]models.ChatFolder, error) {
	var folders []models.ChatFolder
	err := s.db.Where("room_id = ?", roomID).
		Preload("Children").
		Preload("Files").
		Find(&folders).Error
	return folders, err
}

// Notification Management
func (s *ChatService) CreateNotification(userID uint, userType string, roomID uint, notificationType, message string) (*models.ChatNotification, error) {
	notification := &models.ChatNotification{
		UserID:    userID,
		UserType:  userType,
		RoomID:    roomID,
		Type:      notificationType,
		Message:   message,
		IsRead:    false,
		CreatedAt: time.Now(),
	}

	if err := s.db.Create(notification).Error; err != nil {
		return nil, err
	}

	return notification, nil
}

func (s *ChatService) GetNotifications(userID uint) ([]models.ChatNotification, error) {
	var notifications []models.ChatNotification
	err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

func (s *ChatService) MarkNotificationAsRead(notificationID uint) error {
	return s.db.Model(&models.ChatNotification{}).
		Where("id = ?", notificationID).
		Update("is_read", true).Error
}

// Search functionality
func (s *ChatService) SearchMessages(roomID uint, query string) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	err := s.db.Where("room_id = ? AND content ILIKE ? AND is_deleted = ?", roomID, "%"+query+"%", false).
		Preload("Attachments").
		Preload("Reactions").
		Order("created_at DESC").
		Find(&messages).Error
	return messages, err
}

// Get message statistics
func (s *ChatService) GetMessageStats(roomID uint) (map[string]interface{}, error) {
	var stats map[string]interface{}

	var totalMessages int64
	s.db.Model(&models.ChatMessage{}).Where("room_id = ? AND is_deleted = ?", roomID, false).Count(&totalMessages)

	var todayMessages int64
	today := time.Now().Truncate(24 * time.Hour)
	s.db.Model(&models.ChatMessage{}).Where("room_id = ? AND is_deleted = ? AND created_at >= ?", roomID, false, today).Count(&todayMessages)

	var totalFiles int64
	s.db.Model(&models.ChatAttachment{}).Joins("JOIN chat_messages ON chat_attachments.message_id = chat_messages.id").
		Where("chat_messages.room_id = ? AND chat_messages.is_deleted = ?", roomID, false).Count(&totalFiles)

	stats = map[string]interface{}{
		"total_messages": totalMessages,
		"today_messages": todayMessages,
		"total_files":    totalFiles,
	}

	return stats, nil
}

// AI Moderation
func (s *ChatService) ModerateMessage(content, imagePath string, userID uint, userType string, roomID uint) (*ModerationResponse, error) {
	reqBody := ModerationRequest{
		Content:   content,
		ImagePath: imagePath,
		UserID:    userID,
		UserType:  userType,
		RoomID:    roomID,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal moderation request: %v", err)
	}

	resp, err := http.Post(s.aiServiceURL+"/moderate", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("AI service error: %v", err)
	}
	defer resp.Body.Close()

	var result ModerationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode AI response: %v", err)
	}

	return &result, nil
}

// Log moderation event
func (s *ChatService) LogModerationEvent(userID uint, userType string, roomID uint, content string, moderationResult *ModerationResponse) error {
	logEntry := &models.ChatModerationLog{
		UserID:    userID,
		UserType:  userType,
		RoomID:    roomID,
		Content:   content,
		Allowed:   moderationResult.Allowed,
		Reason:    moderationResult.Reason,
		Severity:  moderationResult.Severity,
		Flagged:   moderationResult.Flagged,
		CreatedAt: time.Now(),
	}

	return s.db.Create(logEntry).Error
}
