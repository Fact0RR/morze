package domain

import "time"

// MorzeMessage представляет структуру сообщения
type MorzeMessage struct {
	ID          int       `json:"id"`
	ContactID   int       `json:"contact_id"`
	UserID      int       `json:"user_id"`
	Data        string    `json:"data"`
	Additionals []string  `json:"additionals"`
	CreatedAt   time.Time `json:"created_at"`
}
