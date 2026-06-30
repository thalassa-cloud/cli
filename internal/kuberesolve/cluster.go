package kuberesolve

import (
	"context"
	"fmt"
	"strings"

	"github.com/thalassa-cloud/client-go/kubernetes"
)

// KubernetesClusterAPI lists and fetches Kubernetes clusters.
type KubernetesClusterAPI interface {
	GetKubernetesCluster(ctx context.Context, identity string) (*kubernetes.KubernetesCluster, error)
	ListKubernetesClusters(ctx context.Context, request *kubernetes.ListKubernetesClustersRequest) ([]kubernetes.KubernetesCluster, error)
}

// ResolveKubernetesClusterRef resolves a cluster by identity, name, or slug.
func ResolveKubernetesClusterRef(ctx context.Context, api KubernetesClusterAPI, ref string) (*kubernetes.KubernetesCluster, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return nil, fmt.Errorf("kubernetes cluster reference is empty")
	}
	c, err := api.GetKubernetesCluster(ctx, ref)
	if err == nil && c != nil {
		return c, nil
	}
	list, err := api.ListKubernetesClusters(ctx, &kubernetes.ListKubernetesClustersRequest{})
	if err != nil {
		return nil, fmt.Errorf("list kubernetes clusters: %w", err)
	}
	for i := range list {
		cl := &list[i]
		if strings.EqualFold(cl.Identity, ref) || strings.EqualFold(cl.Slug, ref) || strings.EqualFold(cl.Name, ref) {
			return cl, nil
		}
	}
	return nil, fmt.Errorf("kubernetes cluster not found: %s", ref)
}
