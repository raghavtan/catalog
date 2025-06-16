package update

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	rootcmd "github.com/motain/of-catalog/cmd" // Import for centralized services
	"github.com/motain/of-catalog/internal/modules/component/repository/dtos"
	"github.com/motain/of-catalog/internal/modules/component/resources"
	// "github.com/motain/of-catalog/internal/services/compassservice" // No longer init locally
	// "github.com/motain/of-catalog/internal/services/configservice"  // No longer init locally
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// UpdateCmd represents the update command
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing component",
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath, _ := cmd.Flags().GetString("file")
		if filePath == "" {
			return fmt.Errorf("file path (-f) is required for update")
		}

		yamlFile, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("error reading YAML file: %w", err)
		}

		var componentResource resources.Component
		err = yaml.Unmarshal(yamlFile, &componentResource)
		if err != nil {
			return fmt.Errorf("error unmarshaling YAML to component resource: %w", err)
		}

		if componentResource.ID == "" {
			return fmt.Errorf("component ID is missing in the YAML file, which is required for update")
		}

		// Services are now initialized in rootcmd

		inputDTO := dtos.UpdateComponentInput{
			Component: componentResource,
			// CompassCloudID is not part of UpdateComponentInput for components.
			// It's directly on the component itself if needed by the mutation, or handled by service.
		}

		var outputDTO dtos.UpdateComponentOutput

		err = rootcmd.CompassSvc.RunWithDTOs(context.Background(), &inputDTO, &outputDTO)
		if err != nil {
			return fmt.Errorf("error executing update component operation: %w", err)
		}

		if !outputDTO.IsSuccessful() {
			return fmt.Errorf("failed to update component (ID: %s): %v", componentResource.ID, outputDTO.GetErrors())
		}

		jsonResponse, err := json.MarshalIndent(outputDTO, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshaling response to JSON: %w", err)
		}

		fmt.Println(string(jsonResponse))
		return nil
	},
}

// Init initializes the update command.
func Init() *cobra.Command {
	return UpdateCmd
}
