package cmd

import (
	"github.com/motain/of-catalog/internal/modules/component/cmd/apply"
	"github.com/motain/of-catalog/internal/modules/component/cmd/bind"
	"github.com/motain/of-catalog/internal/modules/component/cmd/compute"
	"github.com/spf13/cobra"
)

func Init() *cobra.Command {
	componentCmd := &cobra.Command{
		Use:   "component",
		Short: "component related commands",
	}

	componentCmd.AddCommand(apply.Init())
	componentCmd.AddCommand(bind.Init())
	componentCmd.AddCommand(compute.Init())

	return componentCmd
}
