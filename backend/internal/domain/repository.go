package domain

import "context"

type MorzaRepository interface{
	GetPrivateMessages(ctx context.Context, contactID int, limit int, offset int) ([]MorzaMessage, error)
	PostPrivateMessage(ctx context.Context, contactID int, userID int, data string, additionals []string) (int, error)
}
