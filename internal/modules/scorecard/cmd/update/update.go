package update

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	rootcmd "github.com/motain/of-catalog/cmd" // Import for centralized services
	"github.com/motain/of-catalog/internal/modules/scorecard/repository/dtos"
	"github.com/motain/of-catalog/internal/modules/scorecard/resources"
	// "github.com/motain/of-catalog/internal/services/compassservice" // No longer init locally
	// "github.com/motain/of-catalog/internal/services/configservice"  // No longer init locally
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing scorecard",
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath, _ := cmd.Flags().GetString("file")
		if filePath == "" {
			return fmt.Errorf("file path (-f) is required")
		}
		yamlFile, err := ioutil.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("error reading YAML file: %w", err)
		}

		var scorecardResourceForUpdate resources.Scorecard
		err = yaml.Unmarshal(yamlFile, &scorecardResourceForUpdate)
		if err != nil {
			return fmt.Errorf("error unmarshaling YAML for scorecard update: %w", err)
		}

        if scorecardResourceForUpdate.ID == "" {
            return fmt.Errorf("scorecard ID is missing in YAML, required for update")
        }

		// Services are now initialized in rootcmd

		inputDTO := dtos.UpdateScorecardInput{
			Scorecard:      scorecardResourceForUpdate,
			CreateCriteria: scorecardResourceForUpdate.CreateCriteria,
			UpdateCriteria: scorecardResourceForUpdate.UpdateCriteria,
			DeleteCriteria: scorecardResourceForUpdate.DeleteCriteriaIDs,
			// CompassCloudID is not part of UpdateScorecardInput, assumed to be handled by service or within Scorecard fields
		}
		var outputDTO dtos.UpdateScorecardOutput
		err = rootcmd.CompassSvc.RunWithDTOs(context.Background(), &inputDTO, &outputDTO)
		if err != nil {
			return fmt.Errorf("error executing update scorecard: %w", err)
		}
		if !outputDTO.IsSuccessful() {
			return fmt.Errorf("failed to update scorecard (ID: %s): %v", scorecardResourceForUpdate.ID, outputDTO.GetErrors())
		}
		jsonResponse, _ := json.MarshalIndent(outputDTO, "", "  ")
		fmt.Println(string(jsonResponse))
		return nil
	},
}

func Init() *cobra.Command { return UpdateCmd }
