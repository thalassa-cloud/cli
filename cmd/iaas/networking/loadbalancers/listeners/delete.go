package listeners

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var deleteForce bool

var deleteCmd = &cobra.Command{
	Use:     "delete LISTENER [LISTENER...]",
	Short:   "Delete listener(s)",
	Long:    "Delete one or more listeners from a load balancer.",
	Example: "tcloud networking loadbalancers listeners delete listener-123 --loadbalancer lb-123 --force",
	Aliases: []string{"d", "del", "remove", "rm"},
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if loadbalancer == "" {
			return fmt.Errorf("--loadbalancer is required")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		if !deleteForce {
			fmt.Printf("Are you sure you want to delete %d listener(s) from load balancer %s?\n", len(args), loadbalancer)
			var confirm string
			fmt.Printf("Enter 'yes' to confirm: ")
			fmt.Scanln(&confirm)
			if confirm != "yes" {
				fmt.Println("Aborted")
				return nil
			}
		}

		for _, listenerID := range args {
			fmt.Printf("Deleting listener: %s\n", listenerID)
			if err := client.IaaS().DeleteListener(cmd.Context(), loadbalancer, listenerID); err != nil {
				if tcclient.IsNotFound(err) {
					fmt.Printf("Listener %s not found\n", listenerID)
					continue
				}
				return fmt.Errorf("failed to delete listener: %w", err)
			}
			fmt.Printf("Listener %s deleted successfully\n", listenerID)
		}

		return nil
	},
}

func init() {
	ListenersCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVar(&loadbalancer, LoadbalancerFlag, "", "Load balancer identity")
	deleteCmd.Flags().BoolVar(&deleteForce, "force", false, "Skip confirmation")

	deleteCmd.MarkFlagRequired(LoadbalancerFlag)
	deleteCmd.ValidArgsFunction = completeLoadbalancerListenerID
	deleteCmd.RegisterFlagCompletionFunc(LoadbalancerFlag, completeLoadbalancerID)
}
