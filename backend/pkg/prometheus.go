package pkg

import (
	"context"
	"net/http"
	"time"

	"github.com/Fact0RR/morza/internal/configs"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// see more - https://prometheus.io/docs/guides/go-application/
func InitPrometheus(ctx context.Context, settings *configs.Settings, logger *log.Logger) {
	allOk := true
	if IsNumber(ctx, settings.PrometheusPort, logger) {
		port := ":" + settings.PrometheusPort
		http.Handle(settings.PrometheusURI, promhttp.Handler())

		srv := &http.Server{
			Addr:         port,
			Handler:      nil,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		}

		err := srv.ListenAndServe()
		if err != nil {
			allOk = false
		}
	} else {
		allOk = false
	}

	if !allOk {
		panic("I could not initialize Prometheus!")
	}

	logger.Debug("Prometheus was configured!")
}

func InitFiberPrometheus(settings *configs.Settings, app *fiber.App, logger *log.Logger) {
	if settings.EnabledPrometheus {
		prometheus := fiberprometheus.New(settings.AppName)
		prometheus.RegisterAt(app, "/metrics")
		app.Use(prometheus.Middleware)
		logger.Debug("Prometheus was configured!")
	} else {
		logger.Debug("Prometheus disabled!")
	}
}
