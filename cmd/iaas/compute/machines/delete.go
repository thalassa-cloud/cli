package machines

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/labels"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/iaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	deleteWait    bool
	force         bool
	labelSelector string
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete machine(s)",
	Long:    "Delete machine(s) by identity or label selector. This command will delete the machine(s) and all the services associated with it.",
	Example: "tcloud compute machines delete vm-123\ntcloud compute machines delete vm-123 vm-456 --wait\ntcloud compute machines delete --selector environment=test --force",
	Aliases: []string{"d", "del", "remove"},
	Args:    cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && labelSelector == "" {
			// Try interactive selection if no args and no selector
			machineIdentity, err := getSelectedMachine(args)
			if err != nil {
				return fmt.Errorf("either machine identity(ies) or --selector must be provided")
			}
			args = []string{machineIdentity}
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Collect machines to delete
		machinesToDelete := []iaas.Machine{}

		// If label selector is provided, filter by labels
		if labelSelector != "" {
			allMachines, err := client.IaaS().ListMachines(cmd.Context(), &iaas.ListMachinesRequest{
				Filters: []filters.Filter{
					&filters.LabelFilter{MatchLabels: labels.ParseLabelSelector(labelSelector)},
				},
			})
			if err != nil {
				return fmt.Errorf("failed to list machines: %w", err)
			}
			if len(allMachines) == 0 {
				fmt.Println("No machines found matching the label selector")
				return nil
			}
			machinesToDelete = append(machinesToDelete, allMachines...)
		} else {
			// Get machines by identity
			for _, machineIdentity := range args {
				machine, err := client.IaaS().GetMachine(cmd.Context(), machineIdentity)
				if err != nil {
					if tcclient.IsNotFound(err) {
						fmt.Printf("Machine %s not found\n", machineIdentity)
						continue
					}
					return fmt.Errorf("failed to get machine: %w", err)
				}
				machinesToDelete = append(machinesToDelete, *machine)
			}
		}

		if len(machinesToDelete) == 0 {
			fmt.Println("No machines to delete")
			return nil
		}

		// Ask for confirmation unless --force is provided
		if !force {
			fmt.Printf("Are you sure you want to delete the following machine(s)?\n")
			for _, machine := range machinesToDelete {
				fmt.Printf("  %s (%s)\n", machine.Name, machine.Identity)
			}
			var confirm string
			fmt.Printf("Enter 'yes' to confirm: ")
			fmt.Scanln(&confirm)
			if confirm != "yes" {
				fmt.Println("Aborted")
				return nil
			}
		}

		// Delete each machine
		for _, machine := range machinesToDelete {
			if machine.Status.Status == string(iaas.MachineStateStopped) {
				fmt.Printf("Machine %s is already stopped\n", machine.Identity)
				continue
			}

			fmt.Printf("Deleting machine: %s (%s)\n", machine.Name, machine.Identity)
			err := client.IaaS().DeleteMachine(cmd.Context(), machine.Identity)
			if err != nil {
				return fmt.Errorf("failed to delete machine: %w", err)
			}

			if deleteWait {
				if err := client.IaaS().WaitUntilMachineDeleted(cmd.Context(), machine.Identity); err != nil {
					return fmt.Errorf("failed to wait for machine to be deleted: %w", err)
				}
			}
			fmt.Printf("Machine %s deleted successfully\n", machine.Identity)
		}

		return nil
	},
}

func init() {
	MachinesCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVarP(&deleteWait, "wait", "w", false, "Wait for the machine(s) to be deleted")
	deleteCmd.Flags().BoolVar(&force, "force", false, "Force the deletion and skip the confirmation")
	deleteCmd.Flags().StringVarP(&labelSelector, "selector", "l", "", "Label selector to filter machines (format: key1=value1,key2=value2)")
}
