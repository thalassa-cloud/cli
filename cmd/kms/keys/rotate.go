package keys

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var rotateCmd = &cobra.Command{
	Use:               "rotate <key>",
	Short:             "Rotate a KMS key on demand",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteKmsKeyIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		key, err := client.KMS().RotateKey(cmd.Context(), region, args[0])
		if err != nil {
			return fmt.Errorf("failed to rotate KMS key: %w", err)
		}

		body := [][]string{{
			key.Identity,
			key.Name,
			string(key.Status),
			fmt.Sprintf("%d", key.LatestVersion),
			formattime.FormatTime(key.UpdatedAt.Local(), showExactTime),
		}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Status", "Latest Version", "Updated"}, body)
		}
		return nil
	},
}

func init() {
	KeysCmd.AddCommand(rotateCmd)
	rotateCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	rotateCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
	rotateCmd.Flags().StringVar(&region, "region", "", "Region")
	_ = rotateCmd.MarkFlagRequired("region")
	_ = rotateCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
}
