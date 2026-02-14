package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "Configs-GO web API.",
	Long:  `Please change message on iternal/cmd/root.go:13 to project description`,
}

func Execute() {
	rootCmd.AddCommand(serverStartCmd)

	err := rootCmd.ExecuteContext(context.Background())
	if err != nil {
		panic(err)
	}
}
