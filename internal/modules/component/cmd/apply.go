package cmd

import (
	"fmt"

	"github.com/motain/fact-collector/internal/modules/component/app"
	"github.com/spf13/cobra"
)

func InitApply() *cobra.Command {
	return &cobra.Command{
		Use:   "apply",
		Short: "Apply changes to components",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("apply component command")
			handler := app.InitializeHandler()
			fmt.Println(handler.Apply())
		},
	}
}
