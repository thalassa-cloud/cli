package iamresolve

import (
	"context"
	"fmt"
	"strings"

	clientiam "github.com/thalassa-cloud/client-go/iam"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

// OrganisationRoleAPI is implemented by *iam.Client.
type OrganisationRoleAPI interface {
	GetOrganisationRole(ctx context.Context, identity string) (*clientiam.OrganisationRole, error)
	ListOrganisationRoles(ctx context.Context, req *clientiam.ListOrganisationRolesRequest) ([]clientiam.OrganisationRole, error)
}

// ResolveOrganisationRoleRef resolves a user-supplied role identity, slug, or display name.
func ResolveOrganisationRoleRef(ctx context.Context, api OrganisationRoleAPI, ref string) (*clientiam.OrganisationRole, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return nil, fmt.Errorf("role reference is empty")
	}
	role, err := api.GetOrganisationRole(ctx, ref)
	if err == nil {
		return role, nil
	}
	if !tcclient.IsNotFound(err) {
		return nil, fmt.Errorf("get organisation role: %w", err)
	}
	roles, err := api.ListOrganisationRoles(ctx, &clientiam.ListOrganisationRolesRequest{})
	if err != nil {
		return nil, fmt.Errorf("list organisation roles: %w", err)
	}
	for i := range roles {
		r := &roles[i]
		if strings.EqualFold(r.Name, ref) || strings.EqualFold(r.Identity, ref) || strings.EqualFold(r.Slug, ref) {
			return r, nil
		}
	}
	return nil, fmt.Errorf("role not found: %s", ref)
}
