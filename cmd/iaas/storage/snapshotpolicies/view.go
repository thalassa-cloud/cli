package snapshotpolicies

import (
	"fmt"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
)

var viewOutputFormat string

var viewCmd = &cobra.Command{
	Use:               "view",
	Short:             "View snapshot policy details",
	Aliases:           []string{"show", "get", "describe"},
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteSnapshotPolicyID,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		policy, err := client.IaaS().GetSnapshotPolicy(cmd.Context(), args[0])
		if err != nil {
			return fmt.Errorf("failed to get snapshot policy: %w", err)
		}

		if viewOutputFormat == "yaml" {
			return outputYAML(*policy)
		}

		fmt.Printf("Snapshot Policy Details:\n")
		fmt.Printf("  ID: %s\n", policy.Identity)
		fmt.Printf("  Name: %s\n", policy.Name)
		fmt.Printf("  Description: %s\n", policy.Description)
		fmt.Printf("  Region: %s\n", regionName(*policy))
		fmt.Printf("  Enabled: %t\n", policy.Enabled)
		fmt.Printf("  Schedule: %s\n", policy.Schedule)
		fmt.Printf("  Timezone: %s\n", policy.Timezone)
		fmt.Printf("  TTL: %s\n", policy.Ttl.String())
		fmt.Printf("  Keep Count: %s\n", keepCountString(policy.KeepCount))
		fmt.Printf("  Target: %s\n", targetSummary(policy.Target))
		fmt.Printf("  Created: %s\n", formattime.FormatTime(policy.CreatedAt.Local(), false))
		fmt.Printf("  Updated: %s\n", formattime.FormatTime(policy.UpdatedAt.Local(), false))
		if policy.NextSnapshotAt != nil {
			fmt.Printf("  Next Snapshot: %s\n", formattime.FormatTime(policy.NextSnapshotAt.Local(), true))
		}
		if policy.LastSnapshotAt != nil {
			fmt.Printf("  Last Snapshot: %s\n", formattime.FormatTime(policy.LastSnapshotAt.Local(), true))
		}
		return nil
	},
}

func outputYAML(policy iaas.SnapshotPolicy) error {
	policy.Organisation = nil
	yamlData, err := yaml.Marshal(&policy)
	if err != nil {
		return fmt.Errorf("failed to marshal to YAML: %w", err)
	}
	fmt.Print(string(yamlData))
	return nil
}

func init() {
	SnapshotPoliciesCmd.AddCommand(viewCmd)
	viewCmd.Flags().StringVarP(&viewOutputFormat, "output", "o", "", "Output format (yaml)")
	_ = viewCmd.RegisterFlagCompletionFunc("output", completion.CompleteOutputFormat)
}
