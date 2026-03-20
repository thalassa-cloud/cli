package workloadidentityfederation

import "github.com/spf13/cobra"

// WorkloadIdentityFederationCmd groups workload identity federation helpers.
var WorkloadIdentityFederationCmd = &cobra.Command{
	Use:     "workload-identity-federation",
	Aliases: []string{"wif"},
	Short:   "Bootstrap and manage CI/CD workload identity (OIDC)",
	Long: `Commands to provision federated identity providers, service accounts, role bindings,
and federated identities for GitHub Actions, GitLab CI, and Kubernetes service account OIDC tokens.`,
}

func init() {
	WorkloadIdentityFederationCmd.AddCommand(bootstrapCmd)
}
