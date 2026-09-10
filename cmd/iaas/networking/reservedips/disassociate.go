package reservedips

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var disassociateCmd = &cobra.Command{
	Use:     "disassociate",
	Short:   "Disassociate a reserved IP from its resource",
	Long:    "Detach a reserved IP from its currently associated load balancer or NAT gateway.",
	Example: "tcloud networking reserved-ips disassociate rip-123",
	Aliases: []string{"detach"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		rip, err := client.IaaS().DisassociateReservedIP(cmd.Context(), args[0])
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
	ReservedIPsCmd.AddCommand(disassociateCmd)
	disassociateCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	disassociateCmd.ValidArgsFunction = completeReservedIPID
}
