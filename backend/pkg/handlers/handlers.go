package handlers

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
)

func ExitIfError(ctx context.Context, err error, logger *log.Logger) {
	if err != nil {
		logger.Error("fatal error", err)
		os.Exit(1)
	}
}
