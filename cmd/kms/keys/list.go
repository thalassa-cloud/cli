package keys

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientkms "github.com/thalassa-cloud/client-go/kms"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List KMS keys in a region",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		keys, err := client.KMS().ListKeys(cmd.Context(), region, &clientkms.ListKeysRequest{})
		if err != nil {
			return fmt.Errorf("failed to list KMS keys: %w", err)
		}

		body := make([][]string, 0, len(keys))
		for _, k := range keys {
			body = append(body, []string{
				k.Identity,
				k.Name,
				string(k.KeyType),
				string(k.Status),
				fmt.Sprintf("%v", k.KeyRotationEnabled),
				formattime.FormatTime(k.CreatedAt.Local(), showExactTime),
			})
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Type", "Status", "Rotation", "Created"}, body)
		}
		return nil
	},
}

func init() {
	KeysCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	listCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
	listCmd.Flags().StringVar(&region, "region", "", "Region")
	_ = listCmd.MarkFlagRequired("region")
	_ = listCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
}
