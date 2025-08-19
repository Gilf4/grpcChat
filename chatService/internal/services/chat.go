package services

import "github.com/Gilf4/grpcChat/chat/internal/domain/models"

type ChatRepository interface {
	Create(chat *models.Chat) error
	Get(id string) (*models.Chat, error)
	GetList() ([]*models.Chat, error)
	Delete(id string) error
}

type MessageRepository interface {
	Create(message *models.Message) error
	Get(id string) (*models.Message, error)
	GetList() ([]*models.Message, error)
	Delete(id string) error
}

type ChatService struct {
	chatRepo    ChatRepository
	messageRepo MessageRepository
}

func NewChatService(chatRepo ChatRepository, messageRepo MessageRepository) *ChatService {
	return &ChatService{
		chatRepo:    chatRepo,
		messageRepo: messageRepo,
	}
}
