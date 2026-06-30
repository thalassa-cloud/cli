package rules

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/kubernetes/iam/shared"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/kuberesolve"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/kubernetes"
)

var (
	noHeader        bool
	ruleResources   []string
	ruleVerbs       []string
	ruleAPIGroups   []string
	ruleResNames    []string
	ruleNonResURLs  []string
	ruleNote        string
)

var addCmd = &cobra.Command{
	Use:               "add <role>",
	Short:             "Add a permission rule to a Kubernetes cluster role",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteKubernetesClusterRoleIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if len(ruleVerbs) == 0 {
			return fmt.Errorf("at least one --verb is required (get, list, watch, create, update, delete, patch, *)")
		}
		if len(ruleResources) == 0 && len(ruleNonResURLs) == 0 {
			return fmt.Errorf("specify at least one --resource or --non-resource-url")
		}

		verbs := make([]kubernetes.KubernetesClusterRolePermissionVerb, 0, len(ruleVerbs))
		for _, v := range ruleVerbs {
			verbs = append(verbs, kubernetes.KubernetesClusterRolePermissionVerb(v))
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		role, err := kuberesolve.ResolveKubernetesClusterRoleRef(ctx, client.Kubernetes(), args[0])
		if err != nil {
			return err
		}

		rule, err := client.Kubernetes().AddClusterRoleRule(ctx, role.Identity, kubernetes.AddKubernetesClusterRolePermissionRule{
			Resources:       ruleResources,
			Verbs:           verbs,
			ApiGroups:       ruleAPIGroups,
			ResourceNames:   ruleResNames,
			NonResourceURLs: ruleNonResURLs,
			Note:            ruleNote,
		})
		if err != nil {
			return fmt.Errorf("failed to add rule: %w", err)
		}
		if rule == nil {
			return nil
		}
		verbStrs := make([]string, 0, len(rule.Verbs))
		for _, v := range rule.Verbs {
			verbStrs = append(verbStrs, string(v))
		}
		body := [][]string{{rule.Identity, strings.Join(rule.Resources, ","), strings.Join(verbStrs, ",")}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"Rule ID", "Resources", "Verbs"}, body)
		}
		return nil
	},
}

func init() {
	RulesCmd.AddCommand(addCmd)
	addCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	addCmd.Flags().StringSliceVar(&ruleResources, "resource", nil, "Resource name (repeatable)")
	addCmd.Flags().StringSliceVar(&ruleVerbs, "verb", nil, "Verb: get, list, watch, create, update, delete, patch, or * (repeatable)")
	addCmd.Flags().StringSliceVar(&ruleAPIGroups, "api-group", nil, "API group (repeatable)")
	addCmd.Flags().StringSliceVar(&ruleResNames, "resource-name", nil, "Concrete resource name (repeatable)")
	addCmd.Flags().StringSliceVar(&ruleNonResURLs, "non-resource-url", nil, "Non-resource URL (repeatable)")
	addCmd.Flags().StringVar(&ruleNote, "note", "", "Human-readable note for the rule")
	_ = addCmd.RegisterFlagCompletionFunc("verb", completion.CompleteKubernetesClusterRolePermissionVerb)
}
