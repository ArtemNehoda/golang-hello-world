package services

import "github.com/ArtemNehoda/golang-hello-world/internal/domain/message"

type messageService struct {
	repo MessageRepository
}

type MessageRepository interface {
	GetAllMessages() ([]message.Entity, error)
}

func NewMessageService(repo MessageRepository) *messageService {
	return &messageService{repo: repo}
}

func (s *messageService) GetAllMessages() ([]message.Entity, error) {
	return s.repo.GetAllMessages()
}
