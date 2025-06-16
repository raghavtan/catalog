package delete

import (
	"context"
	"encoding/json"
	"fmt"

	rootcmd "github.com/motain/of-catalog/cmd" // Import for centralized services
	"github.com/motain/of-catalog/internal/modules/component/repository/dtos"
	"github.com/motain/of-catalog/internal/services/compassservice" // Keep for HasNotFoundError, or move HasNotFoundError
	// "github.com/motain/of-catalog/internal/services/configservice"  // No longer init locally
	"github.com/spf13/cobra"
)

// DeleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:   "delete <COMPONENT_ID>",
	Short: "Delete a component by its ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		componentID := args[0]

		// Services are now initialized in rootcmd

		inputDTO := dtos.DeleteComponentInput{
			ComponentID: componentID,
		}

		var outputDTO dtos.DeleteComponentOutput

		err := rootcmd.CompassSvc.RunWithDTOs(context.Background(), &inputDTO, &outputDTO)
		if err != nil {
			// Retain HasNotFoundError logic if it's specific to compassservice types
			if outputDTO.IsSuccessful() && compassservice.HasNotFoundError(outputDTO.GetErrors()) {
				// Silently succeed if DTO considers "not found" as success for deletion.
			} else {
				return fmt.Errorf("error executing delete component operation: %w", err)
			}
		}

		if !outputDTO.IsSuccessful() {
			return fmt.Errorf("failed to delete component (ID: %s): %v", componentID, outputDTO.GetErrors())
		}

		jsonResponse, err := json.MarshalIndent(outputDTO, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshaling response to JSON: %w", err)
		}

		fmt.Println(string(jsonResponse))
		return nil
	},
}

// Init initializes the delete command.
func Init() *cobra.Command {
	return DeleteCmd
}
