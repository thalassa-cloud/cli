package machines

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/config/contextstate"

	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/pkg/iaas"
	"github.com/thalassa-cloud/client-go/pkg/thalassa"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a machine",
	Long:    "Delete a machine. This command will delete the machine and all the services associated with it.",
	Aliases: []string{"d", "delete"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassa.NewClient(
			tcclient.WithBaseURL(contextstate.Server()),
			tcclient.WithOrganisation(contextstate.Organisation()),
			tcclient.WithAuthPersonalToken(contextstate.Token()),
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
			if tcclient.IsNotFound(err) {
				fmt.Println("Machine not found")
				return nil
			}
			return err
		}

		if machine.Status.Status == string(iaas.MachineStateStopped) {
			fmt.Println("Machine is already stopped")
			return nil
		}

		err = client.IaaS().DeleteMachine(cmd.Context(), machine.Identity)
		if err != nil {
			return err
		}
		fmt.Println("Machine is deleting...")
		if wait {
			// wait for machine to be deleted
			for {
				machine, err = client.IaaS().GetMachine(cmd.Context(), machine.Identity)
				if err != nil {
					if tcclient.IsNotFound(err) {
						break
					}
					return err
				}
				if machine.Status.Status == string(iaas.MachineStateDeleted) {
					break
				}
				time.Sleep(1 * time.Second)
			}
			fmt.Println("Machine deleted")
		}
		return nil
	},
}

func init() {
	MachinesCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVarP(&wait, "wait", "w", false, "Wait for the machine to be deleted")
}
