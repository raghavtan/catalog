package cmd

import (
	"fmt"

	"github.com/motain/fact-collector/internal/modules/metric/cmd/apply"
	"github.com/motain/fact-collector/internal/modules/metric/cmd/bind"
	"github.com/spf13/cobra"
)

func Init() *cobra.Command {
	metricCmd := &cobra.Command{
		Use:   "metric",
		Short: "metric related commands",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("metric command")
		},
	}

	metricCmd.AddCommand(apply.Init())
	metricCmd.AddCommand(bind.Init())

	return metricCmd
}
