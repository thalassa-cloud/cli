package zones

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var deleteForce bool

var deleteCmd = &cobra.Command{
	Use:               "delete <zone>",
	Aliases:           []string{"rm", "remove"},
	Short:             "Delete a DNS zone and all of its records",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteDnsZoneIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		zoneIdentity := args[0]
		zone, err := client.DNS().GetZone(cmd.Context(), zoneIdentity)
		if err != nil {
			return fmt.Errorf("failed to get DNS zone: %w", err)
		}

		ok, err := shared.PromptDestructiveUnlessForce(deleteForce, fmt.Sprintf(
			"Are you sure you want to delete this DNS zone and all of its records?\n  ID: %s\n  Name: %s\n",
			zone.Identity, zone.Name,
		))
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}

		if err := client.DNS().DeleteZone(cmd.Context(), zone.Identity); err != nil {
			return fmt.Errorf("failed to delete DNS zone: %w", err)
		}
		fmt.Printf("Deleted DNS zone %s\n", zone.Name)
		return nil
	},
}

func init() {
	ZonesCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVar(&deleteForce, shared.ForceKey, false, "Skip the confirmation prompt and delete")
}
