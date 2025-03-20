package bind

import (
	"github.com/motain/of-catalog/internal/utils/commandcontext"
	"github.com/motain/of-catalog/internal/utils/yaml"
	"github.com/spf13/cobra"
)

func Init() *cobra.Command {
	return &cobra.Command{
		Use:   "bind",
		Short: "Bind components to metrics",
		Run: func(cmd *cobra.Command, args []string) {
			handler := initializeHandler()
			ctx := commandcontext.Init()
			handler.Bind(ctx, yaml.StateLocation)
		},
	}
}
