package cmd

import (
	"fmt"

	"github.com/motain/fact-collector/internal/app"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fact-collector",
	Short: "🦋 fact collector CLI",
	Run: func(cmd *cobra.Command, args []string) {
		handler := app.InitializeHandler()
		fmt.Println(handler.Handle())
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
