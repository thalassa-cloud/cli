package compute

import (
	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/cmd/iaas/compute/machines"
)

// ComputeCmd represents the compute command
var ComputeCmd = &cobra.Command{
	Use:     "compute",
	Aliases: []string{"comp", "c"},
	Short:   "Manage compute resources",
	Long:    "Manage compute resources in the Thalassa Cloud Platform",
}

func init() {
	ComputeCmd.AddCommand(machines.MachinesCmd)
}
