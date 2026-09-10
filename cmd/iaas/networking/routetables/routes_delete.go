package routetables

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var routesDeleteForce bool

var routesDeleteCmd = &cobra.Command{
	Use:               "delete ROUTE_TABLE ROUTE",
	Short:             "Delete a route from a route table",
	Aliases:           []string{"d", "del", "remove", "rm"},
	Args:              cobra.ExactArgs(2),
	ValidArgsFunction: completeRouteTableID,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		route, err := client.IaaS().GetRouteTableRoute(cmd.Context(), args[0], args[1])
		if err != nil {
			return fmt.Errorf("failed to get route: %w", err)
		}

		var summary strings.Builder
		fmt.Fprintf(&summary, "Are you sure you want to delete route %s (%s)?\n", route.Identity, route.DestinationCidrBlock)
		proceed, err := shared.PromptDestructiveUnlessForce(routesDeleteForce, summary.String())
		if err != nil {
			return err
		}
		if !proceed {
			return nil
		}

		if err := client.IaaS().DeleteRouteTableRoute(cmd.Context(), args[0], args[1]); err != nil {
			return fmt.Errorf("failed to delete route: %w", err)
		}
		fmt.Printf("Route %s deleted successfully\n", args[1])
		return nil
	},
}

func init() {
	RoutesCmd.AddCommand(routesDeleteCmd)
	routesDeleteCmd.Flags().BoolVar(&routesDeleteForce, "force", false, "Force the deletion and skip the confirmation")
}
