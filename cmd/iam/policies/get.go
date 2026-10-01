package policies

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/iamresolve"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var getCmd = &cobra.Command{
	Use:               "get <policy>",
	Short:             "Show a policy including rules and bindings summary",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteIAMPolicyIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		policy, err := iamresolve.ResolveIamPolicyRef(ctx, client.IAM(), args[0])
		if err != nil {
			return err
		}
		// Re-fetch by identity so nested rules/bindings are included when list resolution was used.
		policy, err = client.IAM().GetIamPolicy(ctx, policy.Identity)
		if err != nil {
			return fmt.Errorf("failed to get IAM policy: %w", err)
		}

		fmt.Printf("Identity:             %s\n", policy.Identity)
		fmt.Printf("Name:                 %s\n", policy.Name)
		fmt.Printf("Slug:                 %s\n", policy.Slug)
		fmt.Printf("Description:          %s\n", policy.Description)
		fmt.Printf("System:               %v\n", policy.System)
		fmt.Printf("Read-only:            %v\n", policy.IsReadOnly)
		fmt.Printf("Replicate to children: %v\n", policy.ReplicateToChildren)
		fmt.Printf("Created:              %s\n", formattime.FormatTime(policy.CreatedAt.Local(), showExactTime))
		if policy.Project != nil {
			fmt.Printf("Project:              %s (%s)\n", policy.Project.Name, policy.Project.Identity)
		}

		if len(policy.Rules) > 0 {
			fmt.Println("\nRules:")
			body := make([][]string, 0, len(policy.Rules))
			for _, ru := range policy.Rules {
				perms := make([]string, 0, len(ru.Permissions))
				for _, p := range ru.Permissions {
					perms = append(perms, string(p))
				}
				body = append(body, []string{
					ru.Identity,
					strings.Join(ru.Resources, ","),
					strings.Join(ru.ResourceIdentities, ","),
					strings.Join(perms, ","),
					ru.Note,
				})
			}
			table.Print([]string{"ID", "Resources", "Resource IDs", "Permissions", "Note"}, body)
		}
		if len(policy.Bindings) > 0 {
			fmt.Println("\nBindings:")
			body := make([][]string, 0, len(policy.Bindings))
			for _, b := range policy.Bindings {
				subject := ""
				switch {
				case b.User != nil:
					subject = "user:" + shared.UserPtrDisplay(b.User)
				case b.ServiceAccount != nil:
					subject = "service_account:" + b.ServiceAccount.Slug
				}
				expires := ""
				if b.ExpiresAt != nil {
					expires = formattime.FormatTime(b.ExpiresAt.Local(), showExactTime)
				}
				body = append(body, []string{b.Identity, b.Name, subject, expires})
			}
			table.Print([]string{"ID", "Name", "Subject", "Expires"}, body)
		}
		return nil
	},
}

func init() {
	PoliciesCmd.AddCommand(getCmd)
	getCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	getCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
}
