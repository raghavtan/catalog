package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Init() *cobra.Command {
	componentCmd := &cobra.Command{
		Use:   "component",
		Short: "component related commands",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("component command")
		},
	}

	componentCmd.AddCommand(InitApply())

	return componentCmd
}
