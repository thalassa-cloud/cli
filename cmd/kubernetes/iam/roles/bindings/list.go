package bindings

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/kubernetes/iam/shared"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/kuberesolve"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/kubernetes"
)

var noHeader bool

var listCmd = &cobra.Command{
	Use:               "list <role>",
	Short:             "List bindings for a Kubernetes cluster role",
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

		bindings, err := client.Kubernetes().ListClusterRoleBindings(ctx, role.Identity)
		if err != nil {
			return fmt.Errorf("failed to list bindings: %w", err)
		}
		body := make([][]string, 0, len(bindings))
		for _, b := range bindings {
			body = append(body, []string{b.Identity, b.Name, bindingSubject(b)})
		}
		if len(body) == 0 {
			fmt.Println("No bindings found")
			return nil
		}
		if noHeader {
			table.Print(nil, body)
		} else {
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
	BindingsCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
}
