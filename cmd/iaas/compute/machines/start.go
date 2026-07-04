package machines

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/client-go/iaas"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:     "start",
	Short:   "Start a machine",
	Long:    "Start a machine to start it from stopped state. This command will start the machine and all the services associated with it.",
	Aliases: []string{"s", "start"},
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runMachineStateChange(cmd, args, machineStateChange{
			alreadyState:    iaas.MachineStateRunning,
			alreadyMessage:  "Machine is already running",
			progressMessage: "Machine is starting...",
			targetState:     iaas.MachineStateRunning,
			doneMessage:     "Machine started",
			change: func(ctx context.Context, iaasClient *iaas.Client, identity string) error {
				return iaasClient.MachineStart(ctx, identity)
			},
		})
	},
}

func init() {
	MachinesCmd.AddCommand(startCmd)

	startCmd.Flags().BoolVarP(&wait, "wait", "w", false, "Wait for the machine to be started")
}
