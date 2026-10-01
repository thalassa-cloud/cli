package policies

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/labels"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/filters"
	clientiam "github.com/thalassa-cloud/client-go/iam"
)

var policiesListLabelSelector string

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List IAM policies",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		req := &clientiam.ListIamPoliciesRequest{}
		if policiesListLabelSelector != "" {
			req.Filters = []filters.Filter{
				&filters.LabelFilter{MatchLabels: labels.ParseLabelSelector(policiesListLabelSelector)},
			}
		}
		policies, err := client.IAM().ListIamPolicies(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to list IAM policies: %w", err)
		}
		body := make([][]string, 0, len(policies))
		for _, p := range policies {
			sys := "no"
			if p.System {
				sys = "yes"
			}
			ro := "no"
			if p.IsReadOnly {
				ro = "yes"
			}
			replicate := "no"
			if p.ReplicateToChildren {
				replicate = "yes"
			}
			body = append(body, []string{
				p.Identity,
				p.Name,
				p.Slug,
				sys,
				ro,
				replicate,
				fmt.Sprintf("%d", len(p.Rules)),
				fmt.Sprintf("%d", len(p.Bindings)),
			})
		}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Slug", "System", "Read-only", "Replicate", "Rules", "Bindings"}, body)
		}
		return nil
	},
}

func init() {
	PoliciesCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	listCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
	listCmd.Flags().StringVar(&policiesListLabelSelector, "label-selector", "", "Filter by labels (key=value,key2=value2)")
}
