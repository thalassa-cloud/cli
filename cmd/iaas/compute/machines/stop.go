package machines

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/client-go/iaas"
)

var (
	wait bool
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:     "stop",
	Short:   "Stop a machine",
	Long:    "Stop a machine to stop it from running. This command will stop the machine and all the services associated with it.",
	Aliases: []string{"s", "stop"},
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMachineStateChange(cmd, args, machineStateChange{
			alreadyState:    iaas.MachineStateStopped,
			alreadyMessage:  "Machine is already stopped",
			progressMessage: "Machine is stopping...",
			targetState:     iaas.MachineStateStopped,
			doneMessage:     "Machine stopped",
			change: func(ctx context.Context, iaasClient *iaas.Client, identity string) error {
				return iaasClient.MachineStop(ctx, identity)
			},
		})
	},
}

func init() {
	MachinesCmd.AddCommand(stopCmd)

	stopCmd.Flags().BoolVarP(&wait, "wait", "w", false, "Wait for the machine to be stopped")
}
