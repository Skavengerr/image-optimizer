package model

import (
	"errors"
	"time"
)

type Message struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	ImageURL  string    `json:"image_url,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

func (m *Message) Validate() error {
	if m.ID == "" {
		return errors.New("message ID is required")
	}
	if m.Text == "" {
		return errors.New("message text is required")
	}
	if m.Timestamp.IsZero() {
		return errors.New("message timestamp is required")
	}
	return nil
}
