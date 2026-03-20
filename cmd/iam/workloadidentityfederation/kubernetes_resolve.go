package workloadidentityfederation

import (
	"context"
	"fmt"
	"strings"

	"github.com/thalassa-cloud/client-go/kubernetes"
)

func resolveKubernetesCluster(ctx context.Context, kc *kubernetes.Client, ref string) (*kubernetes.KubernetesCluster, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return nil, fmt.Errorf("kubernetes cluster reference is empty")
	}
	c, err := kc.GetKubernetesCluster(ctx, ref)
	if err == nil && c != nil {
		return c, nil
	}
	list, err := kc.ListKubernetesClusters(ctx, &kubernetes.ListKubernetesClustersRequest{})
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
