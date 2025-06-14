package apply

import (
	"fmt"

	"github.com/motain/of-catalog/internal/utils/commandcontext"
	"github.com/motain/of-catalog/internal/utils/yaml"
	"github.com/spf13/cobra"
)

func Init() *cobra.Command {
	var configRootLocation, componentName string
	var recursive bool

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply changes to scorecards",
		Run: func(cmd *cobra.Command, args []string) {
			if configRootLocation == "" {
				fmt.Println("Error: configRootLocation required")
				cmd.Help()
				return
			}

			handler := initializeHandler()
			ctx := commandcontext.Init()
			handler.Apply(ctx, configRootLocation, yaml.StateLocation, recursive, componentName)
		},
	}

	cmd.Flags().StringVarP(&configRootLocation, "configRootLocation", "l", "", "Root location of the config")
	cmd.Flags().StringVarP(&componentName, "component", "c", "", "Name of the component")
	cmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Apply changes recursively")

	return cmd
}
