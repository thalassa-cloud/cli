package machines

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/config/contextstate"

	"github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/pkg/iaas"
	"github.com/thalassa-cloud/client-go/pkg/thalassa"
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
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassa.NewClient(
			client.WithBaseURL(contextstate.Server()),
			client.WithOrganisation(contextstate.Organisation()),
			client.WithAuthPersonalToken(contextstate.Token()),
		)
		if err != nil {
			return err
		}

		machineIdentity, err := getSelectedMachine(args)
		if err != nil {
			return err
		}

		machine, err := client.IaaS().GetMachine(cmd.Context(), machineIdentity)
		if err != nil {
			return err
		}

		if machine.Status.Status == string(iaas.MachineStateStopped) {
			fmt.Println("Machine is already stopped")
			return nil
		}

		err = client.IaaS().MachineStop(cmd.Context(), machine.Identity)
		if err != nil {
			return err
		}
		fmt.Println("Machine is stopping...")

		if wait {
			// wait for machine to be stopped
			for {
				machine, err = client.IaaS().GetMachine(cmd.Context(), machine.Identity)
				if err != nil {
					return err
				}
				if machine.Status.Status == string(iaas.MachineStateStopped) {
					break
				}
				time.Sleep(1 * time.Second)
			}
			fmt.Println("Machine stopped")
		}
		return nil
	},
}

func init() {
	MachinesCmd.AddCommand(stopCmd)

	stopCmd.Flags().BoolVarP(&wait, "wait", "w", false, "Wait for the machine to be stopped")
}
