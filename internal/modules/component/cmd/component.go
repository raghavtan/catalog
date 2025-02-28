package cmd

import (
	"github.com/motain/fact-collector/internal/modules/component/cmd/apply"
	"github.com/spf13/cobra"
)

func Init() *cobra.Command {
	componentCmd := &cobra.Command{
		Use:   "component",
		Short: "component related commands",
	}

	componentCmd.AddCommand(apply.Init())

	return componentCmd
}
