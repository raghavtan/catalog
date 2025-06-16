package create

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	// "os" // No longer needed for local service init

	rootcmd "github.com/motain/of-catalog/cmd" // Import for centralized services
	"github.com/motain/of-catalog/internal/modules/component/repository/dtos"
	"github.com/motain/of-catalog/internal/modules/component/resources"
	// "github.com/motain/of-catalog/internal/services/compassservice" // No longer init locally
	// "github.com/motain/of-catalog/internal/services/configservice"  // No longer init locally
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// CreateCmd represents the create command
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new component",
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath, _ := cmd.Flags().GetString("file")
		if filePath == "" {
			return fmt.Errorf("file path (-f) is required")
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

		// Use centralized services
		// Note: CompassCloudID is obtained from the CompassSvc as it includes ConfigService logic
		inputDTO := dtos.CreateComponentInput{
			CompassCloudID: rootcmd.CompassSvc.GetCompassCloudId(),
			Component:      componentResource,
		}

		var outputDTO dtos.CreateComponentOutput

		// Use the centralized CompassSvc
		err = rootcmd.CompassSvc.RunWithDTOs(context.Background(), &inputDTO, &outputDTO)
		if err != nil {
			return fmt.Errorf("error executing create component operation: %w", err)
		}

		if !outputDTO.IsSuccessful() {
			return fmt.Errorf("failed to create component: %v", outputDTO.GetErrors())
		}

		jsonResponse, err := json.MarshalIndent(outputDTO, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshaling response to JSON: %w", err)
		}

		fmt.Println(string(jsonResponse))
		return nil
	},
}

// Init initializes the create command.
func Init() *cobra.Command {
	return CreateCmd
}
