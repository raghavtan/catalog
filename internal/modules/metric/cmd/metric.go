package cmd

import (
	"github.com/motain/of-catalog/internal/modules/metric/cmd/create"
	"github.com/motain/of-catalog/internal/modules/metric/cmd/delete"
	"github.com/motain/of-catalog/internal/modules/metric/cmd/read"
	"github.com/motain/of-catalog/internal/modules/metric/cmd/update"
	"github.com/spf13/cobra"
)

// InitCmd initializes the metric command and its subcommands.
func InitCmd() *cobra.Command { // Renamed from Init
	metricCmd := &cobra.Command{
		Use:   "metric",
		Short: "metric related commands",
	}

	metricCmd.AddCommand(create.Init())
	metricCmd.AddCommand(read.Init())
	metricCmd.AddCommand(update.Init())
	metricCmd.AddCommand(delete.Init())

	return metricCmd
}
