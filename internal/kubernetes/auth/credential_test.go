package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/thalassa"

	"github.com/thalassa-cloud/cli/internal/kubernetes/auth"
)

func TestWriteExecCredential(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v1/kubernetes/clusters/cluster-id/kubeconfig", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"username": "admin",
			"apiServerUrl": "https://api.example.com",
			"caCertificate": "Y2E=",
			"identity": "session-id",
			"token": "session-token",
			"kubeconfig": ""
		}`))
	}))
	t.Cleanup(server.Close)

	apiClient, err := thalassa.NewClient(
		client.WithBaseURL(server.URL),
		client.WithToken("test-token"),
	)
	require.NoError(t, err)

	var out bytes.Buffer
	err = auth.WriteExecCredential(context.Background(), apiClient, "cluster-id", &out)
	require.NoError(t, err)

	var credential map[string]any
	require.NoError(t, json.Unmarshal(out.Bytes(), &credential))
	assert.Equal(t, "client.authentication.k8s.io/v1", credential["apiVersion"])
	assert.Equal(t, "ExecCredential", credential["kind"])

	status, ok := credential["status"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "session-token", status["token"])
}
