package secrets

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/fzf"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var (
	browseRegion string
	browsePath   string
)

var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "Browse secret prefixes and secrets at a path",
	Long: `Browse Secrets Manager prefixes and secrets.

In a terminal with fzf available, browse is interactive: select prefixes to
descend, ".." to go up, and secrets to view metadata (optionally reveal values
after confirmation). Press Esc to quit.

When stdout is not a terminal, or TC_IGNORE_FZF is set, prints a non-interactive
table for the given --path.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		if fzf.IsInteractiveMode(os.Stdout) {
			return runInteractiveBrowse(cmd.Context(), cmd.OutOrStdout(), client.Secrets(), browseRegion, browsePath)
		}

		result, err := client.Secrets().BrowseSecrets(cmd.Context(), browseRegion, browsePath)
		if err != nil {
			return fmt.Errorf("failed to browse secrets: %w", err)
		}

		if len(result.Prefixes) > 0 {
			prefixBody := make([][]string, 0, len(result.Prefixes))
			for _, p := range result.Prefixes {
				prefixBody = append(prefixBody, []string{p})
			}
			if noHeader {
				table.Print(nil, prefixBody)
			} else {
				fmt.Println("Prefixes:")
				table.Print([]string{"Prefix"}, prefixBody)
			}
		}

		if len(result.Secrets) > 0 {
			if len(result.Prefixes) > 0 {
				fmt.Println()
			}
			secretBody := make([][]string, 0, len(result.Secrets))
			for _, s := range result.Secrets {
				kmsKey := "-"
				if s.KmsKey != nil {
					kmsKey = s.KmsKey.Identity
					if kmsKey == "" {
						kmsKey = s.KmsKey.Name
					}
				}
				secretBody = append(secretBody, []string{
					s.Path,
					fmt.Sprintf("%d", s.CurrentVersion),
					kmsKey,
					formattime.FormatTime(s.UpdatedAt.Local(), showExactTime),
				})
			}
			if noHeader {
				table.Print(nil, secretBody)
			} else {
				fmt.Println("Secrets:")
				table.Print([]string{"Path", "Version", "KMS Key", "Updated"}, secretBody)
			}
		}
		return nil
	},
}

func init() {
	SecretsCmd.AddCommand(browseCmd)
	browseCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	browseCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
	browseCmd.Flags().StringVar(&browseRegion, "region", "", "Region")
	browseCmd.Flags().StringVar(&browsePath, "path", "/", "Path to browse")
	_ = browseCmd.MarkFlagRequired("region")
	_ = browseCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
	_ = browseCmd.RegisterFlagCompletionFunc("path", completion.CompleteSecretPath)
}
