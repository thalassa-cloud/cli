package roles

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/kubernetes/iam/shared"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/kuberesolve"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/kubernetes"
)

var (
	showExactTime bool
)

var getCmd = &cobra.Command{
	Use:               "get <role>",
	Short:             "Show a Kubernetes cluster role including rules and bindings",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteKubernetesClusterRoleIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		role, err := kuberesolve.ResolveKubernetesClusterRoleRef(ctx, client.Kubernetes(), args[0])
		if err != nil {
			return err
		}
		role, err = client.Kubernetes().GetKubernetesClusterRole(ctx, role.Identity)
		if err != nil {
			return fmt.Errorf("failed to get role: %w", err)
		}

		fmt.Printf("Identity:    %s\n", role.Identity)
		fmt.Printf("Name:        %s\n", role.Name)
		fmt.Printf("Slug:        %s\n", role.Slug)
		fmt.Printf("Description: %s\n", role.Description)
		fmt.Printf("System:      %v\n", role.System)
		fmt.Printf("Created:     %s\n", formattime.FormatTime(role.CreatedAt.Local(), showExactTime))

		if len(role.Rules) > 0 {
			fmt.Println("\nRules:")
			body := make([][]string, 0, len(role.Rules))
			for _, ru := range role.Rules {
				verbs := make([]string, 0, len(ru.Verbs))
				for _, v := range ru.Verbs {
					verbs = append(verbs, string(v))
				}
				body = append(body, []string{
					ru.Identity,
					strings.Join(ru.ApiGroups, ","),
					strings.Join(ru.Resources, ","),
					strings.Join(verbs, ","),
					ru.Note,
				})
			}
			table.Print([]string{"ID", "API groups", "Resources", "Verbs", "Note"}, body)
		}
		if len(role.Bindings) > 0 {
			fmt.Println("\nBindings:")
			body := make([][]string, 0, len(role.Bindings))
			for _, b := range role.Bindings {
				body = append(body, []string{b.Identity, b.Name, bindingSubject(b)})
			}
			table.Print([]string{"ID", "Name", "Subject"}, body)
		}
		return nil
	},
}

func bindingSubject(b kubernetes.KubernetesClusterRoleBinding) string {
	switch {
	case b.User != nil:
		return "user:" + shared.UserPtrDisplay(b.User)
	case b.OrganisationTeam != nil:
		return "team:" + b.OrganisationTeam.Slug
	case b.ServiceAccount != nil:
		return "service_account:" + b.ServiceAccount.Slug
	default:
		return ""
	}
}

func init() {
	RolesCmd.AddCommand(getCmd)
	getCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	getCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
}
