package domain

import "time"

// MorzaMessage представляет структуру сообщения
type MorzaMessage struct {
	ID          int       `json:"id"`
	ContactID   int       `json:"contact_id"`
	UserID      int       `json:"user_id"`
	Data        string    `json:"data"`
	Additionals []string  `json:"additionals"`
	CreatedAt   time.Time `json:"created_at"`
}
