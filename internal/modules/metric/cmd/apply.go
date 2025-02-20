package cmd

import (
	"fmt"

	"github.com/motain/fact-collector/internal/modules/metric/app"
	"github.com/spf13/cobra"
)

func InitApply() *cobra.Command {
	return &cobra.Command{
		Use:   "apply",
		Short: "Apply changes to metrics",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("apply metric command")
			handler := app.InitializeHandler()
			fmt.Println(handler.Handle())
		},
	}
}
