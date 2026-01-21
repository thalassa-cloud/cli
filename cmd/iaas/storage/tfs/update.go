package tfs

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/tfs"
)

var (
	updateTfsName             string
	updateTfsDescription      string
	updateTfsLabels           []string
	updateTfsAnnotations      []string
	updateTfsDeleteProtection bool
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:               "update",
	Short:             "Update a TFS instance",
	Long:              "Update properties of an existing TFS instance.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteTfsInstanceID,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		instanceIdentity := args[0]

		// Get current instance
		current, err := client.Tfs().GetTfsInstance(cmd.Context(), instanceIdentity)
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("TFS instance not found: %s", instanceIdentity)
			}
			return fmt.Errorf("failed to get TFS instance: %w", err)
		}

		req := tfs.UpdateTfsInstanceRequest{
			Name:             current.Name,
			Labels:           current.Labels,
			Annotations:      current.Annotations,
			DeleteProtection: current.DeleteProtection,
		}

		if current.Description != nil {
			req.Description = *current.Description
		}

		// Update name if provided
		if cmd.Flags().Changed("name") {
			req.Name = updateTfsName
		}

		// Update description if provided
		if cmd.Flags().Changed("description") {
			req.Description = updateTfsDescription
		}

		// Parse labels from key=value format
		if cmd.Flags().Changed("labels") {
			labels := make(map[string]string)
			for _, label := range updateTfsLabels {
				parts := strings.SplitN(label, "=", 2)
				if len(parts) == 2 {
					labels[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
				}
			}
			req.Labels = labels
		}

		// Parse annotations from key=value format
		if cmd.Flags().Changed("annotations") {
			annotations := make(map[string]string)
			for _, annotation := range updateTfsAnnotations {
				parts := strings.SplitN(annotation, "=", 2)
				if len(parts) == 2 {
					annotations[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
				}
			}
			req.Annotations = annotations
		}

		// Update delete protection if provided
		if cmd.Flags().Changed("delete-protection") {
			req.DeleteProtection = updateTfsDeleteProtection
		}

		instance, err := client.Tfs().UpdateTfsInstance(cmd.Context(), instanceIdentity, req)
		if err != nil {
			return fmt.Errorf("failed to update TFS instance: %w", err)
		}

		// Output in table format
		regionName := ""
		if instance.Region != nil {
			regionName = instance.Region.Name
		}
		body := [][]string{
			{
				instance.Identity,
				instance.Name,
				string(instance.Status),
				regionName,
			},
		}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Status", "Region"}, body)
		}

		return nil
	},
}

func init() {
	TfsCmd.AddCommand(updateCmd)

	updateCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	updateCmd.Flags().StringVar(&updateTfsName, "name", "", "Name of the TFS instance")
	updateCmd.Flags().StringVar(&updateTfsDescription, "description", "", "Description of the TFS instance")
	updateCmd.Flags().StringSliceVar(&updateTfsLabels, "labels", []string{}, "Labels in key=value format (can be specified multiple times)")
	updateCmd.Flags().StringSliceVar(&updateTfsAnnotations, "annotations", []string{}, "Annotations in key=value format (can be specified multiple times)")
	updateCmd.Flags().BoolVar(&updateTfsDeleteProtection, "delete-protection", false, "Enable or disable delete protection")
}
