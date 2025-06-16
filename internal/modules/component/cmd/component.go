package cmd

import (
	"github.com/motain/of-catalog/internal/modules/component/cmd/create"
	"github.com/motain/of-catalog/internal/modules/component/cmd/delete" // Import the new delete package
	"github.com/motain/of-catalog/internal/modules/component/cmd/read"
	"github.com/motain/of-catalog/internal/modules/component/cmd/update"
	"github.com/spf13/cobra"
)

// InitCmd initializes the component command and its subcommands.
func InitCmd() *cobra.Command { // Renamed from Init
	componentCmd := &cobra.Command{
		Use:   "component",
		Short: "component related commands",
		// The RunE for componentCmd might need to be defined if 'ofc component' itself is a valid command.
		// If 'ofc component' should always be followed by a subcommand (create, read, etc.),
		// then a RunE is not strictly necessary here.
	}

	// Add new CRUD subcommands
	componentCmd.AddCommand(create.Init())
	componentCmd.AddCommand(read.Init())
	componentCmd.AddCommand(update.Init())
	componentCmd.AddCommand(delete.Init()) // Add the delete command

	// Remove old subcommands if they are no longer needed
	// componentCmd.AddCommand(apply.Init())
	// componentCmd.AddCommand(bind.Init())
	// componentCmd.AddCommand(compute.Init())

	return componentCmd
}
