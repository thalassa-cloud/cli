package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/thalassa-cloud/client-go/kubernetes"
	"github.com/thalassa-cloud/client-go/thalassa"
)

const execCredentialAPIVersion = "client.authentication.k8s.io/v1"

type execCredential struct {
	APIVersion string               `json:"apiVersion"`
	Kind       string               `json:"kind"`
	Status     execCredentialStatus `json:"status"`
}

type execCredentialStatus struct {
	Token               string     `json:"token"`
	ExpirationTimestamp *time.Time `json:"expirationTimestamp,omitempty"`
}

// WriteExecCredential fetches a cluster session token and writes kubectl ExecCredential JSON.
func WriteExecCredential(ctx context.Context, client thalassa.Client, clusterIdentity string, out io.Writer) error {
	session, err := client.Kubernetes().GetKubernetesClusterKubeconfigWithParams(ctx, clusterIdentity, kubernetes.KubeconfigParams{
		SessionLifetime: 5 * time.Minute,
	})
	if err != nil {
		return err
	}
	if session.Token == "" {
		return fmt.Errorf("cluster %q returned an empty session token", clusterIdentity)
	}

	credential := execCredential{
		APIVersion: execCredentialAPIVersion,
		Kind:       "ExecCredential",
		Status: execCredentialStatus{
			Token: session.Token,
		},
	}

	encoder := json.NewEncoder(out)
	return encoder.Encode(credential)
}

// ExecutablePath returns the absolute path to the current tcloud binary, falling back to argv0.
func ExecutablePath() string {
	path, err := os.Executable()
	if err == nil && path != "" {
		return path
	}
	if len(os.Args) > 0 && os.Args[0] != "" {
		return os.Args[0]
	}
	return "tcloud"
}
