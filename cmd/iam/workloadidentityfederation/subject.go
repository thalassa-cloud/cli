package workloadidentityfederation

import (
	"fmt"
	"strings"
)

// RefKind selects how the OIDC subject is built.
type RefKind string

const (
	RefKindBranch      RefKind = "branch"
	RefKindTag         RefKind = "tag"
	RefKindEnvironment RefKind = "environment"
	// RefKindPullRequest is GitHub's repo:owner/repo:pull_request subject (no ref segment).
	RefKindPullRequest RefKind = "pull_request"
)

// GitHub issuer URL for Actions OIDC tokens.
const GitHubActionsIssuer = "https://token.actions.githubusercontent.com"

// BuildGitHubSubject returns the JWT `sub` claim GitHub uses for the given repo and ref.
// Repo must be "owner/name" (no leading slash).
func BuildGitHubSubject(repo string, kind RefKind, ref string) (string, error) {
	repo = strings.TrimSpace(repo)
	repo = strings.TrimPrefix(repo, "/")
	repo = strings.TrimSuffix(repo, "/")
	if repo == "" {
		return "", fmt.Errorf("repository is required (owner/name)")
	}
	if strings.Count(repo, "/") < 1 {
		return "", fmt.Errorf("repository must be owner/name, got %q", repo)
	}
	ref = strings.TrimSpace(ref)
	if kind != RefKindPullRequest && ref == "" {
		return "", fmt.Errorf("ref value is required (branch name, tag name, or environment name; omit for pull_request)")
	}
	switch kind {
	case RefKindBranch:
		return fmt.Sprintf("repo:%s:ref:refs/heads/%s", repo, ref), nil
	case RefKindTag:
		return fmt.Sprintf("repo:%s:ref:refs/tags/%s", repo, ref), nil
	case RefKindEnvironment:
		return fmt.Sprintf("repo:%s:environment:%s", repo, ref), nil
	case RefKindPullRequest:
		return fmt.Sprintf("repo:%s:pull_request", repo), nil
	default:
		return "", fmt.Errorf("unsupported ref kind %q for GitHub", kind)
	}
}

// BuildGitLabSubject returns the default GitLab CI id_token `sub` claim:
// project_path:group/project:ref_type:<type>:ref:<ref>
// See: https://docs.gitlab.com/ci/secrets/id_token_authentication/
func BuildGitLabSubject(repoPath, refType, ref string) (string, error) {
	repoPath = strings.TrimSpace(repoPath)
	repoPath = strings.TrimPrefix(repoPath, "/")
	repoPath = strings.TrimSuffix(repoPath, "/")
	if repoPath == "" {
		return "", fmt.Errorf("repository path is required (e.g. group/project)")
	}
	refType = strings.TrimSpace(strings.ToLower(refType))
	if refType == "" {
		refType = "branch"
	}
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return "", fmt.Errorf("ref is required (branch or tag name)")
	}
	return fmt.Sprintf("project_path:%s:ref_type:%s:ref:%s", repoPath, refType, ref), nil
}

// BuildKubernetesSubject returns the JWT sub claim for a Kubernetes bound service account token:
// system:serviceaccount:<namespace>:<name>
// Namespace and name must be given as a single argument "namespace/name".
func BuildKubernetesSubject(namespaceAndName string) (string, error) {
	s := strings.TrimSpace(namespaceAndName)
	s = strings.Trim(s, "/")
	if s == "" {
		return "", fmt.Errorf("namespace/name is required (e.g. default/my-app)")
	}
	parts := strings.SplitN(s, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", fmt.Errorf("kubernetes service account must be namespace/name, got %q", namespaceAndName)
	}
	return fmt.Sprintf("system:serviceaccount:%s:%s", parts[0], parts[1]), nil
}

func normalizeIssuer(issuer string) string {
	return strings.TrimSuffix(strings.TrimSpace(issuer), "/")
}
