package workloadidentityfederation

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
)

var (
	flagGitLabRepository string
	flagGitLabRef        string
	flagGitLabRefType    string
	flagGitLabIssuer     string
)

var bootstrapGitLabCmd = &cobra.Command{
	Use:   "gitlab",
	Short: "Bootstrap workload identity for GitLab CI",
	Long: `Creates or reuses a federated identity provider for your GitLab OIDC issuer, then binds the project
and ref to a Thalassa service account.

The GitLab id_token sub uses project_path:<group/project>:ref_type:<type>:ref:<ref>.`,
	Example: `  # GitLab.com, branch main
  tcloud iam workload-identity-federation bootstrap gitlab --repository mygroup/myproject --ref main --role deployer

  # Tag pipeline
  tcloud iam workload-identity-federation bootstrap gitlab --repository mygroup/myproject --ref v1.0.0 --ref-type tag --role deployer

  # Self-managed GitLab
  tcloud iam workload-identity-federation bootstrap gitlab --repository mygroup/myproject --ref main --issuer https://gitlab.example.com --role deployer`,
	Args: cobra.NoArgs,
	RunE: runBootstrapGitLab,
}

func runBootstrapGitLab(cmd *cobra.Command, _ []string) error {
	if strings.TrimSpace(flagGitLabRef) == "" {
		return fmt.Errorf("--ref is required (branch or tag name for the GitLab id_token sub)")
	}
	return executeBootstrap(cmd, BootstrapOptions{
		VCS:           ValueVCSGitLab,
		Repository:    strings.TrimSpace(flagGitLabRepository),
		RefKind:       RefKindBranch,
		Ref:           strings.TrimSpace(flagGitLabRef),
		GitLabRefType: strings.TrimSpace(flagGitLabRefType),
		GitLabIssuer:  strings.TrimSpace(flagGitLabIssuer),
	})
}

func init() {
	f := bootstrapGitLabCmd.Flags()
	f.StringVar(&flagGitLabRepository, "repository", "", "GitLab project path as group/project (required)")
	f.StringVar(&flagGitLabRef, "ref", "", "Git ref segment for id_token sub, e.g. branch or tag name (required)")
	f.StringVar(&flagGitLabRefType, "ref-type", "branch", "ref_type claim in JWT sub: branch, tag, merge_request, etc.")
	f.StringVar(&flagGitLabIssuer, "issuer", "https://gitlab.com", "GitLab OIDC issuer URL (self-managed)")
	_ = bootstrapGitLabCmd.RegisterFlagCompletionFunc("ref-type", completion.CompleteIAMWIFGitLabRefType)

	_ = bootstrapGitLabCmd.MarkFlagRequired("repository")
}
