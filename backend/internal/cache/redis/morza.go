package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

type MorzaCache struct {
	cache  *TypeRedisClient
	ttl    int
	enable bool
	logger *log.Logger
}

// NewMorzaCache создает новый экземпляр MorzaCache.
func NewMorzaCache(redisClient *TypeRedisClient, ttl int, enable bool, logger *log.Logger) *MorzaCache {
	return &MorzaCache{
		cache:  redisClient,
		ttl:    ttl,
		enable: enable,
		logger: logger,
	}
}

// TryGetByServiceName пытается получить конфигурацию по имени сервиса из кэша.
func (c *MorzaCache) TryGetByServiceName(ctx context.Context, serviceName string) ([]byte, error) {
	if !c.enable {
		return nil, nil
	}
	// Проверяем готовность Redis
	if !c.cache.Ready(ctx, c.logger) {
		return nil, fmt.Errorf("redis is not ready")
	}

	// Получаем данные из Redis используя serviceName как ключ
	data, err := c.cache.Rdb.Conn().Get(ctx, serviceName).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		c.logger.Error("Failed to get from cache",
			"serviceName", serviceName,
			"error", err.Error())
		return nil, fmt.Errorf("failed to get from cache: %w", err)
	}

	c.logger.Debug("Successfully retrieved from cache",
		"serviceName", serviceName,
		"dataLength", len(data))

	return data, nil
}

// WarmByServiceName сохраняет конфигурацию в кэш по имени сервиса.
func (c *MorzaCache) WarmByServiceName(ctx context.Context, serviceName string, data []byte) error {
	if !c.enable {
		return nil
	}
	// Проверяем готовность Redis
	if !c.cache.Ready(ctx, c.logger) {
		return fmt.Errorf("redis is not ready")
	}

	ttl := time.Duration(c.ttl) * time.Second
	err := c.cache.Rdb.Conn().Set(ctx, serviceName, data, ttl).Err()
	if err != nil {
		c.logger.Error("Failed to warm cache",
			"serviceName", serviceName,
			"error", err.Error())
		return fmt.Errorf("failed to warm cache: %w", err)
	}

	c.logger.Debug("Successfully warmed cache",
		"serviceName", serviceName,
		"dataLength", len(data),
		"ttl", ttl)

	return nil
}

// CoolByServiceName удаляет конфигурацию из кэша по имени сервиса.
func (c *MorzaCache) CoolByServiceName(ctx context.Context, serviceName string) error {
	if !c.enable {
		return nil
	}
	// Проверяем готовность Redis
	if !c.cache.Ready(ctx, c.logger) {
		return fmt.Errorf("redis is not ready")
	}

	// Удаляем ключ из Redis
	err := c.cache.Rdb.Conn().Del(ctx, serviceName).Err()
	if err != nil {
		c.logger.Error("Failed to cool cache by service name",
			"serviceName", serviceName,
			"error", err.Error())
		return fmt.Errorf("failed to cool cache: %w", err)
	}

	c.logger.Debug("Successfully cooled cache by service name",
		"serviceName", serviceName)

	return nil
}

// CoolAll удаляет все конфигурации из кэша.
func (c *MorzaCache) CoolAll(ctx context.Context) error {
	if !c.enable {
		return nil
	}
	// Проверяем готовность Redis
	if !c.cache.Ready(ctx, c.logger) {
		return fmt.Errorf("redis is not ready")
	}

	// Очищаем всю базу Redis
	err := c.cache.Rdb.Conn().FlushDB(ctx).Err()
	if err != nil {
		c.logger.Error("Failed to cool all cache",
			"error", err.Error())
		return fmt.Errorf("failed to cool all cache: %w", err)
	}

	c.logger.Debug("Successfully cooled all cache")

	return nil
}
