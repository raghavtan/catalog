package compute

import (
	"fmt"

	"github.com/motain/fact-collector/internal/utils/yaml"
	"github.com/spf13/cobra"
)

func Init() *cobra.Command {
	var componentType, componentName, metricName string

	cmd := &cobra.Command{
		Use:   "compute",
		Short: "Compute metrics for components",
		Run: func(cmd *cobra.Command, args []string) {
			if componentType == "" {
				fmt.Println("Error: componentType is required")
				cmd.Help()
				return
			}
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

			fmt.Printf("Tracking metric '%s' for component '%s' of type '%s'\n", metricName, componentName, componentType)
			handler := initializeHandler()
			handler.Compute(componentType, componentName, metricName, yaml.StateLocation)
		},
	}

	cmd.Flags().StringVarP(&componentType, "type", "t", "", "Type of the component (service, cloud-resource, website, application)")
	cmd.Flags().StringVarP(&componentName, "component", "c", "", "Name of the component")
	cmd.Flags().StringVarP(&metricName, "metric", "m", "", "Name of the metric")

	return cmd
}
