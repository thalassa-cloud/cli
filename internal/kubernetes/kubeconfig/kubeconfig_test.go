package kubeconfig_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thalassa-cloud/client-go/kubernetes"

	"github.com/thalassa-cloud/cli/internal/kubernetes/kubeconfig"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

func TestRenderExecKubeconfigOmitsToken(t *testing.T) {
	config, err := kubeconfig.Render(kubeconfig.RenderInput{
		Cluster: &kubernetes.KubernetesCluster{
			Identity:     "cluster-id",
			Slug:         "prod",
			APIServerURL: "https://api.example.com",
		},
		Session: &kubernetes.KubernetesClusterSessionToken{
			Username:      "admin",
			CACertificate: "Y2E=",
			Token:         "secret-token",
		},
		Executable:      "/usr/local/bin/tcloud",
		ClusterIdentity: "cluster-id",
		Scope: thalassaclient.Scope{
			Context:      "dev",
			Organisation: "acme",
			Project:      "my-project",
		},
	}, kubeconfig.AuthModeExec)
	require.NoError(t, err)

	assert.Contains(t, config, "exec:")
	assert.Contains(t, config, `command: "/usr/local/bin/tcloud"`)
	assert.Contains(t, config, "credential")
	assert.Contains(t, config, "cluster-id")
	assert.Contains(t, config, "--organisation")
	assert.Contains(t, config, "acme")
	assert.Contains(t, config, "--context")
	assert.Contains(t, config, "dev")
	assert.Contains(t, config, "--project")
	assert.Contains(t, config, "my-project")
	assert.NotContains(t, config, "secret-token")
}

func TestRenderExecKubeconfigOmitsProjectWhenUnset(t *testing.T) {
	config, err := kubeconfig.Render(kubeconfig.RenderInput{
		Cluster: &kubernetes.KubernetesCluster{
			Identity:     "cluster-id",
			Slug:         "prod",
			APIServerURL: "https://api.example.com",
		},
		Session: &kubernetes.KubernetesClusterSessionToken{
			Username:      "admin",
			CACertificate: "Y2E=",
		},
		Executable:      "/usr/local/bin/tcloud",
		ClusterIdentity: "cluster-id",
		Scope: thalassaclient.Scope{
			Organisation: "acme",
		},
	}, kubeconfig.AuthModeExec)
	require.NoError(t, err)

	assert.Contains(t, config, "--organisation")
	assert.Contains(t, config, "acme")
	assert.NotContains(t, config, "--project")
}

func TestRenderExecKubeconfigRequiresOrganisation(t *testing.T) {
	_, err := kubeconfig.Render(kubeconfig.RenderInput{
		Cluster: &kubernetes.KubernetesCluster{
			Slug:         "prod",
			APIServerURL: "https://api.example.com",
		},
		Session:         &kubernetes.KubernetesClusterSessionToken{Username: "admin", CACertificate: "Y2E="},
		ClusterIdentity: "cluster-id",
	}, kubeconfig.AuthModeExec)
	require.Error(t, err)
}

func TestRenderInlineTokenKubeconfigIncludesToken(t *testing.T) {
	config, err := kubeconfig.Render(kubeconfig.RenderInput{
		Cluster: &kubernetes.KubernetesCluster{
			Slug:         "prod",
			APIServerURL: "https://api.example.com",
		},
		Session: &kubernetes.KubernetesClusterSessionToken{
			Username:      "admin",
			CACertificate: "Y2E=",
			Token:         "secret-token",
		},
		ClusterIdentity: "cluster-id",
	}, kubeconfig.AuthModeToken)
	require.NoError(t, err)

	assert.Contains(t, config, "token: secret-token")
	assert.NotContains(t, strings.ToLower(config), "exec:")
}
