package app

import (
	"context"
	"time"

	"github.com/Fact0RR/morza/internal/cache/redis"
	"github.com/Fact0RR/morza/internal/configs"
	"github.com/Fact0RR/morza/internal/controller"
	"github.com/Fact0RR/morza/internal/middleware"
	database "github.com/Fact0RR/morza/internal/repository/postgres"
	"github.com/Fact0RR/morza/internal/service"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type App struct {
	server   *fiber.App
	db       *database.ClientDB
	cache    *redis.TypeRedisClient
	settings *configs.Settings
	logger   *log.Logger
}

func InitApp(ctx context.Context) *App {
	settings := configs.InitSettings()
	var logger *log.Logger
	logger = log.New()
	if settings.Debug {
		logger.SetLevel(log.DebugLevel)
	} else {
		logger.SetLevel(log.InfoLevel)
	}

	logger.Debug("Time settings ",
		"Debug", settings.Debug,
		"UTC_time", time.Now().UTC(),
		"sentry", settings.EnabledSentry,
	)

	postgresClient := database.InitDB(ctx, &settings, logger)
	var redisClient *redis.TypeRedisClient
	if settings.EnabledRedis {
		redisClient = redis.NewTypeRedisClient(settings.RedisSettings, logger)
	} else {
		logger.Warn("Кеш выключен!")
	}
	repo := database.NewMorzaRepo(postgresClient.DB, logger)
	cache := redis.NewMorzaCache(redisClient, settings.TTL, settings.EnabledRedis, logger)
	// Остужаю кеш при зауске, чтобы сохранить консистентность данных.
	if err := cache.CoolAll(ctx); err != nil {
		logger.Error("Не получилось охладить кеш")
	}
	service := service.NewMorzaService(repo, logger)
	baseController := controller.NewBaseController(&settings)
	morzasController := controller.NewChangeMorzaController(service, logger)

	server := fiber.New(fiber.Config{
		BodyLimit:             settings.BodyLimit,
		ReadTimeout:           settings.ReadTimeout,
		WriteTimeout:          settings.WriteTimeout,
		IdleTimeout:           settings.IdleTimeout,
		DisableStartupMessage: true,
	})

	initFiberMonitoring(&settings, server, logger)
	api := server.Group("/api")
	baseController.RegisterRotes(api)
	morzasController.RegisterRoutes(api, middleware.AuthMiddleware(settings.JwtSigningKeyBytes, logger))

	return &App{
		server:   server,
		db:       postgresClient,
		cache:    redisClient,
		settings: &settings,
		logger:   logger,
	}
}

// RunApp main application.
func (a *App) RunApp(ctx context.Context) error {
	a.logger.Info("Запуск сервера ", "port:", a.settings.AppPort)

	return a.server.Listen(":" + a.settings.AppPort)
}

func (a *App) GracefulShutdownApp(ctx context.Context) error {
	a.db.CloseDB(ctx)
	a.logger.Debug("Database client closed")
	if a.cache != nil {
		err := a.cache.Rdb.Close()
		if err != nil {
			return err
		}
		a.logger.Debug("Redis client closed")
	}
	return a.server.Shutdown()
}
