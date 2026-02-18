package domain

import "context"

type MorzeRepository interface{
	GetPrivateMessages(ctx context.Context, contactID int, limit int, offset int) ([]MorzeMessage, error)
	PostPrivateMessage(ctx context.Context, contactID int, userID int, data string, additionals []string) (int, error)
}
