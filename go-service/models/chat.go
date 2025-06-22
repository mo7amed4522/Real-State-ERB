package models

import (
	"time"
)

// ChatRoom represents a chat room/group
type ChatRoom struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"` // "direct", "group", "company", "developer"
	CreatedBy   uint      `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relationships
	Participants []ChatParticipant `gorm:"foreignKey:RoomID" json:"participants"`
	Messages     []ChatMessage     `gorm:"foreignKey:RoomID" json:"messages"`
}

// ChatParticipant represents a user in a chat room
type ChatParticipant struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	RoomID    uint      `json:"room_id"`
	UserID    uint      `json:"user_id"`
	UserType  string    `json:"user_type"` // "user", "company", "developer"
	Role      string    `json:"role"`      // "admin", "member", "moderator"
	JoinedAt  time.Time `json:"joined_at"`
	IsActive  bool      `json:"is_active"`
	
	// Relationships
	Room ChatRoom `gorm:"foreignKey:RoomID" json:"room"`
}

// ChatMessage represents a message in a chat room
type ChatMessage struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	RoomID      uint      `json:"room_id"`
	SenderID    uint      `json:"sender_id"`
	SenderType  string    `json:"sender_type"` // "user", "company", "developer"
	Content     string    `json:"content"`
	MessageType string    `json:"message_type"` // "text", "file", "image", "emoji"
	ReplyToID   *uint     `json:"reply_to_id"`  // For reply messages
	ReferenceID *uint     `json:"reference_id"` // For referenced messages
	IsEdited    bool      `json:"is_edited"`
	IsDeleted   bool      `json:"is_deleted"`
	IsModerated bool      `json:"is_moderated"` // Whether message was checked by AI
	ModerationStatus string `json:"moderation_status"` // "pending", "approved", "flagged", "blocked"
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relationships
	Room        ChatRoom           `gorm:"foreignKey:RoomID" json:"room"`
	Attachments []ChatAttachment   `gorm:"foreignKey:MessageID" json:"attachments"`
	Reactions   []ChatReaction     `gorm:"foreignKey:MessageID" json:"reactions"`
	ReplyTo     *ChatMessage       `gorm:"foreignKey:ReplyToID" json:"reply_to"`
	Referenced  *ChatMessage       `gorm:"foreignKey:ReferenceID" json:"referenced"`
}

// ChatAttachment represents file attachments in messages
type ChatAttachment struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	MessageID   uint      `json:"message_id"`
	FileName    string    `json:"file_name"`
	FileType    string    `json:"file_type"`
	FileSize    int64     `json:"file_size"`
	EncryptedPath string  `json:"encrypted_path"` // Encrypted file path
	EncryptedKey  string  `json:"encrypted_key"`  // Encrypted file encryption key
	UploadedAt  time.Time `json:"uploaded_at"`
	IsModerated bool      `json:"is_moderated"` // Whether file was checked by AI
	ModerationStatus string `json:"moderation_status"` // "pending", "approved", "flagged", "blocked"
	
	// Relationships
	Message ChatMessage `gorm:"foreignKey:MessageID" json:"message"`
}

// ChatReaction represents emoji reactions to messages
type ChatReaction struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MessageID uint      `json:"message_id"`
	UserID    uint      `json:"user_id"`
	UserType  string    `json:"user_type"`
	Emoji     string    `json:"emoji"`
	CreatedAt time.Time `json:"created_at"`
	
	// Relationships
	Message ChatMessage `gorm:"foreignKey:MessageID" json:"message"`
}

// ChatFolder represents folders for organizing uploaded files
type ChatFolder struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	RoomID      uint      `json:"room_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ParentID    *uint     `json:"parent_id"` // For nested folders
	CreatedBy   uint      `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relationships
	Room     ChatRoom       `gorm:"foreignKey:RoomID" json:"room"`
	Parent   *ChatFolder    `gorm:"foreignKey:ParentID" json:"parent"`
	Children []ChatFolder   `gorm:"foreignKey:ParentID" json:"children"`
	Files    []ChatAttachment `gorm:"foreignKey:FolderID" json:"files"`
}

// ChatNotification represents notifications for chat events
type ChatNotification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	UserType  string    `json:"user_type"`
	RoomID    uint      `json:"room_id"`
	Type      string    `json:"type"` // "message", "mention", "reaction", "file_upload"
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
	
	// Relationships
	Room ChatRoom `gorm:"foreignKey:RoomID" json:"room"`
}

// ChatModerationLog represents moderation events for audit trail
type ChatModerationLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	UserType  string    `json:"user_type"`
	RoomID    uint      `json:"room_id"`
	Content   string    `json:"content"`
	Allowed   bool      `json:"allowed"`
	Reason    string    `json:"reason"`
	Severity  string    `json:"severity"` // "low", "medium", "high"
	Flagged   bool      `json:"flagged"`
	CreatedAt time.Time `json:"created_at"`
} 