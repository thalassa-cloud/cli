package iaas

import (
	"context"
	"fmt"
	"strings"

	"github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/client-go/pkg/client"
)

// FindRegionByIdentitySlugOrName finds a region in the given list by matching against
// identity, slug, or name (case-insensitive). Prefers identity/slug (unique) over name.
// Returns the matching region or nil if not found.
func FindRegionByIdentitySlugOrName(regions []iaas.Region, search string) *iaas.Region {
	if search == "" {
		return nil
	}

	// First, try to match by identity or slug (unique identifiers)
	for _, r := range regions {
		if strings.EqualFold(r.Identity, search) || strings.EqualFold(r.Slug, search) {
			return &r
		}
	}

	// If no match by identity/slug, try by name
	for _, r := range regions {
		if strings.EqualFold(r.Name, search) {
			return &r
		}
	}
	return nil
}

// FindRegionByIdentitySlugOrNameWithError finds a region in the given list by matching against
// identity, slug, or name (case-insensitive). Returns the matching region or an error with
// available region slugs if not found.
func FindRegionByIdentitySlugOrNameWithError(regions []iaas.Region, search string) (*iaas.Region, error) {
	region := FindRegionByIdentitySlugOrName(regions, search)
	if region == nil {
		availableRegions := make([]string, 0, len(regions))
		for _, r := range regions {
			availableRegions = append(availableRegions, r.Slug)
		}
		return nil, fmt.Errorf("region not found: %s, available regions: %s", search, strings.Join(availableRegions, ", "))
	}
	return region, nil
}

// FindVPCByIdentitySlugOrName finds a VPC in the given list by matching against
// identity, slug, or name (case-insensitive). Prefers identity/slug (unique) over name.
// Returns the matching VPC or nil if not found.
func FindVPCByIdentitySlugOrName(vpcs []iaas.Vpc, search string) *iaas.Vpc {
	if search == "" {
		return nil
	}

	// First, try to match by identity or slug (unique identifiers)
	for _, v := range vpcs {
		if strings.EqualFold(v.Identity, search) || strings.EqualFold(v.Slug, search) {
			return &v
		}
	}

	// If no match by identity/slug, try by name
	for _, v := range vpcs {
		if strings.EqualFold(v.Name, search) {
			return &v
		}
	}
	return nil
}

// FindVPCByIdentitySlugOrNameWithError finds a VPC in the given list by matching against
// identity, slug, or name (case-insensitive). Returns the matching VPC or an error if not found.
func FindVPCByIdentitySlugOrNameWithError(vpcs []iaas.Vpc, search string) (*iaas.Vpc, error) {
	vpc := FindVPCByIdentitySlugOrName(vpcs, search)
	if vpc == nil {
		return nil, fmt.Errorf("vpc not found: %s", search)
	}
	return vpc, nil
}

// GetVPCByIdentitySlugOrName attempts to get a VPC by first trying GetVpc (for identity),
// and if that fails with NotFound, falls back to listing all VPCs and searching by identity, slug, or name.
// This is more efficient than always listing all VPCs when the identifier is likely an identity.
func GetVPCByIdentitySlugOrName(ctx context.Context, iaasClient *iaas.Client, search string) (*iaas.Vpc, error) {
	if search == "" {
		return nil, fmt.Errorf("vpc identifier is required")
	}

	// First, try to get the VPC directly (most efficient for identities)
	vpc, err := iaasClient.GetVpc(ctx, search)
	if err == nil {
		return vpc, nil
	}

	// If not found, try searching by slug or name
	if !client.IsNotFound(err) {
		return nil, fmt.Errorf("failed to get vpc: %w", err)
	}

	// Fall back to listing and searching
	vpcs, err := iaasClient.ListVpcs(ctx, &iaas.ListVpcsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to list vpcs: %w", err)
	}

	vpc = FindVPCByIdentitySlugOrName(vpcs, search)
	if vpc == nil {
		return nil, fmt.Errorf("vpc not found: %s", search)
	}
	return vpc, nil
}

// FindSubnetByIdentitySlugOrName finds a subnet in the given list by matching against
// identity, slug, or name (case-insensitive). Prefers identity/slug (unique) over name.
// Returns the matching subnet or nil if not found.
func FindSubnetByIdentitySlugOrName(subnets []iaas.Subnet, search string) *iaas.Subnet {
	if search == "" {
		return nil
	}

	// First, try to match by identity or slug (unique identifiers)
	for _, s := range subnets {
		if strings.EqualFold(s.Identity, search) || strings.EqualFold(s.Slug, search) {
			return &s
		}
	}

	// If no match by identity/slug, try by name
	for _, s := range subnets {
		if strings.EqualFold(s.Name, search) {
			return &s
		}
	}
	return nil
}

// FindSubnetByIdentitySlugOrNameWithError finds a subnet in the given list by matching against
// identity, slug, or name (case-insensitive). Prefers identity/slug (unique) over name.
// Returns the matching subnet or an error if not found.
func FindSubnetByIdentitySlugOrNameWithError(subnets []iaas.Subnet, search string) (*iaas.Subnet, error) {
	subnet := FindSubnetByIdentitySlugOrName(subnets, search)
	if subnet == nil {
		return nil, fmt.Errorf("subnet not found: %s", search)
	}
	return subnet, nil
}

// GetSubnetByIdentitySlugOrName attempts to get a subnet by first trying GetSubnet (for identity),
// and if that fails with NotFound, falls back to listing all subnets and searching by identity, slug, or name.
// This is more efficient than always listing all subnets when the identifier is likely an identity.
func GetSubnetByIdentitySlugOrName(ctx context.Context, iaasClient *iaas.Client, search string) (*iaas.Subnet, error) {
	if search == "" {
		return nil, fmt.Errorf("subnet identifier is required")
	}

	// First, try to get the subnet directly (most efficient for identities)
	subnet, err := iaasClient.GetSubnet(ctx, search)
	if err == nil {
		return subnet, nil
	}

	// If not found, try searching by slug or name
	if !client.IsNotFound(err) {
		return nil, fmt.Errorf("failed to get subnet: %w", err)
	}

	// Fall back to listing and searching
	subnets, err := iaasClient.ListSubnets(ctx, &iaas.ListSubnetsRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to list subnets: %w", err)
	}

	subnet = FindSubnetByIdentitySlugOrName(subnets, search)
	if subnet == nil {
		return nil, fmt.Errorf("subnet not found: %s", search)
	}
	return subnet, nil
}

// FindVolumeTypeByIdentitySlugOrName finds a volume type in the given list by matching against
// identity or name (case-insensitive). Prefers identity (unique) over name.
// Returns the matching volume type or nil if not found.
func FindVolumeTypeByIdentitySlugOrName(volumeTypes []iaas.VolumeType, search string) *iaas.VolumeType {
	if search == "" {
		return nil
	}

	// First, try to match by identity (unique identifier)
	for _, vt := range volumeTypes {
		if strings.EqualFold(vt.Identity, search) {
			return &vt
		}
	}

	// If no match by identity, try by name
	for _, vt := range volumeTypes {
		if strings.EqualFold(vt.Name, search) {
			return &vt
		}
	}
	return nil
}

// FindVolumeTypeByIdentitySlugOrNameWithError finds a volume type in the given list by matching against
// identity or name (case-insensitive). Prefers identity (unique) over name.
// Returns the matching volume type or an error if not found.
func FindVolumeTypeByIdentitySlugOrNameWithError(volumeTypes []iaas.VolumeType, search string) (*iaas.VolumeType, error) {
	volumeType := FindVolumeTypeByIdentitySlugOrName(volumeTypes, search)
	if volumeType == nil {
		return nil, fmt.Errorf("volume type not found: %s", search)
	}
	return volumeType, nil
}
