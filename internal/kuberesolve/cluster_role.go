package kuberesolve

import (
	"context"
	"fmt"
	"strings"

	"github.com/thalassa-cloud/client-go/kubernetes"
)

// KubernetesClusterRoleAPI lists and fetches Kubernetes cluster IAM roles.
type KubernetesClusterRoleAPI interface {
	GetKubernetesClusterRole(ctx context.Context, identity string) (*kubernetes.KubernetesClusterRole, error)
	ListKubernetesClusterRoles(ctx context.Context, request *kubernetes.ListKubernetesClusterRolesRequest) ([]kubernetes.KubernetesClusterRole, error)
}

// ResolveKubernetesClusterRoleRef resolves a cluster role by identity, name, or slug.
func ResolveKubernetesClusterRoleRef(ctx context.Context, api KubernetesClusterRoleAPI, ref string) (*kubernetes.KubernetesClusterRole, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return nil, fmt.Errorf("kubernetes cluster role reference is empty")
	}
	role, err := api.GetKubernetesClusterRole(ctx, ref)
	if err == nil && role != nil {
		return role, nil
	}
	list, err := api.ListKubernetesClusterRoles(ctx, &kubernetes.ListKubernetesClusterRolesRequest{})
	if err != nil {
		return nil, fmt.Errorf("list kubernetes cluster roles: %w", err)
	}
	for i := range list {
		r := &list[i]
		if strings.EqualFold(r.Identity, ref) || strings.EqualFold(r.Slug, ref) || strings.EqualFold(r.Name, ref) {
			return r, nil
		}
	}
	return nil, fmt.Errorf("kubernetes cluster role not found: %s", ref)
}
