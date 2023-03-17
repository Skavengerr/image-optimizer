package queue

import (
	"context"

	"github.com/Skavengerr/image-optimizer/model"
)

type MessageQueue interface {
	SendMessage(ctx context.Context, message *model.Message) error
}
