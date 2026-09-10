package routetables

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var deleteForce bool

var deleteCmd = &cobra.Command{
	Use:               "delete ROUTE_TABLE [ROUTE_TABLE...]",
	Aliases:           []string{"rm", "del"},
	Short:             "Delete route table(s)",
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: completeRouteTableID,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		var summary strings.Builder
		fmt.Fprintf(&summary, "Delete the following route table(s)?\n")
		for _, id := range args {
			rt, err := client.IaaS().GetRouteTable(cmd.Context(), id)
			if err != nil {
				if tcclient.IsNotFound(err) {
					fmt.Printf("Route table %s not found\n", id)
					continue
				}
				return err
			}
			fmt.Fprintf(&summary, "  %s (%s)\n", rt.Name, rt.Identity)
		}

		proceed, err := shared.PromptDestructiveUnlessForce(deleteForce, summary.String())
		if err != nil {
			return err
		}
		if !proceed {
			return nil
		}

		for _, id := range args {
			if err := client.IaaS().DeleteRouteTable(cmd.Context(), id); err != nil {
				return fmt.Errorf("failed to delete route table %s: %w", id, err)
			}
			fmt.Printf("Route table %s deleted\n", id)
		}
		return nil
	},
}

func init() {
	RouteTablesCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVar(&deleteForce, shared.ForceKey, false, "Skip confirmation")
}
