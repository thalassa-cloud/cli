package completion

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/internal/kuberesolve"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/kubernetes"
)

func completeKubernetesClusterRoleIdentities(cmd *cobra.Command) ([]string, cobra.ShellCompDirective) {
	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	roles, err := client.Kubernetes().ListKubernetesClusterRoles(cmd.Context(), &kubernetes.ListKubernetesClusterRolesRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	out := make([]string, 0, len(roles)*3)
	for _, r := range roles {
		desc := r.Name
		if r.System {
			desc += " (system)"
		}
		out = append(out, r.Identity+"\t"+desc)
		if r.Name != "" && r.Name != r.Identity {
			out = append(out, r.Name+"\t"+desc)
		}
		if r.Slug != "" && r.Slug != r.Identity && r.Slug != r.Name {
			out = append(out, r.Slug+"\t"+desc)
		}
	}
	return out, cobra.ShellCompDirectiveNoFileComp
}

// CompleteKubernetesClusterRoleIdentity completes Kubernetes cluster role identities, names, and slugs.
func CompleteKubernetesClusterRoleIdentity(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return completeKubernetesClusterRoleIdentities(cmd)
}

// CompleteKubernetesClusterRolePermissionVerb completes RBAC verbs for cluster role rules.
func CompleteKubernetesClusterRolePermissionVerb(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	verbs := make([]string, 0, len(kubernetes.KubernetesClusterRolePermissionVerbs)+1)
	verbs = append(verbs, string(kubernetes.KubernetesClusterRolePermissionVerbWildcard)+"\twildcard")
	for _, v := range kubernetes.KubernetesClusterRolePermissionVerbs {
		verbs = append(verbs, string(v))
	}
	return verbs, cobra.ShellCompDirectiveNoFileComp
}

// CompleteKubernetesClusterRoleThenBinding completes the role then binding identities for that role.
func CompleteKubernetesClusterRoleThenBinding(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	switch len(args) {
	case 0:
		return CompleteKubernetesClusterRoleIdentity(cmd, args, toComplete)
	case 1:
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		role, err := kuberesolve.ResolveKubernetesClusterRoleRef(cmd.Context(), client.Kubernetes(), args[0])
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		bindings, err := client.Kubernetes().ListClusterRoleBindings(cmd.Context(), role.Identity)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		out := make([]string, 0, len(bindings)*3)
		for _, b := range bindings {
			desc := b.Name
			if desc == "" {
				desc = b.Slug
			}
			out = append(out, b.Identity+"\t"+desc)
			if b.Slug != "" && b.Slug != b.Identity {
				out = append(out, b.Slug+"\t"+desc)
			}
			if b.Name != "" && b.Name != b.Identity && b.Name != b.Slug {
				out = append(out, b.Name+"\t"+desc)
			}
		}
		return out, cobra.ShellCompDirectiveNoFileComp
	default:
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
}

// CompleteKubernetesClusterRoleThenRule completes the role then rule identities on that role.
func CompleteKubernetesClusterRoleThenRule(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	switch len(args) {
	case 0:
		return CompleteKubernetesClusterRoleIdentity(cmd, args, toComplete)
	case 1:
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		role, err := client.Kubernetes().GetKubernetesClusterRole(cmd.Context(), args[0])
		if err != nil {
			role, err = kuberesolve.ResolveKubernetesClusterRoleRef(cmd.Context(), client.Kubernetes(), args[0])
			if err != nil {
				return nil, cobra.ShellCompDirectiveError
			}
		}
		out := make([]string, 0, len(role.Rules))
		for _, ru := range role.Rules {
			note := ru.Note
			if note == "" {
				note = strings.Join(ru.Resources, ",")
			}
			out = append(out, ru.Identity+"\t"+note)
		}
		return out, cobra.ShellCompDirectiveNoFileComp
	default:
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
}
