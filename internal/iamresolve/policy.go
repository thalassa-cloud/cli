package iamresolve

import (
	"context"
	"fmt"
	"strings"

	clientiam "github.com/thalassa-cloud/client-go/iam"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

// IamPolicyAPI is implemented by *iam.Client.
type IamPolicyAPI interface {
	GetIamPolicy(ctx context.Context, identity string) (*clientiam.IamPolicy, error)
	ListIamPolicies(ctx context.Context, req *clientiam.ListIamPoliciesRequest) ([]clientiam.IamPolicy, error)
}

// ResolveIamPolicyRef resolves a user-supplied IAM policy identity, slug, or display name.
func ResolveIamPolicyRef(ctx context.Context, api IamPolicyAPI, ref string) (*clientiam.IamPolicy, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return nil, fmt.Errorf("policy reference is empty")
	}
	policy, err := api.GetIamPolicy(ctx, ref)
	if err == nil {
		return policy, nil
	}
	if !tcclient.IsNotFound(err) {
		return nil, fmt.Errorf("get IAM policy: %w", err)
	}
	policies, err := api.ListIamPolicies(ctx, &clientiam.ListIamPoliciesRequest{})
	if err != nil {
		return nil, fmt.Errorf("list IAM policies: %w", err)
	}
	for i := range policies {
		p := &policies[i]
		if strings.EqualFold(p.Name, ref) || strings.EqualFold(p.Identity, ref) || strings.EqualFold(p.Slug, ref) {
			return p, nil
		}
	}
	return nil, fmt.Errorf("IAM policy not found: %s", ref)
}
