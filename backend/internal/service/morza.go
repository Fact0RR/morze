package service

import (
	"context"

	"github.com/Fact0RR/morza/internal/domain"
	log "github.com/sirupsen/logrus"
)

var ValueTypes map[int]string

func init() {
	ValueTypes = make(map[int]string)
	ValueTypes[0] = "int"
	ValueTypes[1] = "str"
	ValueTypes[2] = "double"
	ValueTypes[3] = "boolean"
}

type MorzaService struct {
	repo   domain.MorzaRepository
	cache  domain.MorzaCache
	logger *log.Logger
}

func NewMorzaService(configRepo domain.MorzaRepository, configCache domain.MorzaCache, logger *log.Logger) *MorzaService {
	return &MorzaService{
		repo:   configRepo,
		cache:  configCache,
		logger: logger,
	}
}

func (s *MorzaService) CoolCacheForMorza(ctx context.Context, serviceName string) error {
	return s.cache.CoolByServiceName(ctx, serviceName)
}

func (s *MorzaService) GoToCacheForMorza(ctx context.Context, serviceName string) ([]byte, error) {
	return s.cache.TryGetByServiceName(ctx, serviceName)
}

func (s *MorzaService) WarmCacheByNewMorza(ctx context.Context, serviceName string, data []byte) error {
	return s.cache.WarmByServiceName(ctx, serviceName, data)
}
