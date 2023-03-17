package repository

import (
	"context"

	"github.com/Skavengerr/image-optimizer/model"
)

type MessageRepository interface {
	SaveMessage(ctx context.Context, message *model.Message) error
	DeleteMessage(ctx context.Context, messageID string) error
}
