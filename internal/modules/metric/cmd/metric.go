package cmd

import (
	"github.com/motain/fact-collector/internal/modules/metric/cmd/apply"
	"github.com/motain/fact-collector/internal/modules/metric/cmd/bind"
	"github.com/motain/fact-collector/internal/modules/metric/cmd/compute"
	"github.com/spf13/cobra"
)

func Init() *cobra.Command {
	metricCmd := &cobra.Command{
		Use:   "metric",
		Short: "metric related commands",
	}

	metricCmd.AddCommand(apply.Init())
	metricCmd.AddCommand(bind.Init())
	metricCmd.AddCommand(compute.Init())

	return metricCmd
}
