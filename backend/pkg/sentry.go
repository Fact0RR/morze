package pkg

import (
	"fmt"
	"time"

	"github.com/Fact0RR/morze/internal/configs"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/contrib/fibersentry"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	log "github.com/sirupsen/logrus"
)

// Initialize sentry client with custom settings
// see more - https://docs.sentry.io/platforms/go/
func InitSentry(settings *configs.Settings, logger *log.Logger) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              settings.SentryDsn,
		Debug:            settings.Debug,
		SampleRate:       settings.SentrySamplerRate,
		TracesSampleRate: settings.SentryTracesSamplerRate,
		Release:          "morze@" + settings.Version,
	})
	if err != nil {
		logger.Error("sentry.Init", err)
	} else {
		logger.Info("Sentry was configured!")
	}

	defer sentry.Flush(time.Duration(settings.SentryFlushSecondsCount) * time.Second)
}

func InitFiberSentry(settings *configs.Settings, app *fiber.App, logger *log.Logger) {
	if settings.EnabledSentry {
		_ = sentry.Init(sentry.ClientOptions{
			Dsn: settings.SentryDsn,
			BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
				if hint.Context != nil {
					if c, ok := hint.Context.Value(sentry.RequestContextKey).(*fiber.Ctx); ok {
						// You have access to the original Context if it panicked.
						fmt.Println(utils.CopyString(c.Hostname()))
					}
				}
				fmt.Println(event)
				return event
			},
			Debug:            true,
			AttachStacktrace: true,
			Release:          "morze@" + settings.Version,
			TracesSampleRate: settings.SentryTracesSamplerRate,
			SampleRate:       settings.SentrySamplerRate,
		})
		app.Use(fibersentry.New(fibersentry.Config{
			Repanic:         true,
			WaitForDelivery: true,
		}))
		logger.Debug("Sentry configured.")
	} else {
		logger.Warn("Sentry disabled.")
	}
}
