package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func InitMetric() *cobra.Command {
	metricCmd := &cobra.Command{
		Use:   "metric",
		Short: "metric related commands",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("metric command")
		},
	}

	metricCmd.AddCommand(InitApply())

	return metricCmd
}
