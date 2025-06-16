package cmd // Changed from main

import (
	"fmt"
	"os"

	// Aliased imports for module commands to avoid collision if they also have 'cmd' package
	componentcmd "github.com/motain/of-catalog/internal/modules/component/cmd"
	metriccmd "github.com/motain/of-catalog/internal/modules/metric/cmd"
	scorecardcmd "github.com/motain/of-catalog/internal/modules/scorecard/cmd"
	"github.com/motain/of-catalog/internal/services/compassservice"
	"github.com/motain/of-catalog/internal/services/configservice"
	"github.com/spf13/cobra"
)

// Exported package-level variables for services
var (
	CfgService   configservice.ConfigServiceInterface
	CompassSvc compassservice.CompassServiceInterface
)

var RootCmd = &cobra.Command{ // Renamed to RootCmd for export
	Use:   "ofc",
	Short: "⚽ onefootball catalog CLI",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		CfgService = configservice.NewConfigService()

		gqlClient := compassservice.NewGraphQLClient(CfgService)
		httpClient := compassservice.NewHTTPClient(CfgService)
		// Basic nil check for httpClient, though a real app might have more robust error handling or fallbacks
		if httpClient == nil {
			return fmt.Errorf("failed to initialize HTTPClient, it is nil")
		}

		CompassSvc = compassservice.NewCompassService(CfgService, gqlClient, httpClient)
		if CompassSvc == nil {
			return fmt.Errorf("failed to initialize CompassService, it is nil")
		}
		return nil
	},
}

// Execute function is now part of the 'cmd' package
func Execute() {
	// Assumes InitCmd functions will be available in the respective module cmd packages
	RootCmd.AddCommand(componentcmd.InitCmd())
	RootCmd.AddCommand(metriccmd.InitCmd())
	RootCmd.AddCommand(scorecardcmd.InitCmd())

	// Define persistent flags applicable to all subcommands.
	RootCmd.PersistentFlags().StringP("file", "f", "", "Path to the resource definition YAML file")

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// main function removed, will be in a new main.go at project root
// Accessor functions are not strictly necessary if CfgService and CompassSvc are exported
// and the root cmd package is imported by subcommands.
// func GetCompassService() compassservice.CompassServiceInterface { return CompassSvc }
// func GetConfigService() configservice.ConfigServiceInterface { return CfgService }
