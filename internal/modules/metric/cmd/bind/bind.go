package bind

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Init() *cobra.Command {
	return &cobra.Command{
		Use:   "bind",
		Short: "Bind metrics to components",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Bind metric command")
			handler := initializeHandler()
			fmt.Println(handler.Bind())
		},
	}
}
