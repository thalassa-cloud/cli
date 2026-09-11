package kuberesolve

import (
	"fmt"
	"os"
	"strings"

	"github.com/thalassa-cloud/cli/internal/dirconfig"
)

const (
	EnvClusterIdentity = "TCLOUD_CLUSTER_IDENTITY"
	EnvClusterSlug     = "TCLOUD_CLUSTER_SLUG"
	EnvClusterName     = "TCLOUD_CLUSTER_NAME"
)

// PreferredClusterRef resolves a cluster reference from explicit input, connect-shell env, or directory config.
func PreferredClusterRef(explicit string) (string, bool) {
	if ref := strings.TrimSpace(explicit); ref != "" {
		return ref, true
	}
	for _, env := range []string{EnvClusterIdentity, EnvClusterSlug, EnvClusterName} {
		if v := strings.TrimSpace(os.Getenv(env)); v != "" {
			return v, true
		}
	}
	if cluster := strings.TrimSpace(dirconfig.Current().Config.Kubernetes.Cluster); cluster != "" {
		return cluster, true
	}
	return "", false
}

// RequireClusterRef is PreferredClusterRef with an error when no cluster is configured.
func RequireClusterRef(explicit string) (string, error) {
	ref, ok := PreferredClusterRef(explicit)
	if !ok {
		return "", fmt.Errorf("cluster is required (argument, --cluster, or kubernetes.cluster in .thalassa)")
	}
	return ref, nil
}
