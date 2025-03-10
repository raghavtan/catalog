package cmd

import (
	"github.com/motain/of-catalog/internal/modules/scorecard/cmd/apply"
	"github.com/spf13/cobra"
)

func Init() *cobra.Command {
	componentCmd := &cobra.Command{
		Use:   "scorecard",
		Short: "scorecards related commands",
	}

	componentCmd.AddCommand(apply.Init())

	return componentCmd
}
