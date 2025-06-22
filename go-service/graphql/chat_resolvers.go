package graphql

import (
	"context"
	"mime/multipart"
	"my-property/go-service/models"
	"my-property/go-service/services"
	"my-property/go-service/utils"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

type ChatResolver struct {
	chatService *services.ChatService
}

func NewChatResolver(chatService *services.ChatService) *ChatResolver {
	return &ChatResolver{
		chatService: chatService,
	}
}

// Room Resolvers
func (r *ChatResolver) CreateRoom(ctx context.Context, input CreateRoomInput) (*models.ChatRoom, error) {
	// Get user ID from context (assuming it's set by middleware)
	userID := ctx.Value("user_id").(uint)

	room, err := r.chatService.CreateRoom(
		input.Name,
		input.Description,
		input.Type,
		userID,
		input.ParticipantIDs,
	)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (r *ChatResolver) GetRoomsByUser(ctx context.Context) ([]*models.ChatRoom, error) {
	userID := ctx.Value("user_id").(uint)

	rooms, err := r.chatService.GetRoomsByUser(userID)
	if err != nil {
		return nil, err
	}

	// Convert to pointers
	var roomPtrs []*models.ChatRoom
	for i := range rooms {
		roomPtrs = append(roomPtrs, &rooms[i])
	}

	return roomPtrs, nil
}

func (r *ChatResolver) GetRoomByID(ctx context.Context, roomID string) (*models.ChatRoom, error) {
	id, err := strconv.ParseUint(roomID, 10, 32)
	if err != nil {
		return nil, err
	}

	room, err := r.chatService.GetRoomByID(uint(id))
	if err != nil {
		return nil, err
	}

	return room, nil
}

// Message Resolvers
func (r *ChatResolver) SendMessage(ctx context.Context, input SendMessageInput) (*models.ChatMessage, error) {
	userID := ctx.Value("user_id").(uint)
	userType := ctx.Value("user_type").(string)

	roomID, err := strconv.ParseUint(input.RoomID, 10, 32)
	if err != nil {
		return nil, err
	}

	var replyToID, referenceID *uint
	if input.ReplyToID != nil {
		id, err := strconv.ParseUint(*input.ReplyToID, 10, 32)
		if err != nil {
			return nil, err
		}
		replyToID = &[]uint{uint(id)}[0]
	}

	if input.ReferenceID != nil {
		id, err := strconv.ParseUint(*input.ReferenceID, 10, 32)
		if err != nil {
			return nil, err
		}
		referenceID = &[]uint{uint(id)}[0]
	}

	message, err := r.chatService.SendMessage(
		uint(roomID),
		userID,
		userType,
		input.Content,
		input.MessageType,
		replyToID,
		referenceID,
	)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (r *ChatResolver) GetMessages(ctx context.Context, roomID string, limit *int, offset *int) ([]*models.ChatMessage, error) {
	id, err := strconv.ParseUint(roomID, 10, 32)
	if err != nil {
		return nil, err
	}

	l := 50 // default limit
	if limit != nil {
		l = *limit
	}

	o := 0 // default offset
	if offset != nil {
		o = *offset
	}

	messages, err := r.chatService.GetMessages(uint(id), l, o)
	if err != nil {
		return nil, err
	}

	// Convert to pointers
	var messagePtrs []*models.ChatMessage
	for i := range messages {
		messagePtrs = append(messagePtrs, &messages[i])
	}

	return messagePtrs, nil
}

func (r *ChatResolver) EditMessage(ctx context.Context, input EditMessageInput) (*models.ChatMessage, error) {
	userID := ctx.Value("user_id").(uint)

	messageID, err := strconv.ParseUint(input.MessageID, 10, 32)
	if err != nil {
		return nil, err
	}

	message, err := r.chatService.EditMessage(uint(messageID), userID, input.NewContent)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (r *ChatResolver) DeleteMessage(ctx context.Context, messageID string) (bool, error) {
	userID := ctx.Value("user_id").(uint)

	id, err := strconv.ParseUint(messageID, 10, 32)
	if err != nil {
		return false, err
	}

	err = r.chatService.DeleteMessage(uint(id), userID)
	if err != nil {
		return false, err
	}

	return true, nil
}

// File Upload Resolvers
func (r *ChatResolver) UploadFile(ctx context.Context, input UploadFileInput) (*models.ChatAttachment, error) {
	messageID, err := strconv.ParseUint(input.MessageID, 10, 32)
	if err != nil {
		return nil, err
	}

	// Manually create a multipart.FileHeader
	fileHeader := &multipart.FileHeader{
		Filename: input.File.Filename,
		Size:     input.File.Size,
		Header:   make(map[string][]string),
	}
	fileHeader.Header.Set("Content-Type", input.File.ContentType)

	attachment, err := r.chatService.UploadFile(uint(messageID), fileHeader, input.File.File)
	if err != nil {
		return nil, err
	}

	return attachment, nil
}

func (r *ChatResolver) DownloadFile(ctx context.Context, attachmentID string) (*FileDownload, error) {
	id, err := strconv.ParseUint(attachmentID, 10, 32)
	if err != nil {
		return nil, err
	}

	data, filename, err := r.chatService.DownloadFile(uint(id))
	if err != nil {
		return nil, err
	}

	return &FileDownload{
		Data:     data,
		Filename: filename,
	}, nil
}

// Reaction Resolvers
func (r *ChatResolver) AddReaction(ctx context.Context, input AddReactionInput) (*models.ChatReaction, error) {
	userID := ctx.Value("user_id").(uint)
	userType := ctx.Value("user_type").(string)

	messageID, err := strconv.ParseUint(input.MessageID, 10, 32)
	if err != nil {
		return nil, err
	}

	reaction, err := r.chatService.AddReaction(uint(messageID), userID, userType, input.Emoji)
	if err != nil {
		return nil, err
	}

	return reaction, nil
}

func (r *ChatResolver) RemoveReaction(ctx context.Context, input RemoveReactionInput) (bool, error) {
	userID := ctx.Value("user_id").(uint)
	userType := ctx.Value("user_type").(string)

	messageID, err := strconv.ParseUint(input.MessageID, 10, 32)
	if err != nil {
		return false, err
	}

	err = r.chatService.RemoveReaction(uint(messageID), userID, userType)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Folder Resolvers
func (r *ChatResolver) CreateFolder(ctx context.Context, input CreateFolderInput) (*models.ChatFolder, error) {
	userID := ctx.Value("user_id").(uint)

	roomID, err := strconv.ParseUint(input.RoomID, 10, 32)
	if err != nil {
		return nil, err
	}

	var parentID *uint
	if input.ParentID != nil {
		id, err := strconv.ParseUint(*input.ParentID, 10, 32)
		if err != nil {
			return nil, err
		}
		parentID = &[]uint{uint(id)}[0]
	}

	folder, err := r.chatService.CreateFolder(uint(roomID), input.Name, input.Description, parentID, userID)
	if err != nil {
		return nil, err
	}

	return folder, nil
}

func (r *ChatResolver) GetFolders(ctx context.Context, roomID string) ([]*models.ChatFolder, error) {
	id, err := strconv.ParseUint(roomID, 10, 32)
	if err != nil {
		return nil, err
	}

	folders, err := r.chatService.GetFolders(uint(id))
	if err != nil {
		return nil, err
	}

	// Convert to pointers
	var folderPtrs []*models.ChatFolder
	for i := range folders {
		folderPtrs = append(folderPtrs, &folders[i])
	}

	return folderPtrs, nil
}

// Notification Resolvers
func (r *ChatResolver) GetNotifications(ctx context.Context) ([]*models.ChatNotification, error) {
	userID := ctx.Value("user_id").(uint)

	notifications, err := r.chatService.GetNotifications(userID)
	if err != nil {
		return nil, err
	}

	// Convert to pointers
	var notificationPtrs []*models.ChatNotification
	for i := range notifications {
		notificationPtrs = append(notificationPtrs, &notifications[i])
	}

	return notificationPtrs, nil
}

func (r *ChatResolver) MarkNotificationAsRead(ctx context.Context, notificationID string) (bool, error) {
	id, err := strconv.ParseUint(notificationID, 10, 32)
	if err != nil {
		return false, err
	}

	err = r.chatService.MarkNotificationAsRead(uint(id))
	if err != nil {
		return false, err
	}

	return true, nil
}

// Search Resolvers
func (r *ChatResolver) SearchMessages(ctx context.Context, roomID string, query string) ([]*models.ChatMessage, error) {
	id, err := strconv.ParseUint(roomID, 10, 32)
	if err != nil {
		return nil, err
	}

	messages, err := r.chatService.SearchMessages(uint(id), query)
	if err != nil {
		return nil, err
	}

	// Convert to pointers
	var messagePtrs []*models.ChatMessage
	for i := range messages {
		messagePtrs = append(messagePtrs, &messages[i])
	}

	return messagePtrs, nil
}

// Statistics Resolvers
func (r *ChatResolver) GetMessageStats(ctx context.Context, roomID string) (*MessageStats, error) {
	id, err := strconv.ParseUint(roomID, 10, 32)
	if err != nil {
		return nil, err
	}

	stats, err := r.chatService.GetMessageStats(uint(id))
	if err != nil {
		return nil, err
	}

	return &MessageStats{
		TotalMessages: int(stats["total_messages"].(int64)),
		TodayMessages: int(stats["today_messages"].(int64)),
		TotalFiles:    int(stats["total_files"].(int64)),
	}, nil
}

// Subscription Resolvers
func (r *ChatResolver) MessageAdded(ctx context.Context, roomID string) (<-chan *models.ChatMessage, error) {
	id, err := strconv.ParseUint(roomID, 10, 32)
	if err != nil {
		return nil, err
	}

	// Create channel for real-time messages
	messageChan := make(chan *models.ChatMessage)

	// Subscribe to Kafka topic for this room
	err = r.chatService.KafkaService().SubscribeToRoom(uint(id), func(message *utils.ChatMessage) error {
		// Convert Kafka message to model
		chatMessage := &models.ChatMessage{
			ID:          message.ID,
			RoomID:      message.RoomID,
			SenderID:    message.SenderID,
			SenderType:  message.SenderType,
			Content:     message.Content,
			MessageType: message.MessageType,
			ReplyToID:   message.ReplyToID,
			ReferenceID: message.ReferenceID,
			CreatedAt:   message.CreatedAt,
			UpdatedAt:   message.CreatedAt,
		}

		// Send to subscription channel
		select {
		case messageChan <- chatMessage:
		default:
			// Channel is full, skip this message
		}

		return nil
	})

	if err != nil {
		close(messageChan)
		return nil, err
	}

	// Clean up when context is cancelled
	go func() {
		<-ctx.Done()
		close(messageChan)
	}()

	return messageChan, nil
}

// Input Types
type CreateRoomInput struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	Type           string `json:"type"`
	ParticipantIDs []uint `json:"participantIds"`
}

type SendMessageInput struct {
	RoomID      string  `json:"roomId"`
	Content     string  `json:"content"`
	MessageType string  `json:"messageType"`
	ReplyToID   *string `json:"replyToId"`
	ReferenceID *string `json:"referenceId"`
}

type EditMessageInput struct {
	MessageID  string `json:"messageId"`
	NewContent string `json:"newContent"`
}

type UploadFileInput struct {
	MessageID string         `json:"messageId"`
	File      graphql.Upload `json:"file"`
}

type AddReactionInput struct {
	MessageID string `json:"messageId"`
	Emoji     string `json:"emoji"`
}

type RemoveReactionInput struct {
	MessageID string `json:"messageId"`
}

type CreateFolderInput struct {
	RoomID      string  `json:"roomId"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ParentID    *string `json:"parentId"`
}

type FileDownload struct {
	Data     []byte `json:"data"`
	Filename string `json:"filename"`
}

type MessageStats struct {
	TotalMessages int `json:"totalMessages"`
	TodayMessages int `json:"todayMessages"`
	TotalFiles    int `json:"totalFiles"`
}
