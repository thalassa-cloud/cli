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

var (
	rotationEnabled    bool
	rotationPeriodDays int
)

var rotationCmd = &cobra.Command{
	Use:               "rotation <key>",
	Short:             "Update automatic rotation settings for a KMS key",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteKmsKeyIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		req := clientkms.UpdateRotationRequest{}
		if cmd.Flags().Changed("enabled") {
			enabled := rotationEnabled
			req.KeyRotationEnabled = &enabled
		}
		if cmd.Flags().Changed("period-days") {
			period := rotationPeriodDays
			req.RotationPeriodInDays = &period
		}
		if req.KeyRotationEnabled == nil && req.RotationPeriodInDays == nil {
			return fmt.Errorf("at least one of --enabled or --period-days is required")
		}

		key, err := client.KMS().UpdateRotation(cmd.Context(), region, args[0], req)
		if err != nil {
			return fmt.Errorf("failed to update KMS key rotation: %w", err)
		}

		period := "-"
		if key.RotationPeriodInDays != nil {
			period = fmt.Sprintf("%d", *key.RotationPeriodInDays)
		}
		body := [][]string{{
			key.Identity,
			key.Name,
			fmt.Sprintf("%v", key.KeyRotationEnabled),
			period,
			formattime.FormatTime(key.UpdatedAt.Local(), showExactTime),
		}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Rotation Enabled", "Period Days", "Updated"}, body)
		}
		return nil
	},
}

func init() {
	KeysCmd.AddCommand(rotationCmd)
	rotationCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	rotationCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
	rotationCmd.Flags().StringVar(&region, "region", "", "Region")
	rotationCmd.Flags().BoolVar(&rotationEnabled, "enabled", false, "Enable or disable automatic rotation")
	rotationCmd.Flags().IntVar(&rotationPeriodDays, "period-days", 0, "Automatic rotation period in days")
	_ = rotationCmd.MarkFlagRequired("region")
	_ = rotationCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
}
