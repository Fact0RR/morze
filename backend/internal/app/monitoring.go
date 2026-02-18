package app

import (
	"github.com/Fact0RR/morze/internal/configs"
	"github.com/Fact0RR/morze/pkg"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// Initialize Sentry, Opentelemetry, Prometheus.
func initFiberMonitoring(settings *configs.Settings, app *fiber.App, logger *log.Logger) {
	if settings.EnabledPrometheus {
		pkg.InitFiberPrometheus(settings, app, logger)
		setPrometheusServiceVersion(settings)
	} else {
		logger.Warn("Prometheus is disabled!")
	}

	if settings.EnabledSentry {
		pkg.InitFiberSentry(settings, app, logger)
	} else {
		logger.Warn("Sentry is disabled!")
	}
}

func setPrometheusServiceVersion(settings *configs.Settings) {
	PromServiceVersion := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "morze",
			Name:      "service_version",
			Help:      "The information of a service version",
			ConstLabels: prometheus.Labels{
				"version":     settings.Version,
				"commit_hash": settings.CommitHash,
				"environment": settings.Environment,
			},
		})

	prometheus.MustRegister(PromServiceVersion)
}
