package bind

import (
	"github.com/motain/fact-collector/internal/utils/yaml"
	"github.com/spf13/cobra"
)

func Init() *cobra.Command {
	return &cobra.Command{
		Use:   "bind",
		Short: "Bind metrics to components",
		Run: func(cmd *cobra.Command, args []string) {
			handler := initializeHandler()
			handler.Bind(yaml.StateLocation)
		},
	}
}
