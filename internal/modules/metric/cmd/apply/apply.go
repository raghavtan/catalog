package apply

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Init() *cobra.Command {
	return &cobra.Command{
		Use:   "apply",
		Short: "Apply changes to metrics",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("apply metric command")
			handler := initializeHandler()
			fmt.Println(handler.Apply())
		},
	}
}
