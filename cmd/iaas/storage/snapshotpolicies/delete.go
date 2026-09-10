package snapshotpolicies

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/labels"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/iaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	deleteForce         bool
	deleteLabelSelector string
)

var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete snapshot policy(ies)",
	Aliases: []string{"d", "del", "remove", "rm"},
	Args:    cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && deleteLabelSelector == "" {
			return fmt.Errorf("either snapshot policy identity(ies) or --selector must be provided")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		toDelete := []iaas.SnapshotPolicy{}
		if deleteLabelSelector != "" {
			all, err := client.IaaS().ListSnapshotPolicies(cmd.Context(), &iaas.ListSnapshotPoliciesRequest{
				Filters: []filters.Filter{
					&filters.LabelFilter{MatchLabels: labels.ParseLabelSelector(deleteLabelSelector)},
				},
			})
			if err != nil {
				return fmt.Errorf("failed to list snapshot policies: %w", err)
			}
			if len(all) == 0 {
				fmt.Println("No snapshot policies found matching the label selector")
				return nil
			}
			toDelete = append(toDelete, all...)
		} else {
			for _, identity := range args {
				policy, err := client.IaaS().GetSnapshotPolicy(cmd.Context(), identity)
				if err != nil {
					if tcclient.IsNotFound(err) {
						fmt.Printf("Snapshot policy %s not found\n", identity)
						continue
					}
					return fmt.Errorf("failed to get snapshot policy: %w", err)
				}
				toDelete = append(toDelete, *policy)
			}
		}

		if len(toDelete) == 0 {
			fmt.Println("No snapshot policies to delete")
			return nil
		}

		var summary strings.Builder
		fmt.Fprintf(&summary, "Are you sure you want to delete the following snapshot policy(ies)?\n")
		for _, policy := range toDelete {
			fmt.Fprintf(&summary, "  %s (%s)\n", policy.Name, policy.Identity)
		}
		proceed, err := shared.PromptDestructiveUnlessForce(deleteForce, summary.String())
		if err != nil {
			return err
		}
		if !proceed {
			return nil
		}

		for _, policy := range toDelete {
			fmt.Printf("Deleting snapshot policy: %s (%s)\n", policy.Name, policy.Identity)
			if err := client.IaaS().DeleteSnapshotPolicy(cmd.Context(), policy.Identity); err != nil {
				return fmt.Errorf("failed to delete snapshot policy: %w", err)
			}
			fmt.Printf("Snapshot policy %s deleted successfully\n", policy.Identity)
		}
		return nil
	},
}

func init() {
	SnapshotPoliciesCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVar(&deleteForce, "force", false, "Force the deletion and skip the confirmation")
	deleteCmd.Flags().StringVarP(&deleteLabelSelector, "selector", "l", "", "Label selector to filter snapshot policies")
	deleteCmd.ValidArgsFunction = completion.CompleteSnapshotPolicyID
}
