package delete

import (
	"context"
	"encoding/json"
	"fmt"

	rootcmd "github.com/motain/of-catalog/cmd" // Import for centralized services
	"github.com/motain/of-catalog/internal/modules/scorecard/repository/dtos"
	// "github.com/motain/of-catalog/internal/services/compassservice" // No longer init locally
	// "github.com/motain/of-catalog/internal/services/configservice"  // No longer init locally
	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete <SCORECARD_ID>",
	Short: "Delete a scorecard by its ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		scorecardID := args[0]

		// Services are now initialized in rootcmd

		inputDTO := dtos.DeleteScorecardInput{
			ScorecardID: scorecardID,
			// CompassCloudID is not part of DeleteScorecardInput
		}
		var outputDTO dtos.DeleteScorecardOutput
		err := rootcmd.CompassSvc.RunWithDTOs(context.Background(), &inputDTO, &outputDTO)
		if err != nil {
			return fmt.Errorf("error executing delete scorecard: %w", err)
		}
		if !outputDTO.IsSuccessful() {
			return fmt.Errorf("failed to delete scorecard (ID: %s): %v", scorecardID, outputDTO.GetErrors())
		}
		jsonResponse, _ := json.MarshalIndent(outputDTO, "", "  ")
		fmt.Println(string(jsonResponse))
		return nil
	},
}

func Init() *cobra.Command { return DeleteCmd }
