package repository

import (
	"encoding/json"

	"github.com/Skavengerr/image-optimizer/model"
	"github.com/streadway/amqp"
)

func NewImageRepository(rabbitMQConn *amqp.Connection) *ImageRepository {
	return &ImageRepository{
		rabbitMQConn: rabbitMQConn,
		queueName:    "image_queue",
	}
}

type ImageRepository struct {
	rabbitMQConn *amqp.Connection
	queueName    string
}

func (ir *ImageRepository) Create(image *model.Image) error {
	ch, err := ir.rabbitMQConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		ir.queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(image)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
