package cmd

import (
	"github.com/motain/of-catalog/internal/modules/scorecard/cmd/create"
	"github.com/motain/of-catalog/internal/modules/scorecard/cmd/delete"
	// Placeholder for read when DTO is identified
	"github.com/motain/of-catalog/internal/modules/scorecard/cmd/update"
	"github.com/spf13/cobra"
)

// InitCmd initializes the scorecard command and its subcommands.
func InitCmd() *cobra.Command { // Renamed from Init
	scorecardCmd := &cobra.Command{
		Use:   "scorecard",
		Short: "scorecard related commands",
	}

	scorecardCmd.AddCommand(create.Init())
	scorecardCmd.AddCommand(update.Init())
	scorecardCmd.AddCommand(delete.Init())
	// scorecardCmd.AddCommand(read.Init()) // Add when read is implemented

	return scorecardCmd
}
