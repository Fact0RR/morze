package postgres

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

type PGConfig struct {
	PGDatabaseURL             string
	MaxIdlePgConnections      int32
	IdleLifetimePgConnections int
	QueryExecMode             int32
	WithOpenTelemetry         bool
	RetriesLeft               int
}

func OpenPgPool(ctx context.Context, dbConf PGConfig, logger *log.Logger) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dbConf.PGDatabaseURL)
	if err != nil {
		logger.Error(InvalidConnectionErrMsg, err)

		return nil, fmt.Errorf(ErrorFmtTemplate, InvalidConnectionErrMsg, err)
	}

	cfg.MaxConns = dbConf.MaxIdlePgConnections
	cfg.MaxConnIdleTime = time.Duration(dbConf.IdleLifetimePgConnections) * time.Second
	cfg.ConnConfig.DefaultQueryExecMode = pgx.QueryExecMode(dbConf.QueryExecMode)

	if dbConf.WithOpenTelemetry {
		cfg.ConnConfig.Tracer = otelpgx.NewTracer()
	}

	dbPool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		logger.Error(UnablePoolErrMsg, err)

		return nil, fmt.Errorf(ErrorFmtTemplate, UnablePoolErrMsg, err)
	}

	return dbPool, nil
}

func NewPgPool(
	ctx context.Context,
	dbConf PGConfig,
	retriesLeft int,
	logger *log.Logger,
) (*pgxpool.Pool, error) {
	var (
		pgpool          *pgxpool.Pool
		err             error
		retriesInterval = 10 * time.Second
	)

	pgpool, err = OpenPgPool(ctx, dbConf, logger)
	if err != nil {
		logger.Error("cannot open connection to postgres", err)

		return nil, err
	}

	err = pgpool.Ping(ctx)
	if err != nil {
		if retriesLeft > 0 {
			logger.Error(
				"cannot talk to postgres, retrying in 10 seconds",
				"attempts left", retriesLeft,
				err,
			)
			time.Sleep(retriesInterval)

			return NewPgPool(ctx, dbConf, retriesLeft-1, logger)
		}

		logger.Error(CannotTalkPostgresErrMsg,
			err,
		)

		return nil, fmt.Errorf(ErrorFmtTemplate, CannotTalkPostgresErrMsg, err)
	}

	return pgpool, nil
}

func GetPGPool(
	ctx context.Context,
	dbConf PGConfig,
	logger *log.Logger,
) (*pgxpool.Pool, error) {
	var (
		pgInstance *pgxpool.Pool
		pgOnce     sync.Once
		err        error
	)

	pgOnce.Do(func() {
		pgInstance, err = NewPgPool(ctx, dbConf, dbConf.RetriesLeft, logger)
	})

	if err != nil {
		return nil, fmt.Errorf(ErrorFmtTemplate, CannotGetPGPoolErrMsg, err)
	}

	return pgInstance, nil
}
