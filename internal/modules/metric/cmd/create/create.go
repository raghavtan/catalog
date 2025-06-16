package create

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	rootcmd "github.com/motain/of-catalog/cmd" // Import for centralized services
	"github.com/motain/of-catalog/internal/modules/metric/repository/dtos"
	"github.com/motain/of-catalog/internal/modules/metric/resources"
	// "github.com/motain/of-catalog/internal/services/compassservice" // No longer init locally
	// "github.com/motain/of-catalog/internal/services/configservice"  // No longer init locally
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new metric",
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath, _ := cmd.Flags().GetString("file")
		if filePath == "" {
			return fmt.Errorf("file path (-f) is required")
		}
		yamlFile, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("error reading YAML file: %w", err)
		}
		var metricResource resources.Metric
		err = yaml.Unmarshal(yamlFile, &metricResource)
		if err != nil {
			return fmt.Errorf("error unmarshaling YAML: %w", err)
		}

		// Services are now initialized in rootcmd

		inputDTO := dtos.CreateMetricInput{
			CompassCloudID: rootcmd.CompassSvc.GetCompassCloudId(),
			Metric:         metricResource,
		}
		var outputDTO dtos.CreateMetricOutput
		err = rootcmd.CompassSvc.RunWithDTOs(context.Background(), &inputDTO, &outputDTO)
		if err != nil {
			return fmt.Errorf("error executing create metric: %w", err)
		}
		if !outputDTO.IsSuccessful() {
			return fmt.Errorf("failed to create metric: %v", outputDTO.GetErrors())
		}
		jsonResponse, _ := json.MarshalIndent(outputDTO, "", "  ")
		fmt.Println(string(jsonResponse))
		return nil
	},
}

func Init() *cobra.Command { return CreateCmd }
