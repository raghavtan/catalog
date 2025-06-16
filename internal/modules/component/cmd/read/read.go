package read

import (
	"context"
	"encoding/json"
	"fmt"

	rootcmd "github.com/motain/of-catalog/cmd" // Import for centralized services
	"github.com/motain/of-catalog/internal/modules/component/repository/dtos"
	// "github.com/motain/of-catalog/internal/services/compassservice" // No longer init locally
	// "github.com/motain/of-catalog/internal/services/configservice"  // No longer init locally
	"github.com/spf13/cobra"
)

// ReadCmd represents the read command
var ReadCmd = &cobra.Command{
	Use:   "read <COMPONENT_SLUG>",
	Short: "Read a component by its slug",
	Args:  cobra.ExactArgs(1), // Expects one argument: COMPONENT_SLUG
	RunE: func(cmd *cobra.Command, args []string) error {
		componentSlug := args[0]

		// Services are now initialized in rootcmd

		inputDTO := dtos.ComponentByReferenceInput{
			CompassCloudID: rootcmd.CompassSvc.GetCompassCloudId(),
			Slug:           componentSlug,
		}

		var outputDTO dtos.ComponentByReferenceOutput

		err := rootcmd.CompassSvc.RunWithDTOs(context.Background(), &inputDTO, &outputDTO)
		if err != nil {
			return fmt.Errorf("error executing read component operation: %w", err)
		}

		if !outputDTO.IsSuccessful() {
			return fmt.Errorf("failed to read component (ID: %s). The component might not exist or there was an issue fetching it", componentSlug)
		}

		jsonResponse, err := json.MarshalIndent(outputDTO, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshaling response to JSON: %w", err)
		}

		fmt.Println(string(jsonResponse))
		return nil
	},
}

// Init initializes the read command.
func Init() *cobra.Command {
	return ReadCmd
}
