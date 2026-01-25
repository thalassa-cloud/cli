package machines

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/thalassaclient"
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

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		machineIdentity, err := getSelectedMachine(args)
		if err != nil {
			return err
		}

		machine, err := client.IaaS().GetMachine(cmd.Context(), machineIdentity)
		if err != nil {
			return err
		}

		if machine.Status.Status == string(iaas.MachineStateRunning) {
			fmt.Println("Machine is already running")
			return nil
		}

		err = client.IaaS().MachineStart(cmd.Context(), machine.Identity)
		if err != nil {
			return err
		}
		fmt.Println("Machine is starting...")

		if wait {
			// wait for machine to be started
			for {
				machine, err = client.IaaS().GetMachine(cmd.Context(), machine.Identity)
				if err != nil {
					return err
				}
				if machine.Status.Status == string(iaas.MachineStateRunning) {
					break
				}
				time.Sleep(1 * time.Second)
			}
			fmt.Println("Machine started")
		}
		return nil
	},
}

func init() {
	MachinesCmd.AddCommand(startCmd)

	startCmd.Flags().BoolVarP(&wait, "wait", "w", false, "Wait for the machine to be started")
}
