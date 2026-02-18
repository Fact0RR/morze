package redis

import (
	"context"
	"runtime"
	"sync"

	"github.com/Fact0RR/morze/internal/configs"
	"github.com/gofiber/storage/redis/v3"

	log "github.com/sirupsen/logrus"
)

type TypeRedisClient struct {
	Rdb  *redis.Storage
	mu   sync.Mutex
	done bool
}

func NewTypeRedisClient(settings configs.RedisSettings, logger *log.Logger) *TypeRedisClient {
	client := new(TypeRedisClient)

	cfg := redis.Config{
		URL:              settings.RedisURL,
		Reset:            false,
		TLSConfig:        nil,
		PoolSize:         10 * runtime.GOMAXPROCS(0),
		Addrs:            []string{},
		MasterName:       "",
		ClientName:       "",
		SentinelUsername: "",
		SentinelPassword: "",
	}

	client.Rdb = redis.New(cfg)
	return client
}

func (c *TypeRedisClient) Close() error {
	defer c.mu.Unlock()

	c.mu.Lock()
	if !c.done {
		if err := c.Rdb.Close(); err != nil {
			return err
		}
		c.done = true
	}

	return nil
}

func (c *TypeRedisClient) Ready(ctx context.Context, logger *log.Logger) bool {
	defer c.mu.Unlock()

	c.mu.Lock()
	if c.done {
		return false
	}

	if err := c.Rdb.Conn().Ping(ctx).Err(); err != nil {
		logger.Error("Failed ping REDIS_URL", "Error", err.Error())
		return false
	}

	return true
}
