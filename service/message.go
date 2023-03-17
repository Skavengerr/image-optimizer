package service

import (
	"context"

	"github.com/Skavengerr/image-optimizer/model"
	"github.com/Skavengerr/image-optimizer/queue"
	"github.com/Skavengerr/image-optimizer/repository"
)

type MessageService interface {
	SendMessage(ctx context.Context, message *model.Message) error
}

type messageService struct {
	repo  repository.MessageRepository
	queue queue.MessageQueue
}

func (s *messageService) SendMessage(ctx context.Context, message *model.Message) error {
	if err := message.Validate(); err != nil {
		return err
	}

	if err := s.repo.SaveMessage(ctx, message); err != nil {
		return err
	}

	if err := s.queue.SendMessage(ctx, message); err != nil {
		_ = s.repo.DeleteMessage(ctx, message.ID)
		return err
	}

	return nil
}

func NewMessageService(repo repository.MessageRepository, queue queue.MessageQueue) MessageService {
	return &messageService{
		repo:  repo,
		queue: queue,
	}
}
