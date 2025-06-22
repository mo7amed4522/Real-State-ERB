package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

type KafkaService struct {
	producer sarama.SyncProducer
	consumer sarama.Consumer
	brokers  []string
	mu       sync.RWMutex
	handlers map[string][]MessageHandler
}

type MessageHandler func(message *ChatMessage) error

type ChatMessage struct {
	ID          uint      `json:"id"`
	RoomID      uint      `json:"room_id"`
	SenderID    uint      `json:"sender_id"`
	SenderType  string    `json:"sender_type"`
	Content     string    `json:"content"`
	MessageType string    `json:"message_type"`
	ReplyToID   *uint     `json:"reply_to_id"`
	ReferenceID *uint     `json:"reference_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type KafkaMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func NewKafkaService(brokers []string) (*KafkaService, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %v", err)
	}

	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %v", err)
	}

	return &KafkaService{
		producer: producer,
		consumer: consumer,
		brokers:  brokers,
		handlers: make(map[string][]MessageHandler),
	}, nil
}

// PublishMessage publishes a message to a specific topic
func (k *KafkaService) PublishMessage(topic string, message *ChatMessage) error {
	msg := KafkaMessage{
		Type:    "chat_message",
		Payload: message,
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	producerMessage := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(jsonData),
	}

	partition, offset, err := k.producer.SendMessage(producerMessage)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	log.Printf("Message sent to topic %s, partition %d, offset %d", topic, partition, offset)
	return nil
}

// SubscribeToTopic subscribes to a topic and handles messages
func (k *KafkaService) SubscribeToTopic(topic string, handler MessageHandler) error {
	partitionConsumer, err := k.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return fmt.Errorf("failed to create partition consumer: %v", err)
	}

	k.mu.Lock()
	k.handlers[topic] = append(k.handlers[topic], handler)
	k.mu.Unlock()

	go func() {
		defer partitionConsumer.Close()
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				var kafkaMsg KafkaMessage
				if err := json.Unmarshal(msg.Value, &kafkaMsg); err != nil {
					log.Printf("Failed to unmarshal message: %v", err)
					continue
				}

				if kafkaMsg.Type == "chat_message" {
					chatMsg, ok := kafkaMsg.Payload.(*ChatMessage)
					if !ok {
						// Try to unmarshal from JSON
						jsonData, _ := json.Marshal(kafkaMsg.Payload)
						if err := json.Unmarshal(jsonData, &chatMsg); err != nil {
							log.Printf("Failed to parse chat message: %v", err)
							continue
						}
					}

					k.mu.RLock()
					handlers := k.handlers[topic]
					k.mu.RUnlock()

					for _, h := range handlers {
						if err := h(chatMsg); err != nil {
							log.Printf("Handler error: %v", err)
						}
					}
				}
			case err := <-partitionConsumer.Errors():
				log.Printf("Consumer error: %v", err)
			}
		}
	}()

	return nil
}

// CreateTopic creates a new Kafka topic
func (k *KafkaService) CreateTopic(topic string) error {
	admin, err := sarama.NewClusterAdmin(k.brokers, nil)
	if err != nil {
		return fmt.Errorf("failed to create cluster admin: %v", err)
	}
	defer admin.Close()

	err = admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}, false)
	if err != nil {
		return fmt.Errorf("failed to create topic: %v", err)
	}

	return nil
}

// Close closes the Kafka service
func (k *KafkaService) Close() error {
	if err := k.producer.Close(); err != nil {
		return fmt.Errorf("failed to close producer: %v", err)
	}
	if err := k.consumer.Close(); err != nil {
		return fmt.Errorf("failed to close consumer: %v", err)
	}
	return nil
}

// PublishRoomMessage publishes a message to a specific room
func (k *KafkaService) PublishRoomMessage(roomID uint, message *ChatMessage) error {
	topic := fmt.Sprintf("chat_room_%d", roomID)
	return k.PublishMessage(topic, message)
}

// SubscribeToRoom subscribes to messages from a specific room
func (k *KafkaService) SubscribeToRoom(roomID uint, handler MessageHandler) error {
	topic := fmt.Sprintf("chat_room_%d", roomID)
	
	// Create topic if it doesn't exist
	if err := k.CreateTopic(topic); err != nil {
		log.Printf("Topic creation failed (might already exist): %v", err)
	}
	
	return k.SubscribeToTopic(topic, handler)
} 