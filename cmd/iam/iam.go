package iam

import (
	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/federatedidentities"
	"github.com/thalassa-cloud/cli/cmd/iam/federatedidentityproviders"
	"github.com/thalassa-cloud/cli/cmd/iam/invites"
	"github.com/thalassa-cloud/cli/cmd/iam/members"
	"github.com/thalassa-cloud/cli/cmd/iam/roles"
	"github.com/thalassa-cloud/cli/cmd/iam/serviceaccounts"
	"github.com/thalassa-cloud/cli/cmd/iam/teams"
	"github.com/thalassa-cloud/cli/cmd/iam/workloadidentityfederation"
)

// IamCmd is the root identity and access management command.
var IamCmd = &cobra.Command{
	Use:   "iam",
	Short: "Identity and access management for your organisation",
	Long: `Manage teams, organisation members, custom roles, federated OIDC identities,
and related resources. Commands apply to the organisation selected in your context
(or the --organisation / -O flag).`,
}

func init() {
	IamCmd.AddCommand(teams.TeamsCmd)
	IamCmd.AddCommand(members.MembersCmd)
	IamCmd.AddCommand(roles.RolesCmd)
	IamCmd.AddCommand(federatedidentities.FederatedIdentitiesCmd)
	IamCmd.AddCommand(federatedidentityproviders.FederatedIdentityProvidersCmd)
	IamCmd.AddCommand(invites.InvitesCmd)
	IamCmd.AddCommand(serviceaccounts.ServiceAccountsCmd)
	IamCmd.AddCommand(workloadidentityfederation.WorkloadIdentityFederationCmd)
}
