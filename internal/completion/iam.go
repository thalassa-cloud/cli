package completion

import (
	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/iamresolve"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientiam "github.com/thalassa-cloud/client-go/iam"
	"github.com/thalassa-cloud/client-go/pkg/base"
)

// CompleteIAMWIFGitHubRefKind completes --ref-kind for GitHub bootstrap.
func CompleteIAMWIFGitHubRefKind(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{
		"branch\tGit branch",
		"tag\tGit tag",
		"environment\tGitHub environment",
		"pull_request\tPull request",
	}, cobra.ShellCompDirectiveNoFileComp
}

// CompleteIAMWIFGitLabRefType completes common --ref-type values for GitLab bootstrap.
func CompleteIAMWIFGitLabRefType(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{
		"branch\tBranch",
		"tag\tTag",
		"merge_request\tMerge request",
		"pipeline\tPipeline",
	}, cobra.ShellCompDirectiveNoFileComp
}

// CompleteIAMPermissionType completes permission values for role rules.
func CompleteIAMPermissionType(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{
		"create\tCreate",
		"read\tRead",
		"update\tUpdate",
		"delete\tDelete",
		"list\tList",
		"*\tAll actions",
	}, cobra.ShellCompDirectiveNoFileComp
}

func iamUserDesc(u *base.AppUser) string {
	if u == nil {
		return ""
	}
	if u.Email != "" {
		return u.Email
	}
	if u.Name != "" {
		return u.Name
	}
	return u.Subject
}

// CompleteIAMOrganisationMemberIdentity completes organisation member record identities.
func CompleteIAMOrganisationMemberIdentity(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	members, err := client.IAM().ListOrganisationMembers(cmd.Context(), &clientiam.ListMembersRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	out := make([]string, 0, len(members))
	for _, m := range members {
		desc := iamUserDesc(m.User)
		if desc != "" {
			out = append(out, m.Identity+"\t"+desc)
		} else {
			out = append(out, m.Identity)
		}
	}
	return out, cobra.ShellCompDirectiveNoFileComp
}

// CompleteIAMAppUserSubject completes user subjects from organisation members (for --user / --user-identity).
func CompleteIAMAppUserSubject(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	members, err := client.IAM().ListOrganisationMembers(cmd.Context(), &clientiam.ListMembersRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	out := make([]string, 0, len(members))
	for _, m := range members {
		if m.User == nil || m.User.Subject == "" {
			continue
		}
		desc := iamUserDesc(m.User)
		if desc != "" && desc != m.User.Subject {
			out = append(out, m.User.Subject+"\t"+desc)
		} else {
			out = append(out, m.User.Subject)
		}
	}
	return out, cobra.ShellCompDirectiveNoFileComp
}

// CompleteIAMOrganisationMemberType completes OWNER / MEMBER.
func CompleteIAMOrganisationMemberType(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{
		string(clientiam.OrganisationMemberTypeOwner) + "\tOrganisation owner",
		string(clientiam.OrganisationMemberTypeMember) + "\tOrganisation member",
	}, cobra.ShellCompDirectiveNoFileComp
}

func completeIAMTeamIdentities(cmd *cobra.Command) ([]string, cobra.ShellCompDirective) {
	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	teams, err := client.IAM().ListTeams(cmd.Context(), &clientiam.ListTeamsRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	out := make([]string, 0, len(teams)*3)
	for _, t := range teams {
		desc := t.Name
		if desc == "" {
			desc = t.Slug
		}
		out = append(out, t.Identity+"\t"+desc)
		if t.Slug != "" && t.Slug != t.Identity {
			out = append(out, t.Slug+"\t"+desc)
		}
		if t.Name != "" && t.Name != t.Identity && t.Name != t.Slug {
			out = append(out, t.Name+"\t"+desc)
		}
	}
	return out, cobra.ShellCompDirectiveNoFileComp
}

// CompleteIAMTeamIdentity completes team identities, slugs, and names (first positional only).
func CompleteIAMTeamIdentity(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return completeIAMTeamIdentities(cmd)
}

// CompleteIAMTeamIdentityFlag completes team identities for flag values (ignores positional args).
func CompleteIAMTeamIdentityFlag(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return completeIAMTeamIdentities(cmd)
}

func completeIAMOrganisationRoleIdentities(cmd *cobra.Command) ([]string, cobra.ShellCompDirective) {
	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	roles, err := client.IAM().ListOrganisationRoles(cmd.Context(), &clientiam.ListOrganisationRolesRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	out := make([]string, 0, len(roles)*3)
	for _, r := range roles {
		desc := r.Name
		if desc == "" {
			desc = r.Slug
		}
		out = append(out, r.Identity+"\t"+desc)
		if r.Slug != "" && r.Slug != r.Identity {
			out = append(out, r.Slug+"\t"+desc)
		}
		if r.Name != "" && r.Name != r.Identity && r.Name != r.Slug {
			out = append(out, r.Name+"\t"+desc)
		}
	}
	return out, cobra.ShellCompDirectiveNoFileComp
}

// CompleteIAMOrganisationRoleIdentity completes organisation role identities and slugs (first positional only).
func CompleteIAMOrganisationRoleIdentity(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return completeIAMOrganisationRoleIdentities(cmd)
}

// CompleteIAMOrganisationRoleIdentityFlag completes organisation roles for flag values (ignores positional args).
func CompleteIAMOrganisationRoleIdentityFlag(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return completeIAMOrganisationRoleIdentities(cmd)
}

func completeIAMServiceAccountIdentities(cmd *cobra.Command) ([]string, cobra.ShellCompDirective) {
	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	list, err := client.IAM().ListServiceAccounts(cmd.Context(), &clientiam.ListServiceAccountsRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	out := make([]string, 0, len(list)*3)
	for _, sa := range list {
		desc := sa.Name
		if desc == "" {
			desc = sa.Slug
		}
		out = append(out, sa.Identity+"\t"+desc)
		if sa.Slug != "" && sa.Slug != sa.Identity {
			out = append(out, sa.Slug+"\t"+desc)
		}
		if sa.Name != "" && sa.Name != sa.Identity && sa.Name != sa.Slug {
			out = append(out, sa.Name+"\t"+desc)
		}
	}
	return out, cobra.ShellCompDirectiveNoFileComp
}

// CompleteIAMServiceAccountIdentity completes service account identities (first positional only).
func CompleteIAMServiceAccountIdentity(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return completeIAMServiceAccountIdentities(cmd)
}

// CompleteIAMServiceAccountIdentityFlag completes service accounts for flag values (ignores positional args).
func CompleteIAMServiceAccountIdentityFlag(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return completeIAMServiceAccountIdentities(cmd)
}

func completeIAMFederatedIdentityIdentities(cmd *cobra.Command) ([]string, cobra.ShellCompDirective) {
	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	list, err := client.IAM().ListFederatedIdentities(cmd.Context(), &clientiam.ListFederatedIdentitiesRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	out := make([]string, 0, len(list)*2)
	for _, fi := range list {
		if fi.Name != "" && fi.Name != fi.Identity {
			out = append(out, fi.Identity+"\t"+fi.Name)
			out = append(out, fi.Name+"\t"+fi.Identity)
		} else {
			out = append(out, fi.Identity)
		}
	}
	return out, cobra.ShellCompDirectiveNoFileComp
}

// CompleteIAMFederatedIdentityIdentity completes federated identity identifiers (first positional only).
func CompleteIAMFederatedIdentityIdentity(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return completeIAMFederatedIdentityIdentities(cmd)
}

// CompleteIAMFederatedIdentityIdentityFlag completes federated identities for flag values.
func CompleteIAMFederatedIdentityIdentityFlag(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return completeIAMFederatedIdentityIdentities(cmd)
}

func completeIAMFederatedIdentityProviderIdentities(cmd *cobra.Command) ([]string, cobra.ShellCompDirective) {
	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	list, err := client.IAM().ListFederatedIdentityProviders(cmd.Context(), &clientiam.ListFederatedIdentityProvidersRequest{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	out := make([]string, 0, len(list)*2)
	for _, p := range list {
		if p.Name != "" && p.Name != p.Identity {
			out = append(out, p.Identity+"\t"+p.Name)
			out = append(out, p.Name+"\t"+p.Identity)
		} else {
			out = append(out, p.Identity)
		}
	}
	return out, cobra.ShellCompDirectiveNoFileComp
}

// CompleteIAMFederatedIdentityProviderIdentity completes federated identity provider identifiers (first positional only).
func CompleteIAMFederatedIdentityProviderIdentity(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	return completeIAMFederatedIdentityProviderIdentities(cmd)
}

// CompleteIAMFederatedIdentityProviderIdentityFlag completes federated identity providers for flag values.
func CompleteIAMFederatedIdentityProviderIdentityFlag(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return completeIAMFederatedIdentityProviderIdentities(cmd)
}

// CompleteIAMRoleThenBinding completes the role then binding identities for that role.
func CompleteIAMRoleThenBinding(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	switch len(args) {
	case 0:
		return CompleteIAMOrganisationRoleIdentity(cmd, args, toComplete)
	case 1:
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		role, err := iamresolve.ResolveOrganisationRoleRef(cmd.Context(), client.IAM(), args[0])
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		bindings, err := client.IAM().ListRoleBindings(cmd.Context(), role.Identity, &clientiam.ListRoleBindingsRequest{})
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		out := make([]string, 0, len(bindings)*3)
		for _, b := range bindings {
			desc := b.Name
			if desc == "" {
				desc = b.Slug
			}
			out = append(out, b.Identity+"\t"+desc)
			if b.Slug != "" && b.Slug != b.Identity {
				out = append(out, b.Slug+"\t"+desc)
			}
			if b.Name != "" && b.Name != b.Identity && b.Name != b.Slug {
				out = append(out, b.Name+"\t"+desc)
			}
		}
		return out, cobra.ShellCompDirectiveNoFileComp
	default:
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
}

// CompleteIAMRoleThenRule completes the role then rule identities on that role.
func CompleteIAMRoleThenRule(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	switch len(args) {
	case 0:
		return CompleteIAMOrganisationRoleIdentity(cmd, args, toComplete)
	case 1:
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		role, err := client.IAM().GetOrganisationRole(cmd.Context(), args[0])
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		out := make([]string, 0, len(role.Rules))
		for _, ru := range role.Rules {
			note := ru.Note
			if note == "" {
				note = "rule"
			}
			out = append(out, ru.Identity+"\t"+note)
		}
		return out, cobra.ShellCompDirectiveNoFileComp
	default:
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
}

// CompleteIAMTeamThenTeamMember completes the team then team-member identities for that team.
func CompleteIAMTeamThenTeamMember(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	switch len(args) {
	case 0:
		return CompleteIAMTeamIdentity(cmd, args, toComplete)
	case 1:
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		team, err := client.IAM().GetTeam(cmd.Context(), args[0], &clientiam.GetTeamRequest{})
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		out := make([]string, 0, len(team.Members))
		for _, m := range team.Members {
			desc := iamUserDesc(&m.User)
			if desc != "" {
				out = append(out, m.Identity+"\t"+desc)
			} else {
				out = append(out, m.Identity)
			}
		}
		return out, cobra.ShellCompDirectiveNoFileComp
	default:
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
}
