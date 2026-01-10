package vpcpeering

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
)

const (
	UpdateFlagName        = "name"
	UpdateFlagDescription = "description"
	UpdateFlagLabels      = "labels"
	UpdateFlagAnnotations = "annotations"
)

var (
	updateName        string
	updateDescription string
	updateLabels      []string
	updateAnnotations []string
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a VPC peering connection",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		connectionIdentity := args[0]

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Get current connection to preserve values if not provided
		current, err := client.IaaS().GetVpcPeeringConnection(cmd.Context(), connectionIdentity)
		if err != nil {
			return fmt.Errorf("failed to get VPC peering connection: %w", err)
		}

		req := iaas.UpdateVpcPeeringConnectionRequest{
			Name:        current.Name,
			Description: current.Description,
			Labels:      current.Labels,
			Annotations: current.Annotations,
		}

		// Update only provided fields
		if cmd.Flags().Changed(UpdateFlagName) {
			req.Name = updateName
		}
		if cmd.Flags().Changed(UpdateFlagDescription) {
			req.Description = updateDescription
		}

		// Parse labels from key=value format
		if cmd.Flags().Changed(UpdateFlagLabels) {
			labels := make(map[string]string)
			for _, label := range updateLabels {
				parts := strings.SplitN(label, "=", 2)
				if len(parts) == 2 {
					labels[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
				}
			}
			req.Labels = labels
		}

		// Parse annotations from key=value format
		if cmd.Flags().Changed(UpdateFlagAnnotations) {
			annotations := make(map[string]string)
			for _, annotation := range updateAnnotations {
				parts := strings.SplitN(annotation, "=", 2)
				if len(parts) == 2 {
					annotations[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
				}
			}
			req.Annotations = annotations
		}

		connection, err := client.IaaS().UpdateVpcPeeringConnection(cmd.Context(), connectionIdentity, req)
		if err != nil {
			return fmt.Errorf("failed to update VPC peering connection: %w", err)
		}

		body := make([][]string, 0, 1)
		body = append(body, []string{
			connection.Identity,
			connection.Name,
			string(connection.Status),
			formattime.FormatTime(connection.UpdatedAt.Local(), false),
		})
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Status", "Updated"}, body)
		}
		return nil
	},
}

func init() {
	VpcPeeringCmd.AddCommand(updateCmd)

	updateCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	updateCmd.Flags().StringVar(&updateName, UpdateFlagName, "", "Name of the VPC peering connection")
	updateCmd.Flags().StringVar(&updateDescription, UpdateFlagDescription, "", "Description of the VPC peering connection")
	updateCmd.Flags().StringSliceVar(&updateLabels, UpdateFlagLabels, []string{}, "Labels in key=value format")
	updateCmd.Flags().StringSliceVar(&updateAnnotations, UpdateFlagAnnotations, []string{}, "Annotations in key=value format")
}
