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

type fakeIamPolicyAPI struct {
	getFn  func(ctx context.Context, identity string) (*clientiam.IamPolicy, error)
	listFn func(ctx context.Context, req *clientiam.ListIamPoliciesRequest) ([]clientiam.IamPolicy, error)
}

func (f *fakeIamPolicyAPI) GetIamPolicy(ctx context.Context, identity string) (*clientiam.IamPolicy, error) {
	return f.getFn(ctx, identity)
}

func (f *fakeIamPolicyAPI) ListIamPolicies(ctx context.Context, req *clientiam.ListIamPoliciesRequest) ([]clientiam.IamPolicy, error) {
	return f.listFn(ctx, req)
}

func TestResolveIamPolicyRef(t *testing.T) {
	ctx := context.Background()
	notFound := fmt.Errorf("policy missing: %w", tcclient.ErrNotFound)

	tests := []struct {
		name    string
		api     *fakeIamPolicyAPI
		ref     string
		wantID  string
		wantErr string
	}{
		{
			name: "direct get",
			api: &fakeIamPolicyAPI{
				getFn: func(_ context.Context, id string) (*clientiam.IamPolicy, error) {
					if id == "pid" {
						return &clientiam.IamPolicy{Identity: "pid", Name: "Deploy", Slug: "deploy"}, nil
					}
					return nil, notFound
				},
				listFn: func(context.Context, *clientiam.ListIamPoliciesRequest) ([]clientiam.IamPolicy, error) {
					return nil, errors.New("list should not be called")
				},
			},
			ref:    "pid",
			wantID: "pid",
		},
		{
			name: "resolve by slug",
			api: &fakeIamPolicyAPI{
				getFn: func(context.Context, string) (*clientiam.IamPolicy, error) {
					return nil, notFound
				},
				listFn: func(context.Context, *clientiam.ListIamPoliciesRequest) ([]clientiam.IamPolicy, error) {
					return []clientiam.IamPolicy{
						{Identity: "pid", Name: "Deploy", Slug: "deploy"},
					}, nil
				},
			},
			ref:    "deploy",
			wantID: "pid",
		},
		{
			name: "resolve by name",
			api: &fakeIamPolicyAPI{
				getFn: func(context.Context, string) (*clientiam.IamPolicy, error) {
					return nil, notFound
				},
				listFn: func(context.Context, *clientiam.ListIamPoliciesRequest) ([]clientiam.IamPolicy, error) {
					return []clientiam.IamPolicy{
						{Identity: "pid", Name: "Deploy Access", Slug: "deploy"},
					}, nil
				},
			},
			ref:    "Deploy Access",
			wantID: "pid",
		},
		{
			name: "empty ref",
			api: &fakeIamPolicyAPI{
				getFn: func(context.Context, string) (*clientiam.IamPolicy, error) {
					return nil, errors.New("should not be called")
				},
				listFn: func(context.Context, *clientiam.ListIamPoliciesRequest) ([]clientiam.IamPolicy, error) {
					return nil, errors.New("should not be called")
				},
			},
			ref:     "  ",
			wantErr: "policy reference is empty",
		},
		{
			name: "not found",
			api: &fakeIamPolicyAPI{
				getFn: func(context.Context, string) (*clientiam.IamPolicy, error) {
					return nil, notFound
				},
				listFn: func(context.Context, *clientiam.ListIamPoliciesRequest) ([]clientiam.IamPolicy, error) {
					return []clientiam.IamPolicy{}, nil
				},
			},
			ref:     "missing",
			wantErr: "IAM policy not found: missing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveIamPolicyRef(ctx, tt.api, tt.ref)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tt.wantID, got.Identity)
		})
	}
}
