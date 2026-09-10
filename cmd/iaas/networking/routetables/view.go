package routetables

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var viewExactTime bool

var viewCmd = &cobra.Command{
	Use:               "view ROUTE_TABLE",
	Short:             "View a route table and its routes",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeRouteTableID,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		rt, err := client.IaaS().GetRouteTable(cmd.Context(), args[0])
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("route table not found: %s", args[0])
			}
			return fmt.Errorf("failed to get route table: %w", err)
		}

		body := [][]string{
			{"ID", rt.Identity},
			{"Name", rt.Name},
			{"Slug", rt.Slug},
			{"Description", descriptionString(rt.Description)},
			{"VPC", vpcName(*rt)},
			{"Default", fmt.Sprintf("%t", rt.IsDefault)},
			{"Created", formattime.FormatTime(rt.CreatedAt.Local(), viewExactTime)},
		}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"Field", "Value"}, body)
		}

		if len(rt.Routes) > 0 {
			fmt.Println()
			routeBody := make([][]string, 0, len(rt.Routes))
			for _, route := range rt.Routes {
				routeBody = append(routeBody, []string{
					route.Identity,
					route.DestinationCidrBlock,
					routeTarget(route),
					route.Type,
				})
			}
			if noHeader {
				table.Print(nil, routeBody)
			} else {
				table.Print([]string{"ID", "Destination", "Target", "Type"}, routeBody)
			}
		}
		return nil
	},
}

func init() {
	RouteTablesCmd.AddCommand(viewCmd)
	viewCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	viewCmd.Flags().BoolVar(&viewExactTime, "exact-time", false, "Show full timestamps instead of relative time")
}
