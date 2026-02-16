package database

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/Fact0RR/morza/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

var (
	//go:embed queries/get_private_messages.sql
	queryGetPrivateMessages string

	//go:embed queries/post_private_message.sql
	queryPostPrivateMessage string
)

var (
	errInsertOrUpdateMorza = errors.New("ошибка при обновлении конфига")
	errGetMorzaMessages    = errors.New("ошибка при получении сообщений")
)

// wrapError оборачивает ошибку с дополнительным контекстом.
func wrapError(errBase error, err error) error {
	return fmt.Errorf("%w: %w", errBase, err)
}

type MorzaRepo struct {
	db     *pgxpool.Pool
	logger *log.Logger
}

func NewMorzaRepo(db *pgxpool.Pool, logger *log.Logger) *MorzaRepo {
	return &MorzaRepo{
		db:     db,
		logger: logger,
	}
}

// GetPrivateMessages получает последние сообщения для указанного contact_id с пагинацией
func (r *MorzaRepo) GetPrivateMessages(ctx context.Context, contactID int, limit int, offset int) ([]domain.MorzaMessage, error) {
	r.logger.WithFields(log.Fields{
		"contact_id": contactID,
		"limit":      limit,
		"offset":     offset,
	}).Debug("Запрос на получение сообщений")

	rows, err := r.db.Query(ctx, queryGetPrivateMessages, contactID, limit, offset)
	if err != nil {
		r.logger.Errorf("Не удалось выполнить запрос: %v, contact_id: %d", err, contactID)
		return nil, wrapError(errGetMorzaMessages, err)
	}
	defer rows.Close()

	var messages []domain.MorzaMessage
	for rows.Next() {
		var msg domain.MorzaMessage
		err := rows.Scan(
			&msg.ID,
			&msg.ContactID,
			&msg.UserID,
			&msg.Data,
			&msg.Additionals,
			&msg.CreatedAt,
		)
		if err != nil {
			r.logger.Errorf("Не удалось прочитать строку: %v", err)
			return nil, wrapError(errGetMorzaMessages, err)
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		r.logger.Errorf("Ошибка при итерации по строкам: %v", err)
		return nil, wrapError(errGetMorzaMessages, err)
	}

	r.logger.WithFields(log.Fields{
		"count":      len(messages),
		"contact_id": contactID,
	}).Debug("Сообщения успешно получены")

	return messages, nil
}

// PostPrivateMessage добавляет новое сообщение
func (r *MorzaRepo) PostPrivateMessage(ctx context.Context, contactID int, userID int, data string, additionals []string) (int, error) {
	r.logger.WithFields(log.Fields{
		"contact_id": contactID,
		"user_id":    userID,
	}).Debug("Запрос на добавление сообщения")

	var messageID int
	err := r.db.QueryRow(ctx, queryPostPrivateMessage, contactID, userID, data, additionals).Scan(&messageID)
	if err != nil {
		r.logger.Errorf("Не удалось добавить сообщение: %v", err)
		return 0, wrapError(errInsertOrUpdateMorza, err)
	}

	r.logger.WithField("message_id", messageID).Debug("Сообщение успешно добавлено")
	return messageID, nil
}