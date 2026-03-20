package workloadidentityfederation

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/internal/shared"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/config/contextstate"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

// Persistent flags shared by all bootstrap <platform> subcommands.
var (
	flagRole             string
	flagTrustedAudiences []string
	flagScopes           []string
	flagProviderName     string
	flagProviderDesc     string
	flagBootstrapName    string
	flagDryRun           bool
	flagNoHints          bool
)

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Provision workload identity for GitHub, GitLab, or Kubernetes",
	Long: fmt.Sprintf(`Creates (when missing) a federated OIDC identity provider, a Thalassa service account,
a role binding to your organisation role, and a federated identity for the workload JWT subject.

Resources are labelled %s=%s and %s=<github|gitlab|kubernetes>.

Subcommands:
  github      — GitHub Actions (issuer %s)
  gitlab      — GitLab CI id_token
  kubernetes  — in-cluster service account JWTs
`,
		LabelManagedBy, ValueManagedBy, LabelWIFVCS, GitHubActionsIssuer),
}

func executeBootstrap(cmd *cobra.Command, opts BootstrapOptions) error {
	opts.RoleRef = strings.TrimSpace(flagRole)
	opts.ProviderDisplayName = strings.TrimSpace(flagProviderName)
	opts.ProviderDescription = strings.TrimSpace(flagProviderDesc)
	opts.ResourceName = strings.TrimSpace(flagBootstrapName)
	opts.DryRun = flagDryRun

	scopes, err := shared.ParseAccessCredentialScopes(flagScopes)
	if err != nil {
		return err
	}
	opts.AllowedScopes = scopes

	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	apiBase := strings.TrimSuffix(strings.TrimSpace(contextstate.Server()), "/")
	if apiBase == "" {
		apiBase = strings.TrimSuffix(contextstate.DefaultAPIURL, "/")
	}
	trusted := append([]string(nil), flagTrustedAudiences...)
	if len(trusted) == 0 {
		trusted = []string{apiBase}
	}
	opts.TrustedAudiences = trusted

	res, err := RunBootstrap(cmd.Context(), client, opts)
	if err != nil {
		return err
	}

	printBootstrapOutcome(opts.VCS, res, flagDryRun)
	if !flagNoHints && !flagDryRun {
		printBootstrapHints(opts, res, apiBase, strings.TrimSpace(contextstate.Organisation()))
	} else if flagDryRun && !flagNoHints {
		fmt.Println("\n(dry-run: hints omitted; run without --dry-run after applying changes)")
	}
	return nil
}

// parseGitHubRefKind parses --ref-kind for the github bootstrap subcommand.
func parseGitHubRefKind(s string) (RefKind, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "" {
		s = string(RefKindBranch)
	}
	switch RefKind(s) {
	case RefKindBranch, RefKindTag, RefKindEnvironment, RefKindPullRequest:
		return RefKind(s), nil
	default:
		return "", fmt.Errorf("--ref-kind must be branch, tag, environment, or pull_request (got %q)", s)
	}
}

func init() {
	p := bootstrapCmd.PersistentFlags()
	p.StringVar(&flagRole, "role", "", "Organisation role identity, slug, or name (required)")
	p.StringSliceVar(&flagTrustedAudiences, "trusted-audience", nil, "JWT aud values to trust (repeatable; default: current context API URL, e.g. https://api.thalassa.cloud)")
	p.StringSliceVar(&flagScopes, "scope", nil, "Federated identity allowed scopes: api:read, api:write, kubernetes, objectStorage (default: api:read,api:write)")
	p.StringVar(&flagProviderName, "provider-name", "", "Optional display name when creating the federated identity provider")
	p.StringVar(&flagProviderDesc, "provider-description", "", "Optional description when creating the federated identity provider")
	p.StringVar(&flagBootstrapName, "name", "", "Base name for the Thalassa service account and federated identity (federated identity becomes <name>-fi; default: wif-<platform>-<key>)")
	p.BoolVar(&flagDryRun, "dry-run", false, "Print planned changes without calling the API")
	p.BoolVar(&flagNoHints, "no-hints", false, "Do not print platform hints after bootstrap")

	_ = bootstrapCmd.MarkPersistentFlagRequired("role")
	_ = bootstrapCmd.RegisterFlagCompletionFunc("role", completion.CompleteIAMOrganisationRoleIdentityFlag)

	bootstrapCmd.AddCommand(bootstrapGitHubCmd, bootstrapGitLabCmd, bootstrapKubernetesCmd)
}
