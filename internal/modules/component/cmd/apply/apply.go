package apply

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Init() *cobra.Command {
	return &cobra.Command{
		Use:   "apply",
		Short: "Apply changes to components",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("apply component command")
			handler := initializeHandler()
			fmt.Println(handler.Apply())
		},
	}
}
