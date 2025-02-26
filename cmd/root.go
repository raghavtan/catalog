package main

import (
	"fmt"

	component "github.com/motain/fact-collector/internal/modules/component/cmd"
	metric "github.com/motain/fact-collector/internal/modules/metric/cmd"
	scorecard "github.com/motain/fact-collector/internal/modules/scorecard/cmd"
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
	rootCmd.AddCommand(component.Init())
	rootCmd.AddCommand(metric.Init())
	rootCmd.AddCommand(scorecard.Init())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func main() {
	Execute()
}
