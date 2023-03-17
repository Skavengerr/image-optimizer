package service

import (
	"encoding/json"

	"github.com/Skavengerr/image-optimizer/model"
	repo "github.com/Skavengerr/image-optimizer/repository"
	"github.com/streadway/amqp"
)

func NewImageService(imageRepo *repo.ImageRepository, rabbitMQConn *amqp.Connection) *ImageService {
	return &ImageService{
		imageRepo:    imageRepo,
		rabbitMQConn: rabbitMQConn,
	}
}

type ImageService struct {
	imageRepo    *repo.ImageRepository
	rabbitMQConn *amqp.Connection
}

func (is *ImageService) Create(image *model.Image) error {
	err := is.imageRepo.Create(image)
	if err != nil {
		return err
	}

	ch, err := is.rabbitMQConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"image_process_queue",
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
