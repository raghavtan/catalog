package main

import (
	"fmt"

	component "github.com/motain/of-catalog/internal/modules/component/cmd"
	metric "github.com/motain/of-catalog/internal/modules/metric/cmd"
	scorecard "github.com/motain/of-catalog/internal/modules/scorecard/cmd"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ofc",
	Short: "⚽ onefootball catalog CLI",
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
