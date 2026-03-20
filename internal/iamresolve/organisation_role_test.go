package iamresolve

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	clientiam "github.com/thalassa-cloud/client-go/iam"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

type fakeOrgRoleAPI struct {
	getFn  func(ctx context.Context, identity string) (*clientiam.OrganisationRole, error)
	listFn func(ctx context.Context, req *clientiam.ListOrganisationRolesRequest) ([]clientiam.OrganisationRole, error)
}

func (f *fakeOrgRoleAPI) GetOrganisationRole(ctx context.Context, identity string) (*clientiam.OrganisationRole, error) {
	return f.getFn(ctx, identity)
}

func (f *fakeOrgRoleAPI) ListOrganisationRoles(ctx context.Context, req *clientiam.ListOrganisationRolesRequest) ([]clientiam.OrganisationRole, error) {
	return f.listFn(ctx, req)
}

func TestResolveOrganisationRoleRef(t *testing.T) {
	ctx := context.Background()
	notFound := fmt.Errorf("role missing: %w", tcclient.ErrNotFound)

	tests := []struct {
		name    string
		api     *fakeOrgRoleAPI
		ref     string
		wantID  string
		wantErr string
	}{
		{
			name: "direct get",
			api: &fakeOrgRoleAPI{
				getFn: func(_ context.Context, id string) (*clientiam.OrganisationRole, error) {
					if id == "rid" {
						return &clientiam.OrganisationRole{Identity: "rid", Name: "Editor", Slug: "editor"}, nil
					}
					return nil, notFound
				},
				listFn: func(context.Context, *clientiam.ListOrganisationRolesRequest) ([]clientiam.OrganisationRole, error) {
					return nil, errors.New("list should not be called")
				},
			},
			ref:    "rid",
			wantID: "rid",
		},
		{
			name: "resolve by slug after not found",
			api: &fakeOrgRoleAPI{
				getFn: func(_ context.Context, _ string) (*clientiam.OrganisationRole, error) {
					return nil, notFound
				},
				listFn: func(context.Context, *clientiam.ListOrganisationRolesRequest) ([]clientiam.OrganisationRole, error) {
					return []clientiam.OrganisationRole{
						{Identity: "uuid-1", Name: "Deploy", Slug: "deploy"},
					}, nil
				},
			},
			ref:    "deploy",
			wantID: "uuid-1",
		},
		{
			name: "resolve by name case insensitive",
			api: &fakeOrgRoleAPI{
				getFn: func(_ context.Context, _ string) (*clientiam.OrganisationRole, error) {
					return nil, notFound
				},
				listFn: func(context.Context, *clientiam.ListOrganisationRolesRequest) ([]clientiam.OrganisationRole, error) {
					return []clientiam.OrganisationRole{
						{Identity: "uuid-2", Name: "Owner Role", Slug: "owner"},
					}, nil
				},
			},
			ref:    "owner role",
			wantID: "uuid-2",
		},
		{
			name: "get error not notfound",
			api: &fakeOrgRoleAPI{
				getFn: func(_ context.Context, _ string) (*clientiam.OrganisationRole, error) {
					return nil, errors.New("server down")
				},
				listFn: func(context.Context, *clientiam.ListOrganisationRolesRequest) ([]clientiam.OrganisationRole, error) {
					return nil, errors.New("list should not be called")
				},
			},
			ref:     "x",
			wantErr: "get organisation role",
		},
		{
			name: "not in list",
			api: &fakeOrgRoleAPI{
				getFn: func(_ context.Context, _ string) (*clientiam.OrganisationRole, error) {
					return nil, notFound
				},
				listFn: func(context.Context, *clientiam.ListOrganisationRolesRequest) ([]clientiam.OrganisationRole, error) {
					return []clientiam.OrganisationRole{{Identity: "a", Name: "A", Slug: "a"}}, nil
				},
			},
			ref:     "missing",
			wantErr: "role not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveOrganisationRoleRef(ctx, tt.api, tt.ref)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				assert.Nil(t, got)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tt.wantID, got.Identity)
		})
	}
}
