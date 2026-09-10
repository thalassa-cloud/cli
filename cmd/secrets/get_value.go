package secrets

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var (
	getValueRegion  string
	getValuePath    string
	getValueVersion int
)

var getValueCmd = &cobra.Command{
	Use:   "get-value",
	Short: "Print secret material to stdout",
	Long:  "Fetches and prints the secret value. Prefer piping to a file and avoid logging the output.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		var version *int
		if cmd.Flags().Changed("version") {
			v := getValueVersion
			version = &v
		}

		result, err := client.Secrets().GetSecretValue(cmd.Context(), getValueRegion, getValuePath, version)
		if err != nil {
			return fmt.Errorf("failed to get secret value: %w", err)
		}

		switch {
		case result.SecretString != "":
			fmt.Print(result.SecretString)
			if len(result.SecretString) == 0 || result.SecretString[len(result.SecretString)-1] != '\n' {
				fmt.Println()
			}
		case len(result.SecretKeyValues) > 0:
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetIndent("", "  ")
			if err := enc.Encode(result.SecretKeyValues); err != nil {
				return err
			}
		default:
			fmt.Println("{}")
		}
		return nil
	},
}

func init() {
	SecretsCmd.AddCommand(getValueCmd)
	getValueCmd.Flags().StringVar(&getValueRegion, "region", "", "Region")
	getValueCmd.Flags().StringVar(&getValuePath, "path", "", "Secret path")
	getValueCmd.Flags().IntVar(&getValueVersion, "version", 0, "Specific version (defaults to current)")
	_ = getValueCmd.MarkFlagRequired("region")
	_ = getValueCmd.MarkFlagRequired("path")
	_ = getValueCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
	_ = getValueCmd.RegisterFlagCompletionFunc("path", completion.CompleteSecretPath)
	_ = getValueCmd.RegisterFlagCompletionFunc("version", completion.CompleteSecretVersion)
}
