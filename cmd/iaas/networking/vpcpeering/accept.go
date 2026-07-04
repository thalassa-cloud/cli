package vpcpeering

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
)

var (
	acceptForce bool
)

// acceptCmd represents the accept command
var acceptCmd = &cobra.Command{
	Use:     "accept",
	Short:   "Accept a VPC peering connection",
	Long:    "Accept a pending VPC peering connection request",
	Example: "tcloud networking vpc-peering accept vpcpc-123\ntcloud networking vpc-peering accept vpcpc-123 --force",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		connectionIdentity := args[0]

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Get connection details for confirmation
		connection, err := client.IaaS().GetVpcPeeringConnection(cmd.Context(), connectionIdentity)
		if err != nil {
			return fmt.Errorf("failed to get VPC peering connection: %w", err)
		}

		// Ask for confirmation unless --force is provided
		var summary strings.Builder
		fmt.Fprintf(&summary, "Are you sure you want to accept the following VPC peering connection?\n")
		fmt.Fprintf(&summary, "  ID: %s\n", connection.Identity)
		fmt.Fprintf(&summary, "  Name: %s\n", connection.Name)
		fmt.Fprintf(&summary, "  Status: %s\n", connection.Status)
		if connection.RequesterVpc != nil {
			fmt.Fprintf(&summary, "  Requester VPC: %s (%s)\n", connection.RequesterVpc.Name, connection.RequesterVpc.Identity)
		}
		if connection.AccepterVpc != nil {
			fmt.Fprintf(&summary, "  Accepter VPC: %s (%s)\n", connection.AccepterVpc.Name, connection.AccepterVpc.Identity)
		}
		proceed, err := shared.PromptDestructiveUnlessForce(acceptForce, summary.String())
		if err != nil {
			return err
		}
		if !proceed {
			return nil
		}

		req := iaas.AcceptVpcPeeringConnectionRequest{}

		connection, err = client.IaaS().AcceptVpcPeeringConnection(cmd.Context(), connectionIdentity, req)
		if err != nil {
			return fmt.Errorf("failed to accept VPC peering connection: %w", err)
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
	VpcPeeringCmd.AddCommand(acceptCmd)

	acceptCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	acceptCmd.Flags().BoolVar(&acceptForce, "force", false, "Force the acceptance and skip the confirmation")

	// Add completion
	acceptCmd.ValidArgsFunction = completion.CompleteVpcPeeringConnectionID
}
