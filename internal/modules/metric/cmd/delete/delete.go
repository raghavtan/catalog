package delete

import (
	"context"
	"encoding/json"
	"fmt"

	rootcmd "github.com/motain/of-catalog/cmd" // Import for centralized services
	"github.com/motain/of-catalog/internal/modules/metric/repository/dtos"
	// "github.com/motain/of-catalog/internal/services/compassservice" // No longer init locally
	// "github.com/motain/of-catalog/internal/services/configservice"  // No longer init locally
	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete <METRIC_ID>",
	Short: "Delete a metric by its ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		metricID := args[0]

		// Services are now initialized in rootcmd

		inputDTO := dtos.DeleteMetricInput{
			// CompassCloudID is not typically part of DeleteMetricInput if ID is globally unique
			// CompassCloudID: rootcmd.CompassSvc.GetCompassCloudId(),
			MetricID:       metricID,
		}
		var outputDTO dtos.DeleteMetricOutput
		err := rootcmd.CompassSvc.RunWithDTOs(context.Background(), &inputDTO, &outputDTO)
		if err != nil {
			return fmt.Errorf("error executing delete metric: %w", err)
		}
		if !outputDTO.IsSuccessful() {
			return fmt.Errorf("failed to delete metric (ID: %s): %v", metricID, outputDTO.GetErrors())
		}
		jsonResponse, _ := json.MarshalIndent(outputDTO, "", "  ")
		fmt.Println(string(jsonResponse))
		return nil
	},
}

func Init() *cobra.Command { return DeleteCmd }
