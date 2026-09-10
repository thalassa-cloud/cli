package routetables

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	updateName        string
	updateDescription string
	updateLabels      []string
	updateAnnotations []string
)

var updateCmd = &cobra.Command{
	Use:               "update ROUTE_TABLE",
	Short:             "Update a route table",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeRouteTableID,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		current, err := client.IaaS().GetRouteTable(cmd.Context(), args[0])
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("route table not found: %s", args[0])
			}
			return fmt.Errorf("failed to get route table: %w", err)
		}

		req := iaas.UpdateRouteTable{
			Labels:      current.Labels,
			Annotations: current.Annotations,
		}
		if cmd.Flags().Changed("name") {
			req.Name = &updateName
		}
		if cmd.Flags().Changed("description") {
			req.Description = &updateDescription
		}
		if cmd.Flags().Changed("labels") {
			req.Labels = shared.KeyValuePairsToMap(updateLabels)
		}
		if cmd.Flags().Changed("annotations") {
			req.Annotations = shared.KeyValuePairsToMap(updateAnnotations)
		}

		rt, err := client.IaaS().UpdateRouteTable(cmd.Context(), current.Identity, req)
		if err != nil {
			return fmt.Errorf("failed to update route table: %w", err)
		}
		fmt.Printf("Route table %s updated\n", rt.Identity)
		return nil
	},
}

func init() {
	RouteTablesCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVar(&updateName, "name", "", "Name")
	updateCmd.Flags().StringVar(&updateDescription, "description", "", "Description")
	updateCmd.Flags().StringSliceVar(&updateLabels, "labels", nil, "Labels as key=value (repeatable)")
	updateCmd.Flags().StringSliceVar(&updateAnnotations, "annotations", nil, "Annotations as key=value (repeatable)")
}
