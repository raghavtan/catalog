package compute

import (
	"fmt"

	"github.com/motain/of-catalog/internal/utils/commandcontext"
	"github.com/motain/of-catalog/internal/utils/yaml"
	"github.com/spf13/cobra"
)

func Init() *cobra.Command {
	var componentName, metricName string
	var all bool

	cmd := &cobra.Command{
		Use:   "compute",
		Short: "Compute metrics for components",
		Run: func(cmd *cobra.Command, args []string) {
			if componentName == "" {
				fmt.Println("Error: componentName is required")
				cmd.Help()
				return
			}
			if !all && metricName == "" {
				fmt.Println("Error: metricName is required")
				cmd.Help()
				return
			}

			handler := initializeHandler()
			ctx := commandcontext.Init()
			handler.Compute(ctx, componentName, all, metricName, yaml.StateLocation)
		},
	}

	cmd.Flags().StringVarP(&componentName, "component", "c", "", "Name of the component")
	cmd.Flags().StringVarP(&metricName, "metric", "m", "", "Name of the metric")
	cmd.Flags().BoolVarP(&all, "all", "a", false, "Compute all metrics for the component")

	return cmd
}
