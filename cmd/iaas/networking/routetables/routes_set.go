package routetables

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
)

var routesSetFile string

var routesSetCmd = &cobra.Command{
	Use:               "set ROUTE_TABLE",
	Short:             "Replace all routes in a route table from a JSON file",
	Long:              "Batch-update routes for a route table. The file must contain a JSON array of UpdateRouteTableRoute objects.",
	Example:           `tcloud networking routetables routes set rt-123 --file routes.json`,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeRouteTableID,
	RunE: func(cmd *cobra.Command, args []string) error {
		if routesSetFile == "" {
			return fmt.Errorf("--file is required")
		}

		data, err := os.ReadFile(routesSetFile)
		if err != nil {
			return fmt.Errorf("failed to read routes file: %w", err)
		}

		var routes []iaas.UpdateRouteTableRoute
		if err := json.Unmarshal(data, &routes); err != nil {
			return fmt.Errorf("failed to parse routes file: %w", err)
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		updated, err := client.IaaS().UpdateRouteTableRoutes(cmd.Context(), args[0], iaas.UpdateRouteTableRoutes{
			Routes: routes,
		})
		if err != nil {
			return err
		}

		body := make([][]string, 0, len(updated))
		for _, route := range updated {
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
	RoutesCmd.AddCommand(routesSetCmd)

	routesSetCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	routesSetCmd.Flags().StringVar(&routesSetFile, "file", "", "JSON file containing an array of routes")
	_ = routesSetCmd.MarkFlagRequired("file")
}
