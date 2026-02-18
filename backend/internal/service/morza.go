package service

import (
	"context"

	"github.com/Fact0RR/morze/internal/domain"
	log "github.com/sirupsen/logrus"
)

type MorzeService struct {
	repo   domain.MorzeRepository
	logger *log.Logger
}

func NewMorzeService(configRepo domain.MorzeRepository, logger *log.Logger) *MorzeService {
	return &MorzeService{
		repo:   configRepo,
		logger: logger,
	}
}

func (s *MorzeService) GetPrivateMessages(ctx context.Context, contactID int, limit int, offset int) ([]domain.MorzeMessage, error) {
	return s.repo.GetPrivateMessages(ctx, contactID, limit, offset)
}

func (s *MorzeService) PostPrivateMessages(ctx context.Context, contactID int, userID int, data string, additionals []string) (int, error) {
	return s.repo.PostPrivateMessage(ctx, contactID, userID, data, additionals)
}
