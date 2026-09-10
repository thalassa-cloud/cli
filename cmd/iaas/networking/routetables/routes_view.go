package routetables

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var routesViewOutputFormat string

var routesViewCmd = &cobra.Command{
	Use:               "view ROUTE_TABLE ROUTE",
	Short:             "View a route in a route table",
	Aliases:           []string{"show", "describe"},
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

		if routesViewOutputFormat == "yaml" {
			yamlData, err := yaml.Marshal(route)
			if err != nil {
				return fmt.Errorf("failed to marshal to YAML: %w", err)
			}
			fmt.Print(string(yamlData))
			return nil
		}

		fmt.Printf("Route Details:\n")
		fmt.Printf("  ID: %s\n", route.Identity)
		fmt.Printf("  Destination: %s\n", route.DestinationCidrBlock)
		fmt.Printf("  Target: %s\n", routeTarget(*route))
		fmt.Printf("  Type: %s\n", route.Type)
		if route.Note != nil {
			fmt.Printf("  Note: %s\n", *route.Note)
		}
		return nil
	},
}

func init() {
	RoutesCmd.AddCommand(routesViewCmd)
	routesViewCmd.Flags().StringVarP(&routesViewOutputFormat, "output", "o", "", "Output format (yaml)")
	_ = routesViewCmd.RegisterFlagCompletionFunc("output", completeOutputFormat)
}
