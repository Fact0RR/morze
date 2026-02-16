package configs

import (
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// All settings for this project.
type Settings struct {
	AppHost             string `default:""                          envconfig:"APP_PORT"`
	AppPort             string `default:"8000"                      envconfig:"APP_PORT"`
	AppName             string `default:"morza"                     envconfig:"APP_NAME"`
	Environment         string `default:"dev"                       envconfig:"APP_ENVIRONMENT"`
	Debug               bool   `default:"false"                     envconfig:"DEBUG"`
	Version             string `default:"0.1.0"                     envconfig:"VERSION"`
	CommitHash          string `default:""                          envconfig:"COMMIT_HASH"`
	JwtSigningKeyString string `default:""                          envconfig:"JWT_SIGNING_KEY"`
	JwtSigningKeyBytes  []byte

	DBSettings
	RedisSettings
	MonitoringSettings
	HTTPConfig
}

type DBSettings struct {
	DatabaseURL               string `default:"postgres://morza:morza@postgres:5432/morza?sslmode=disable"                    envconfig:"DATABASE_URL"`
	MaxIdlePgConnections      int32  `default:"100"                                                                     envconfig:"MAX_IDLE_PG_CONNECTIONS"`
	IdleLifetimePgConnections int    `default:"1200"                                                                    envconfig:"IDLE_LIFETIME_PG_CONNECTIONS"` // sec
	QueryExecMode             int32  `default:"1"                                                                       envconfig:"QUERY_EXEC_MODE"`
	DBConnectionRetries       int    `default:"5"                                                                       envconfig:"DB_CONNECTION_RETRIES"`
}

type RedisSettings struct {
	RedisURL     string `default:"redis://redis:6379/0"  envconfig:"REDIS_URL"`
	EnabledRedis bool   `default:"false"                 envconfig:"ENABLED_REDIS"`
	TTL          int    `default:"300"                   envconfig:"LIFETIME_CONFIG"` // секунды
}

type MonitoringSettings struct {
	EnabledOpentelemetry    bool    `default:"false"                                          envconfig:"ENABLED_OPENTELEMETRY"`
	EnabledSentry           bool    `default:"false"                                          envconfig:"ENABLED_SENTRY"`
	EnabledPrometheus       bool    `default:"false"                                          envconfig:"ENABLED_PROMETHEUS"`
	OpentelemetryEndpoint   string  `default:"localhost:4317"                                 envconfig:"OPENTELEMETRY_ENDPOINT"`
	PrometheusPort          string  `default:"9090"                                           envconfig:"PROMETHEUS_PORT"`
	PrometheusURI           string  `default:"/metrics"                                       envconfig:"PROMETHEUS_URI"`
	SentryDsn               string  `default:"https://examplePublicKey@o0.ingest.sentry.io/0" envconfig:"SENTRY_DSN"`
	SentrySamplerRate       float64 `default:"0.05"                                           envconfig:"SENTRY_SAMPLER_RATE"`
	SentryTracesSamplerRate float64 `default:"0.05"                                           envconfig:"SENTRY_TRACES_SAMPLER_RATE"`
	SentryFlushSecondsCount int     `default:"2"                                              envconfig:"SENTRY_FLUSH_SECONDS_COUNT"`
}

type HTTPConfig struct {
	BodyLimit    int           `envconfig:"BODY_LIMIT"    default:"4194304"`
	ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT"  default:"5s"`
	WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT" default:"10s"`
	IdleTimeout  time.Duration `envconfig:"IDLE_TIMEOUT"  default:"60s"`
}

// InitSettings Read .env file and write values from environment variables to settings.
func InitSettings() Settings {
	var settings Settings
	logger := log.New()
	if err := godotenv.Load(); err != nil {
		logger.Warnf("Не удалось загрузить .env файл. Ошибка %q", err.Error())
	}

	if err := envconfig.Process("MORZA", &settings); err != nil {
		logger.Warnf("Не удалось загрузить переменные окружения. Ошибка %q", err.Error())
	}

	settings.JwtSigningKeyBytes = []byte(settings.JwtSigningKeyString)
	return settings
}
