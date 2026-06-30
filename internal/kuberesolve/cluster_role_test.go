package kuberesolve

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thalassa-cloud/client-go/kubernetes"
)

type fakeClusterRoleAPI struct {
	getFn  func(ctx context.Context, identity string) (*kubernetes.KubernetesClusterRole, error)
	listFn func(ctx context.Context, req *kubernetes.ListKubernetesClusterRolesRequest) ([]kubernetes.KubernetesClusterRole, error)
}

func (f *fakeClusterRoleAPI) GetKubernetesClusterRole(ctx context.Context, identity string) (*kubernetes.KubernetesClusterRole, error) {
	return f.getFn(ctx, identity)
}

func (f *fakeClusterRoleAPI) ListKubernetesClusterRoles(ctx context.Context, req *kubernetes.ListKubernetesClusterRolesRequest) ([]kubernetes.KubernetesClusterRole, error) {
	return f.listFn(ctx, req)
}

func TestResolveKubernetesClusterRoleRef(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		ref     string
		api     *fakeClusterRoleAPI
		wantID  string
		wantErr string
	}{
		{
			name: "by identity from get",
			ref:  "rid",
			api: &fakeClusterRoleAPI{
				getFn: func(_ context.Context, id string) (*kubernetes.KubernetesClusterRole, error) {
					if id == "rid" {
						return &kubernetes.KubernetesClusterRole{Identity: "rid", Name: "Viewer", Slug: "viewer"}, nil
					}
					return nil, assert.AnError
				},
				listFn: func(context.Context, *kubernetes.ListKubernetesClusterRolesRequest) ([]kubernetes.KubernetesClusterRole, error) {
					return nil, nil
				},
			},
			wantID: "rid",
		},
		{
			name: "by slug from list",
			ref:  "editor",
			api: &fakeClusterRoleAPI{
				getFn: func(_ context.Context, _ string) (*kubernetes.KubernetesClusterRole, error) {
					return nil, assert.AnError
				},
				listFn: func(context.Context, *kubernetes.ListKubernetesClusterRolesRequest) ([]kubernetes.KubernetesClusterRole, error) {
					return []kubernetes.KubernetesClusterRole{
						{Identity: "eid", Name: "Editor", Slug: "editor"},
					}, nil
				},
			},
			wantID: "eid",
		},
		{
			name: "not found",
			ref:  "missing",
			api: &fakeClusterRoleAPI{
				getFn: func(_ context.Context, _ string) (*kubernetes.KubernetesClusterRole, error) {
					return nil, assert.AnError
				},
				listFn: func(context.Context, *kubernetes.ListKubernetesClusterRolesRequest) ([]kubernetes.KubernetesClusterRole, error) {
					return []kubernetes.KubernetesClusterRole{{Identity: "a", Name: "A", Slug: "a"}}, nil
				},
			},
			wantErr: "kubernetes cluster role not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ResolveKubernetesClusterRoleRef(context.Background(), tt.api, tt.ref)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantID, got.Identity)
		})
	}
}
