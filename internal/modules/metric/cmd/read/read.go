package read

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	rootcmd "github.com/motain/of-catalog/cmd" // Import for centralized services
	"github.com/motain/of-catalog/internal/modules/metric/repository/dtos"
	"github.com/motain/of-catalog/internal/modules/metric/resources" // For Metric resource structure
	// "github.com/motain/of-catalog/internal/services/compassservice" // No longer init locally
	// "github.com/motain/of-catalog/internal/services/configservice"  // No longer init locally
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var ReadCmd = &cobra.Command{
	Use:   "read",
	Short: "Read/Search metrics. Optionally use -f to specify a metric name for searching.",
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath, _ := cmd.Flags().GetString("file")
        var metricResource resources.Metric // Defaults to empty, searching all

        if filePath != "" {
		yamlFile, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("error reading YAML file for search: %w", err)
		}
		err = yaml.Unmarshal(yamlFile, &metricResource)
		if err != nil {
			return fmt.Errorf("error unmarshaling YAML for search: %w", err)
		}
        }

		// Services are now initialized in rootcmd

		inputDTO := dtos.SearchMetricsInput{
			CompassCloudID: rootcmd.CompassSvc.GetCompassCloudId(),
			Metric:         metricResource,
		}
		var outputDTO dtos.SearchMetricsOutput
		err := rootcmd.CompassSvc.RunWithDTOs(context.Background(), &inputDTO, &outputDTO)
		if err != nil {
			return fmt.Errorf("error executing search metric: %w", err)
		}
		if !outputDTO.IsSuccessful() {
			return fmt.Errorf("failed to search metric: %v", outputDTO.GetErrors())
		}
		jsonResponse, _ := json.MarshalIndent(outputDTO, "", "  ")
		fmt.Println(string(jsonResponse))
		return nil
	},
}

func Init() *cobra.Command { return ReadCmd }
