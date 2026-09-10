package secrets

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var (
	viewRegion   string
	viewPath     string
	viewVersions bool
)

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View secret metadata (does not print secret material)",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		secret, err := client.Secrets().GetSecret(cmd.Context(), viewRegion, viewPath, viewVersions)
		if err != nil {
			return fmt.Errorf("failed to get secret: %w", err)
		}

		kmsKey := "-"
		if secret.KmsKey != nil {
			kmsKey = secret.KmsKey.Identity
			if kmsKey == "" {
				kmsKey = secret.KmsKey.Name
			}
		}

		body := [][]string{
			{"Path", secret.Path},
			{"Description", secret.Description},
			{"Current Version", fmt.Sprintf("%d", secret.CurrentVersion)},
			{"KMS Key", kmsKey},
			{"Created", formattime.FormatTime(secret.CreatedAt.Local(), showExactTime)},
			{"Updated", formattime.FormatTime(secret.UpdatedAt.Local(), showExactTime)},
		}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"Field", "Value"}, body)
		}

		if viewVersions && len(secret.Versions) > 0 {
			fmt.Println()
			versionBody := make([][]string, 0, len(secret.Versions))
			for _, v := range secret.Versions {
				destroyed := "-"
				if v.DestroyedAt != nil {
					destroyed = formattime.FormatTime(v.DestroyedAt.Local(), showExactTime)
				}
				versionBody = append(versionBody, []string{
					fmt.Sprintf("%d", v.Version),
					v.Status,
					formattime.FormatTime(v.CreatedAt.Local(), showExactTime),
					destroyed,
				})
			}
			if noHeader {
				table.Print(nil, versionBody)
			} else {
				table.Print([]string{"Version", "Status", "Created", "Destroyed"}, versionBody)
			}
		}
		return nil
	},
}

func init() {
	SecretsCmd.AddCommand(viewCmd)
	viewCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	viewCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
	viewCmd.Flags().StringVar(&viewRegion, "region", "", "Region")
	viewCmd.Flags().StringVar(&viewPath, "path", "", "Secret path")
	viewCmd.Flags().BoolVar(&viewVersions, "versions", false, "Include version history")
	_ = viewCmd.MarkFlagRequired("region")
	_ = viewCmd.MarkFlagRequired("path")
	_ = viewCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
	_ = viewCmd.RegisterFlagCompletionFunc("path", completion.CompleteSecretPath)
}
