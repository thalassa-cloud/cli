package kubeconfig

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/thalassa-cloud/cli/internal/kubernetes/auth"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/kubernetes"
)

//go:embed templates/kubeconfig-exec.yaml
var execTemplate string

//go:embed templates/kubeconfig-token.yaml
var tokenTemplate string

type AuthMode string

const (
	AuthModeExec  AuthMode = "exec"
	AuthModeToken AuthMode = "token"
)

type RenderInput struct {
	Cluster         *kubernetes.KubernetesCluster
	Session         *kubernetes.KubernetesClusterSessionToken
	Executable      string
	ClusterIdentity string
	Scope           thalassaclient.Scope
}

// Render builds a kubeconfig for the given cluster session.
func Render(input RenderInput, mode AuthMode) (string, error) {
	if input.Cluster == nil {
		return "", fmt.Errorf("cluster is required")
	}
	if input.Session == nil {
		return "", fmt.Errorf("session is required")
	}
	if mode == AuthModeExec && strings.TrimSpace(input.Scope.Organisation) == "" {
		return "", fmt.Errorf("organisation is required for exec kubeconfig")
	}

	templateInput := struct {
		APIServerURL        string
		Base64CACertificate string
		SessionToken        string
		Cluster             string
		User                string
		Executable          string
		ExecArgs            []string
	}{
		APIServerURL:        input.Cluster.APIServerURL,
		Base64CACertificate: input.Session.CACertificate,
		SessionToken:        input.Session.Token,
		Cluster:             input.Cluster.Slug,
		User:                input.Session.Username,
		Executable:          input.Executable,
		ExecArgs:            execCredentialArgs(input.ClusterIdentity, input.Scope),
	}

	source := execTemplate
	if mode == AuthModeToken {
		source = tokenTemplate
	}

	tmpl, err := template.New("kubeconfig").Parse(source)
	if err != nil {
		return "", fmt.Errorf("parse kubeconfig template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateInput); err != nil {
		return "", fmt.Errorf("render kubeconfig template: %w", err)
	}
	return strings.TrimSpace(buf.String()) + "\n", nil
}

// NewRenderInput builds render input using the current CLI binary path and context scope.
func NewRenderInput(cluster *kubernetes.KubernetesCluster, session *kubernetes.KubernetesClusterSessionToken) RenderInput {
	return RenderInput{
		Cluster:         cluster,
		Session:         session,
		Executable:      auth.ExecutablePath(),
		ClusterIdentity: cluster.Identity,
		Scope:           thalassaclient.ScopeFromContext(),
	}
}

func execCredentialArgs(clusterIdentity string, scope thalassaclient.Scope) []string {
	args := []string{
		"kubernetes",
		"credential",
		"--cluster", clusterIdentity,
		"--organisation", scope.Organisation,
	}
	if scope.Project != "" {
		args = append(args, "--project", scope.Project)
	}
	return args
}
