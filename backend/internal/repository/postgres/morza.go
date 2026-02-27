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

	//go:embed queries/delete_private_message.sql
	queryDeletePrivateMessage string

	//go:embed queries/update_private_message.sql
	queryUpdatePrivateMessage string
)

var (
	errInsertOrUpdateMorze = errors.New("ошибка при обновлении конфига")
	errGetMorzeMessages    = errors.New("ошибка при получении сообщений")
	errDeleteMorzeMessage  = errors.New("ошибка при удалении сообщения")
	errUpdateMorzeMessage  = errors.New("ошибка при обновлении сообщения")
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
			&msg.Updated,
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

// DeletePrivateMessage мягко удаляет сообщение по ID
func (r *MorzeRepo) DeletePrivateMessage(ctx context.Context, messageID int) error {
	r.logger.WithField("message_id", messageID).Debug("Запрос на удаление сообщения")

	commandTag, err := r.db.Exec(ctx, queryDeletePrivateMessage, messageID)
	if err != nil {
		r.logger.Errorf("Не удалось удалить сообщение: %v, message_id: %d", err, messageID)
		return wrapError(errDeleteMorzeMessage, err)
	}

	if commandTag.RowsAffected() == 0 {
		r.logger.WithField("message_id", messageID).Warn("Сообщение не найдено для удаления")
		return wrapError(errDeleteMorzeMessage, fmt.Errorf("сообщение с id %d не найдено", messageID))
	}

	r.logger.WithField("message_id", messageID).Debug("Сообщение успешно удалено")
	return nil
}

// UpdatePrivateMessage обновляет data сообщения по ID
func (r *MorzeRepo) UpdatePrivateMessage(ctx context.Context, messageID int, newData string) error {
	r.logger.WithFields(log.Fields{
		"message_id": messageID,
		"new_data":   newData,
	}).Debug("Запрос на обновление сообщения")

	commandTag, err := r.db.Exec(ctx, queryUpdatePrivateMessage, messageID, newData)
	if err != nil {
		r.logger.Errorf("Не удалось обновить сообщение: %v, message_id: %d", err, messageID)
		return wrapError(errUpdateMorzeMessage, err)
	}

	if commandTag.RowsAffected() == 0 {
		r.logger.WithField("message_id", messageID).Warn("Сообщение не найдено для обновления")
		return wrapError(errUpdateMorzeMessage, fmt.Errorf("сообщение с id %d не найдено", messageID))
	}

	r.logger.WithField("message_id", messageID).Debug("Сообщение успешно обновлено")
	return nil
}
