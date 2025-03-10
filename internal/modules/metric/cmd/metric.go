package cmd

import (
	"github.com/motain/of-catalog/internal/modules/metric/cmd/apply"
	"github.com/spf13/cobra"
)

func Init() *cobra.Command {
	metricCmd := &cobra.Command{
		Use:   "metric",
		Short: "metric related commands",
	}

	metricCmd.AddCommand(apply.Init())

	return metricCmd
}
