package cmd

import (
	"fmt"

	"github.com/motain/fact-collector/internal/modules/scorecard/cmd/apply"
	"github.com/spf13/cobra"
)

func Init() *cobra.Command {
	componentCmd := &cobra.Command{
		Use:   "scorecard",
		Short: "scorecards related commands",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("scorecards command")
		},
	}

	componentCmd.AddCommand(apply.Init())

	return componentCmd
}
