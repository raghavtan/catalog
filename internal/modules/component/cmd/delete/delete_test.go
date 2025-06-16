package delete

import (
	"bytes"
	"context"
	"encoding/json"
	// "io/ioutil" // Not used in the final version of the test code
	// "os" // Not used
	"strings"
	"testing"
	"fmt" // Added for service error test case

	// Import the root cmd package to set up services for testing
	rootcmd "github.com/motain/of-catalog/cmd"
	"github.com/motain/of-catalog/internal/modules/component/repository/dtos"
	"github.com/motain/of-catalog/internal/services/compassservice"
	dtoiface "github.com/motain/of-catalog/internal/services/compassservice/dtos" // Alias for compassservice.dtos
	"github.com/motain/of-catalog/internal/services/configservice"
	"github.com/spf13/cobra"
)

// MockCompassService is a mock implementation of CompassServiceInterface
type MockCompassService struct {
	RunFunc                   func(ctx context.Context, query string, variables map[string]interface{}, response interface{}) error
	RunWithDTOsFunc           func(ctx context.Context, input compassservice.InputDTOInterface, output compassservice.OutputDTOInterface) error
	SendMetricFunc            func(ctx context.Context, body map[string]string) (string, error)
	SendAPISpecificationsFunc func(ctx context.Context, input dtoiface.APISpecificationsInput) (string, error) // Corrected DTO import path
	GetCompassCloudIdFunc     func() string
}

func (m *MockCompassService) Run(ctx context.Context, query string, variables map[string]interface{}, response interface{}) error {
	if m.RunFunc != nil {
		return m.RunFunc(ctx, query, variables, response)
	}
	return nil
}

func (m *MockCompassService) RunWithDTOs(ctx context.Context, input compassservice.InputDTOInterface, output compassservice.OutputDTOInterface) error {
	if m.RunWithDTOsFunc != nil {
		return m.RunWithDTOsFunc(ctx, input, output)
	}
	// Default mock behavior for delete: populate output to simulate success
    if out, ok := output.(*dtos.DeleteComponentOutput); ok {
        // Simulate the nested structure if IsSuccessful() depends on it.
        // Based on DTOs, Success is often at Compass.DeleteComponent.Success
        if out.Compass.DeleteComponent == nil { // Ensure nested struct is initialized
            out.Compass.DeleteComponent = &struct{Success bool `json:"success"`}{}
        }
        out.Compass.DeleteComponent.Success = true
    }
	return nil
}
func (m *MockCompassService) SendMetric(ctx context.Context, body map[string]string) (string, error) { return "", nil }
func (m *MockCompassService) SendAPISpecifications(ctx context.Context, input dtoiface.APISpecificationsInput) (string, error) { return "", nil} // Corrected DTO import path
func (m *MockCompassService) GetCompassCloudId() string {
	if m.GetCompassCloudIdFunc != nil {
		return m.GetCompassCloudIdFunc()
	}
	return "test-cloud-id"
}


// Helper function to execute cobra commands and capture output
func executeCommand(cmd *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)
	// Must call PersistentPreRunE manually if the test is on a subcommand and PersistentPreRunE is on a parent.
	// However, here we are testing Init() which returns the DeleteCmd itself.
	// If DeleteCmd or its parents have PersistentPreRun functions, they should be called by Execute().
	// rootcmd.RootCmd itself has the PersistentPreRunE for service init.
	// To make this test self-contained for the subcommand, service init is handled in TestDeleteCmd.
	err := cmd.Execute()
	return buf.String(), err
}

func TestDeleteCmd(t *testing.T) {
	// Setup: Initialize the root command's services with mocks
	originalCompassSvc := rootcmd.CompassSvc
	mockService := &MockCompassService{}

    // Critical: Ensure that the services in rootcmd are set *before* Init() is called if Init() somehow uses them.
    // For this command structure, Init() just returns the command object; RunE is where services are used.
	rootcmd.CompassSvc = mockService

    if rootcmd.CfgService == nil {
        rootcmd.CfgService = configservice.NewConfigService()
    }

	// Restore original services after test
	defer func() { rootcmd.CompassSvc = originalCompassSvc }()

	// Get the command to test (DeleteCmd)
    // We call Init() to get the command instance, similar to how it's added to its parent.
    deleteCmdInstance := Init()

	t.Run("successful delete", func(t *testing.T) {
		mockService.RunWithDTOsFunc = func(ctx context.Context, input compassservice.InputDTOInterface, output compassservice.OutputDTOInterface) error {
			if delInput, ok := input.(*dtos.DeleteComponentInput); ok {
				if delInput.ComponentID != "test-component-id" {
					t.Errorf("Expected component ID 'test-component-id', got '%s'", delInput.ComponentID)
				}
			} else {
				t.Errorf("Unexpected input DTO type")
			}

			if out, ok := output.(*dtos.DeleteComponentOutput); ok {
                if out.Compass.DeleteComponent == nil {
                     out.Compass.DeleteComponent = &struct{Success bool `json:"success"`}{}
                }
				out.Compass.DeleteComponent.Success = true
			} else {
				t.Errorf("Unexpected output DTO type for delete")
			}
			return nil
		}

		output, err := executeCommand(deleteCmdInstance, "test-component-id")
		if err != nil {
			t.Fatalf("Expected no error, got %v. Output: %s", err, output)
		}

        var resp dtos.DeleteComponentOutput
        if errJson := json.Unmarshal([]byte(output), &resp); errJson != nil {
            t.Fatalf("Failed to unmarshal JSON output: %v. Output was: %s", errJson, output)
        }

		if !resp.IsSuccessful() { // Relies on IsSuccessful method of the DTO
			t.Errorf("Expected successful delete in JSON response (via IsSuccessful), got false. Output: %s", output)
		}
        if !strings.Contains(output, `"success": true`) { // Check raw JSON string
             t.Errorf("Expected output to contain '\"success\": true', got: %s", output)
        }
	})

	t.Run("delete with no ID", func(t *testing.T) {
		// For this test, we re-initialize deleteCmdInstance to ensure args are fresh for Cobra's parsing
		cmdInstanceForNoIDTest := Init()
		output, err := executeCommand(cmdInstanceForNoIDTest)
		if err == nil {
			t.Fatalf("Expected error for missing component ID, got nil. Output: %s", output)
		}
        // Cobra's error message for ExactArgs(1) might be "accepts 1 arg(s), received 0" or similar.
		if !strings.Contains(err.Error(), "accepts 1 arg(s)") && !strings.Contains(err.Error(), "requires at least 1 arg(s)") {
             t.Errorf("Expected error message about missing arguments, got: %v", err.Error())
        }
	})

    t.Run("service error", func(t *testing.T) {
		mockService.RunWithDTOsFunc = func(ctx context.Context, input compassservice.InputDTOInterface, output compassservice.OutputDTOInterface) error {
			return fmt.Errorf("internal service error")
		}

		output, err := executeCommand(deleteCmdInstance, "another-id")
		if err == nil {
			t.Fatalf("Expected an error from command execution, got nil. Output: %s", output)
		}
        if !strings.Contains(err.Error(), "internal service error") {
            t.Errorf("Expected error message containing 'internal service error', got: %v. Output: %s", err, output)
        }
	})
}
