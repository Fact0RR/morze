package database

import (
	_ "embed"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

var errInsertOrUpdateMorza = errors.New("ошибка при обновлении конфига")

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
