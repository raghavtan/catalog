package compute

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Init() *cobra.Command {
	var componentType, componentName, metricName string

	cmd := &cobra.Command{
		Use:   "compute",
		Short: "Compute metrics for components",
		Run: func(cmd *cobra.Command, args []string) {
			if componentType != "service" && componentType != "cloud-resource" && componentType != "website" && componentType != "application" {
				fmt.Println("Invalid component type. Accepted values are: service, cloud-resource, website, application")
				return
			}
			fmt.Printf("Tracking metric '%s' for component '%s' of type '%s'\n", metricName, componentName, componentType)
			handler := initializeHandler()
			fmt.Println(handler.Compute(componentType, componentName, metricName))
		},
	}

	cmd.Flags().StringVarP(&componentType, "type", "t", "", "Type of the component (service, cloud-resource, website, application)")
	cmd.Flags().StringVarP(&componentName, "component", "c", "", "Name of the component")
	cmd.Flags().StringVarP(&metricName, "metric", "m", "", "Name of the metric")

	return cmd
}
