package keys

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientkms "github.com/thalassa-cloud/client-go/kms"
)

func runKeyStatusMutation(ctx context.Context, keyIdentity string, mutate func(context.Context, string, string) (*clientkms.KmsKey, error), action string) error {
	key, err := mutate(ctx, region, keyIdentity)
	if err != nil {
		return fmt.Errorf("failed to %s KMS key: %w", action, err)
	}

	body := [][]string{{
		key.Identity,
		key.Name,
		string(key.Status),
		formattime.FormatTime(key.UpdatedAt.Local(), showExactTime),
	}}
	if noHeader {
		table.Print(nil, body)
	} else {
		table.Print([]string{"ID", "Name", "Status", "Updated"}, body)
	}
	return nil
}

func registerKeyRegionFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	cmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
	cmd.Flags().StringVar(&region, "region", "", "Region")
	_ = cmd.MarkFlagRequired("region")
	_ = cmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
}

var enableCmd = &cobra.Command{
	Use:               "enable <key>",
	Short:             "Enable a KMS key",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteKmsKeyIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		return runKeyStatusMutation(cmd.Context(), args[0], client.KMS().EnableKey, "enable")
	},
}

var disableCmd = &cobra.Command{
	Use:               "disable <key>",
	Short:             "Disable a KMS key",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteKmsKeyIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		return runKeyStatusMutation(cmd.Context(), args[0], client.KMS().DisableKey, "disable")
	},
}

func init() {
	KeysCmd.AddCommand(enableCmd)
	registerKeyRegionFlags(enableCmd)

	KeysCmd.AddCommand(disableCmd)
	registerKeyRegionFlags(disableCmd)
}
