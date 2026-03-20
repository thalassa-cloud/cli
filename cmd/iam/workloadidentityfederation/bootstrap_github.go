package workloadidentityfederation

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
)

var (
	flagGitHubRepository string
	flagGitHubRefKind    string
	flagGitHubRef        string
)

var bootstrapGitHubCmd = &cobra.Command{
	Use:   "github",
	Short: "Bootstrap workload identity for GitHub Actions",
	Long: `Creates or reuses a federated identity provider for GitHub's OIDC issuer, then binds your repository
and ref to a Thalassa service account via a federated identity.

The JWT issuer is https://token.actions.githubusercontent.com. Match subjects with --repository (owner/name),
--ref-kind, and --ref (omit --ref only when --ref-kind is pull_request).`,
	Example: `  # Main branch (JWT aud defaults to context API URL)
  tcloud iam workload-identity-federation bootstrap github --repository acme/api --ref main --role deployer

  # Specific ref kind
  tcloud iam workload-identity-federation bootstrap github --repository acme/api --ref-kind branch --ref main --role deployer

  # Pull request workflows (subject repo:owner/repo:pull_request)
  tcloud iam workload-identity-federation bootstrap github --repository acme/api --ref-kind pull_request --role deployer

  # Custom Thalassa service account / federated identity names (federated identity: platform-ci-fi)
  tcloud iam workload-identity-federation bootstrap github --repository acme/api --ref main --name platform-ci --role deployer`,
	Args: cobra.NoArgs,
	RunE: runBootstrapGitHub,
}

func runBootstrapGitHub(cmd *cobra.Command, _ []string) error {
	refKind, err := parseGitHubRefKind(flagGitHubRefKind)
	if err != nil {
		return err
	}
	if refKind != RefKindPullRequest && strings.TrimSpace(flagGitHubRef) == "" {
		return fmt.Errorf("--ref is required unless --ref-kind is pull_request")
	}
	return executeBootstrap(cmd, BootstrapOptions{
		VCS:        ValueVCSGitHub,
		Repository: strings.TrimSpace(flagGitHubRepository),
		RefKind:    refKind,
		Ref:        strings.TrimSpace(flagGitHubRef),
	})
}

func init() {
	f := bootstrapGitHubCmd.Flags()
	f.StringVar(&flagGitHubRepository, "repository", "", "GitHub repository as owner/name (required)")
	f.StringVar(&flagGitHubRefKind, "ref-kind", "branch", "branch, tag, environment, or pull_request")
	f.StringVar(&flagGitHubRef, "ref", "", "Branch, tag, or environment name (omit when --ref-kind pull_request)")
	_ = bootstrapGitHubCmd.RegisterFlagCompletionFunc("ref-kind", completion.CompleteIAMWIFGitHubRefKind)

	_ = bootstrapGitHubCmd.MarkFlagRequired("repository")
}
