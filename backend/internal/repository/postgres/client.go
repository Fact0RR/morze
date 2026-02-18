package database

import (
	"context"

	"github.com/Fact0RR/morze/internal/configs"
	"github.com/Fact0RR/morze/pkg/handlers"
	"github.com/Fact0RR/morze/pkg/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

type ClientDB struct {
	DB *pgxpool.Pool
}

func InitDB(ctx context.Context, settings *configs.Settings, logger *log.Logger) *ClientDB {
	var err error

	conf := postgres.PGConfig{
		PGDatabaseURL:             settings.DatabaseURL,
		MaxIdlePgConnections:      settings.MaxIdlePgConnections,
		IdleLifetimePgConnections: settings.IdleLifetimePgConnections,
		QueryExecMode:             settings.QueryExecMode,
		WithOpenTelemetry:         settings.EnabledOpentelemetry,
		RetriesLeft:               settings.DBConnectionRetries,
	}
	db, err := postgres.GetPGPool(ctx, conf, logger)
	handlers.ExitIfError(ctx, err, logger)
	logger.Debug("ConfigsgoDB connected")

	return &ClientDB{
		DB: db,
	}
}

func (client *ClientDB) CloseDB(ctx context.Context) {
	if client.DB != nil {
		client.DB.Close()
	}
}
