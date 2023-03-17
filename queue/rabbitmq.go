package queue

import (
	"context"
	"encoding/json"

	"github.com/Skavengerr/image-optimizer/model"
	"github.com/streadway/amqp"
)

type rabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQ(connectionString string) (*rabbitMQ, error) {
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &rabbitMQ{
		conn: conn,
		ch:   ch,
	}, nil
}

func (q *rabbitMQ) SendMessage(ctx context.Context, message *model.Message) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = q.ch.Publish(
		"",
		"messages",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (q *rabbitMQ) Close() error {
	if err := q.ch.Close(); err != nil {
		return err
	}
	if err := q.conn.Close(); err != nil {
		return err
	}
	return nil
}
