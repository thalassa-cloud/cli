package keys

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var viewCmd = &cobra.Command{
	Use:               "view <key>",
	Aliases:           []string{"get", "show"},
	Short:             "View a KMS key",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteKmsKeyIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		key, err := client.KMS().GetKey(cmd.Context(), region, args[0])
		if err != nil {
			return fmt.Errorf("failed to get KMS key: %w", err)
		}

		rotationPeriod := "-"
		if key.RotationPeriodInDays != nil {
			rotationPeriod = fmt.Sprintf("%d", *key.RotationPeriodInDays)
		}

		body := [][]string{{
			key.Identity,
			key.Name,
			string(key.KeyType),
			string(key.Status),
			fmt.Sprintf("%v", key.ExportAllowed),
			fmt.Sprintf("%v", key.KeyRotationEnabled),
			rotationPeriod,
			fmt.Sprintf("%d", key.LatestVersion),
			formatMap(key.Labels),
			formattime.FormatTime(key.CreatedAt.Local(), showExactTime),
		}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Type", "Status", "Exportable", "Rotation", "Period Days", "Latest Version", "Labels", "Created"}, body)
		}
		return nil
	},
}

func formatMap(m map[string]string) string {
	if len(m) == 0 {
		return "-"
	}
	parts := make([]string, 0, len(m))
	for k, v := range m {
		parts = append(parts, k+"="+v)
	}
	sort.Strings(parts)
	return strings.Join(parts, ",")
}

func init() {
	KeysCmd.AddCommand(viewCmd)
	viewCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	viewCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
	viewCmd.Flags().StringVar(&region, "region", "", "Region")
	_ = viewCmd.MarkFlagRequired("region")
	_ = viewCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
}
