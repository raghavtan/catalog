package compute

import (
	"fmt"

	"github.com/motain/fact-collector/internal/utils/yaml"
	"github.com/spf13/cobra"
)

func Init() *cobra.Command {
	var componentName, metricName string

	cmd := &cobra.Command{
		Use:   "compute",
		Short: "Compute metrics for components",
		Run: func(cmd *cobra.Command, args []string) {
			if componentName == "" {
				fmt.Println("Error: componentName is required")
				cmd.Help()
				return
			}
			if metricName == "" {
				fmt.Println("Error: metricName is required")
				cmd.Help()
				return
			}

			fmt.Printf("Tracking metric '%s' for component '%s'\n", metricName, componentName)
			handler := initializeHandler()
			handler.Compute(componentName, metricName, yaml.StateLocation)
		},
	}

	cmd.Flags().StringVarP(&componentName, "component", "c", "", "Name of the component")
	cmd.Flags().StringVarP(&metricName, "metric", "m", "", "Name of the metric")

	return cmd
}
