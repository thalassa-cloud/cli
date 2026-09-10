package reservedips

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/table"
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
	Use:     "update",
	Short:   "Update a reserved IP address",
	Long:    "Update metadata of an existing reserved IP address. Unspecified fields are preserved from the current resource.",
	Example: "tcloud networking reserved-ips update rip-123 --name new-name\ntcloud networking reserved-ips update rip-123 --description 'updated'",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		current, err := client.IaaS().GetReservedIP(cmd.Context(), args[0])
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("reserved IP not found: %s", args[0])
			}
			return fmt.Errorf("failed to get reserved IP: %w", err)
		}

		req := iaas.UpdateReservedIpRequest{
			Name:        current.Name,
			Description: current.Description,
			Labels:      current.Labels,
			Annotations: current.Annotations,
		}
		if cmd.Flags().Changed("name") {
			req.Name = updateName
		}
		if cmd.Flags().Changed("description") {
			req.Description = updateDescription
		}
		if cmd.Flags().Changed("labels") {
			req.Labels = parseKeyValueSlice(updateLabels)
		}
		if cmd.Flags().Changed("annotations") {
			req.Annotations = parseKeyValueSlice(updateAnnotations)
		}

		rip, err := client.IaaS().UpdateReservedIP(cmd.Context(), current.Identity, req)
		if err != nil {
			return err
		}

		body := [][]string{{
			rip.Identity,
			rip.Name,
			regionName(*rip),
			string(rip.Status),
			dashIfEmpty(rip.IPv4Address),
			dashIfEmpty(rip.IPv6Address),
			attachedTo(*rip),
		}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Region", "Status", "IPv4", "IPv6", "AttachedTo"}, body)
		}
		return nil
	},
}

func init() {
	ReservedIPsCmd.AddCommand(updateCmd)

	updateCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	updateCmd.Flags().StringVar(&updateName, "name", "", "Name of the reserved IP")
	updateCmd.Flags().StringVar(&updateDescription, "description", "", "Description of the reserved IP")
	updateCmd.Flags().StringSliceVar(&updateLabels, "labels", []string{}, "Labels in key=value format")
	updateCmd.Flags().StringSliceVar(&updateAnnotations, "annotations", []string{}, "Annotations in key=value format")

	updateCmd.ValidArgsFunction = completeReservedIPID
}
