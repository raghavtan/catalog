package main

import (
	"fmt"

	metric "github.com/motain/fact-collector/internal/modules/metric/cmd"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fact-collector",
	Short: "🦋 fact collector CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fact-collector command")
	},
}

func Execute() {
	rootCmd.AddCommand(metric.InitMetric())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func main() {
	Execute()
}
