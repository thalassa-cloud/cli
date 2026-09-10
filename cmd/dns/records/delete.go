package records

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var deleteForce bool

var deleteCmd = &cobra.Command{
	Use:     "delete <record>",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a DNS record",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		recordIdentity := args[0]
		record, err := client.DNS().GetRecord(cmd.Context(), zoneIdentity, recordIdentity)
		if err != nil {
			return fmt.Errorf("failed to get DNS record: %w", err)
		}

		ok, err := shared.PromptDestructiveUnlessForce(deleteForce, fmt.Sprintf(
			"Are you sure you want to delete this DNS record?\n  ID: %s\n  Name: %s\n  Type: %s\n",
			record.Identity, record.Name, record.Type,
		))
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}

		if err := client.DNS().DeleteRecord(cmd.Context(), zoneIdentity, record.Identity); err != nil {
			return fmt.Errorf("failed to delete DNS record: %w", err)
		}
		fmt.Printf("Deleted DNS record %s\n", record.Name)
		return nil
	},
}

func init() {
	RecordsCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVar(&deleteForce, shared.ForceKey, false, "Skip the confirmation prompt and delete")
	deleteCmd.Flags().StringVar(&zoneIdentity, "zone", "", "DNS zone identity")
	_ = deleteCmd.MarkFlagRequired("zone")
	_ = deleteCmd.RegisterFlagCompletionFunc("zone", completion.CompleteDnsZoneIdentity)
}
