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
	listRegion string
	listPrefix string
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List secrets under a path prefix (metadata only)",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		secrets, err := client.Secrets().ListSecrets(cmd.Context(), listRegion, listPrefix)
		if err != nil {
			return fmt.Errorf("failed to list secrets: %w", err)
		}

		body := make([][]string, 0, len(secrets))
		for _, s := range secrets {
			kmsKey := "-"
			if s.KmsKey != nil {
				kmsKey = s.KmsKey.Identity
				if kmsKey == "" {
					kmsKey = s.KmsKey.Name
				}
			}
			body = append(body, []string{
				s.Path,
				s.Description,
				fmt.Sprintf("%d", s.CurrentVersion),
				kmsKey,
				formattime.FormatTime(s.CreatedAt.Local(), showExactTime),
			})
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"Path", "Description", "Version", "KMS Key", "Created"}, body)
		}
		return nil
	},
}

func init() {
	SecretsCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	listCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
	listCmd.Flags().StringVar(&listRegion, "region", "", "Region")
	listCmd.Flags().StringVar(&listPrefix, "prefix", "/", "Path prefix")
	_ = listCmd.MarkFlagRequired("region")
	_ = listCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
}
