package routetables

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var routesListCmd = &cobra.Command{
	Use:               "list ROUTE_TABLE",
	Short:             "List routes in a route table",
	Aliases:           []string{"ls", "get"},
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeRouteTableID,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		rt, err := client.IaaS().GetRouteTable(cmd.Context(), args[0])
		if err != nil {
			return fmt.Errorf("failed to get route table: %w", err)
		}

		body := make([][]string, 0, len(rt.Routes))
		for _, route := range rt.Routes {
			body = append(body, []string{
				route.Identity,
				route.DestinationCidrBlock,
				routeTarget(route),
				route.Type,
			})
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Destination", "Target", "Type"}, body)
		}
		return nil
	},
}

func init() {
	RoutesCmd.AddCommand(routesListCmd)
	routesListCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
}
