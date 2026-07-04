package projectresolve

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thalassa-cloud/client-go/projects"
)

type mockProjectsAPI struct {
	getProject func(ctx context.Context, identity string) (*projects.Project, error)
}

func (m *mockProjectsAPI) GetProject(ctx context.Context, identity string) (*projects.Project, error) {
	return m.getProject(ctx, identity)
}

func TestResolveProjectIdentity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		ref     string
		api     *mockProjectsAPI
		want    string
		wantErr string
	}{
		{
			name: "root reference",
			ref:  "root",
			api:  &mockProjectsAPI{},
			want: "",
		},
		{
			name: "empty ref",
			ref:  "",
			api:  &mockProjectsAPI{},
			want: "",
		},
		{
			name: "resolves slug to identity",
			ref:  "production",
			api: &mockProjectsAPI{
				getProject: func(_ context.Context, ref string) (*projects.Project, error) {
					assert.Equal(t, "production", ref)
					return &projects.Project{Identity: "prj-abc123", Slug: "production"}, nil
				},
			},
			want: "prj-abc123",
		},
		{
			name: "identity passthrough via API",
			ref:  "prj-abc123",
			api: &mockProjectsAPI{
				getProject: func(_ context.Context, ref string) (*projects.Project, error) {
					assert.Equal(t, "prj-abc123", ref)
					return &projects.Project{Identity: "prj-abc123", Slug: "production"}, nil
				},
			},
			want: "prj-abc123",
		},
		{
			name: "get project error",
			ref:  "missing",
			api: &mockProjectsAPI{
				getProject: func(context.Context, string) (*projects.Project, error) {
					return nil, errors.New("not found")
				},
			},
			wantErr: `project not found: "missing"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ResolveProjectIdentity(context.Background(), tt.api, tt.ref)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
