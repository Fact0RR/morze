package database

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/Fact0RR/morze/internal/domain"
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
	errInsertOrUpdateMorze = errors.New("ошибка при обновлении конфига")
	errGetMorzeMessages    = errors.New("ошибка при получении сообщений")
)

// wrapError оборачивает ошибку с дополнительным контекстом.
func wrapError(errBase error, err error) error {
	return fmt.Errorf("%w: %w", errBase, err)
}

type MorzeRepo struct {
	db     *pgxpool.Pool
	logger *log.Logger
}

func NewMorzeRepo(db *pgxpool.Pool, logger *log.Logger) *MorzeRepo {
	return &MorzeRepo{
		db:     db,
		logger: logger,
	}
}

// GetPrivateMessages получает последние сообщения для указанного contact_id с пагинацией
func (r *MorzeRepo) GetPrivateMessages(ctx context.Context, contactID int, limit int, offset int) ([]domain.MorzeMessage, error) {
	r.logger.WithFields(log.Fields{
		"contact_id": contactID,
		"limit":      limit,
		"offset":     offset,
	}).Debug("Запрос на получение сообщений")

	rows, err := r.db.Query(ctx, queryGetPrivateMessages, contactID, limit, offset)
	if err != nil {
		r.logger.Errorf("Не удалось выполнить запрос: %v, contact_id: %d", err, contactID)
		return nil, wrapError(errGetMorzeMessages, err)
	}
	defer rows.Close()

	var messages []domain.MorzeMessage
	for rows.Next() {
		var msg domain.MorzeMessage
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
			return nil, wrapError(errGetMorzeMessages, err)
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		r.logger.Errorf("Ошибка при итерации по строкам: %v", err)
		return nil, wrapError(errGetMorzeMessages, err)
	}

	r.logger.WithFields(log.Fields{
		"count":      len(messages),
		"contact_id": contactID,
	}).Debug("Сообщения успешно получены")

	return messages, nil
}

// PostPrivateMessage добавляет новое сообщение
func (r *MorzeRepo) PostPrivateMessage(ctx context.Context, contactID int, userID int, data string, additionals []string) (int, error) {
	r.logger.WithFields(log.Fields{
		"contact_id": contactID,
		"user_id":    userID,
	}).Debug("Запрос на добавление сообщения")

	var messageID int
	err := r.db.QueryRow(ctx, queryPostPrivateMessage, contactID, userID, data, additionals).Scan(&messageID)
	if err != nil {
		r.logger.Errorf("Не удалось добавить сообщение: %v", err)
		return 0, wrapError(errInsertOrUpdateMorze, err)
	}

	r.logger.WithField("message_id", messageID).Debug("Сообщение успешно добавлено")
	return messageID, nil
}