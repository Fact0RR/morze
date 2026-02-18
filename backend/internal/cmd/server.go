package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Fact0RR/morze/internal/app"
	"github.com/spf13/cobra"
)

var serverStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Run main application Configs-GO.",
	Long:  ``,
	Run: func(_ *cobra.Command, _ []string) {
		application := app.InitApp(rootCmd.Context())

		go func() {
			if err := application.RunApp(rootCmd.Context()); err != nil {
				log.Fatalf("Ошибка при запуске сервера: %v.", err)
			}
		}()

		// Ожидание сигналов для graceful shutdown.
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		// Graceful shutdown.
		if err := application.GracefulShutdownApp(rootCmd.Context()); err != nil {
			log.Fatalf("Ошибка при остановке сервера: %v.", err)
		}
	},
}
