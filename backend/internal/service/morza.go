package service

import (
	"context"

	"github.com/Fact0RR/morza/internal/domain"
	log "github.com/sirupsen/logrus"
)

type MorzaService struct {
	repo   domain.MorzaRepository
	logger *log.Logger
}

func NewMorzaService(configRepo domain.MorzaRepository, logger *log.Logger) *MorzaService {
	return &MorzaService{
		repo:   configRepo,
		logger: logger,
	}
}

func (s *MorzaService) GetPrivateMessages(ctx context.Context, contactID int, limit int, offset int) ([]domain.MorzaMessage, error) {
	return s.repo.GetPrivateMessages(ctx, contactID, limit, offset)
}

func (s *MorzaService) PostPrivateMessages(ctx context.Context, contactID int, userID int, data string, additionals []string) (int, error) {
	return s.repo.PostPrivateMessage(ctx, contactID, userID, data, additionals)
}
