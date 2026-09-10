package zones

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientdns "github.com/thalassa-cloud/client-go/dns"
)

var dnssecCmd = &cobra.Command{
	Use:   "dnssec",
	Short: "Manage DNSSEC for a DNS zone",
}

var (
	dnssecRegion         string
	dnssecKmsKeyIdentity string
	dnssecDisableForce   bool
)

var dnssecGetCmd = &cobra.Command{
	Use:               "get <zone>",
	Aliases:           []string{"view", "show", "status"},
	Short:             "Get DNSSEC status for a zone",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteDnsZoneIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		status, err := client.DNS().GetDnssec(cmd.Context(), args[0])
		if err != nil {
			return fmt.Errorf("failed to get DNSSEC status: %w", err)
		}

		dsCount := fmt.Sprintf("%d", len(status.DsRecords))
		body := [][]string{{
			fmt.Sprintf("%v", status.Enabled),
			fmt.Sprintf("%v", status.DsDelegated),
			dsCount,
			status.Region,
			status.KmsKeyIdentity,
		}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"Enabled", "DS Delegated", "DS Records", "Region", "KMS Key"}, body)
		}

		if len(status.DsRecords) > 0 {
			fmt.Println()
			dsBody := make([][]string, 0, len(status.DsRecords))
			for _, ds := range status.DsRecords {
				dsBody = append(dsBody, []string{
					fmt.Sprintf("%d", ds.KeyTag),
					fmt.Sprintf("%d", ds.Algorithm),
					ds.DigestTypeName,
					ds.KeyRole,
					ds.Record,
				})
			}
			if noHeader {
				table.Print(nil, dsBody)
			} else {
				table.Print([]string{"Key Tag", "Algorithm", "Digest Type", "Key Role", "Record"}, dsBody)
			}
		}
		return nil
	},
}

var dnssecEnableCmd = &cobra.Command{
	Use:               "enable <zone>",
	Short:             "Enable DNSSEC signing for a zone",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteDnsZoneIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		status, err := client.DNS().SetDnssec(cmd.Context(), args[0], clientdns.SetDnssecRequest{
			Region:         dnssecRegion,
			KmsKeyIdentity: dnssecKmsKeyIdentity,
		})
		if err != nil {
			return fmt.Errorf("failed to enable DNSSEC: %w", err)
		}

		body := [][]string{{
			fmt.Sprintf("%v", status.Enabled),
			status.Region,
			status.KmsKeyIdentity,
		}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"Enabled", "Region", "KMS Key"}, body)
		}
		return nil
	},
}

var dnssecDisableCmd = &cobra.Command{
	Use:               "disable <zone>",
	Short:             "Disable DNSSEC signing for a zone",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteDnsZoneIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		ok, err := shared.PromptDestructiveUnlessForce(dnssecDisableForce, fmt.Sprintf(
			"Are you sure you want to disable DNSSEC for zone %s?\n", args[0],
		))
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}

		if err := client.DNS().DeleteDnssec(cmd.Context(), args[0]); err != nil {
			return fmt.Errorf("failed to disable DNSSEC: %w", err)
		}
		fmt.Printf("Disabled DNSSEC for zone %s\n", args[0])
		return nil
	},
}

func init() {
	ZonesCmd.AddCommand(dnssecCmd)
	dnssecCmd.AddCommand(dnssecGetCmd)
	dnssecCmd.AddCommand(dnssecEnableCmd)
	dnssecCmd.AddCommand(dnssecDisableCmd)

	dnssecGetCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")

	dnssecEnableCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	dnssecEnableCmd.Flags().StringVar(&dnssecRegion, "region", "", "Region for DNSSEC KMS key")
	dnssecEnableCmd.Flags().StringVar(&dnssecKmsKeyIdentity, "kms-key", "", "KMS key identity used for DNSSEC signing")
	_ = dnssecEnableCmd.MarkFlagRequired("region")
	_ = dnssecEnableCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
	_ = dnssecEnableCmd.RegisterFlagCompletionFunc("kms-key", completion.CompleteKmsKeyIdentity)

	dnssecDisableCmd.Flags().BoolVar(&dnssecDisableForce, shared.ForceKey, false, "Skip the confirmation prompt")
}
