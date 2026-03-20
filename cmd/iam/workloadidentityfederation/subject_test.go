package workloadidentityfederation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildGitHubSubject(t *testing.T) {
	tests := []struct {
		name    string
		repo    string
		kind    RefKind
		ref     string
		want    string
		wantErr bool
	}{
		{name: "branch", repo: "acme/app", kind: RefKindBranch, ref: "main", want: "repo:acme/app:ref:refs/heads/main"},
		{name: "tag", repo: "acme/app", kind: RefKindTag, ref: "v1.0.0", want: "repo:acme/app:ref:refs/tags/v1.0.0"},
		{name: "environment", repo: "acme/app", kind: RefKindEnvironment, ref: "production", want: "repo:acme/app:environment:production"},
		{name: "trim repo", repo: "/acme/app/", kind: RefKindBranch, ref: "main", want: "repo:acme/app:ref:refs/heads/main"},
		{name: "bad repo", repo: "nope", kind: RefKindBranch, ref: "main", wantErr: true},
		{name: "empty ref", repo: "a/b", kind: RefKindBranch, ref: "  ", wantErr: true},
		{name: "pull_request", repo: "acme/app", kind: RefKindPullRequest, ref: "", want: "repo:acme/app:pull_request"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildGitHubSubject(tt.repo, tt.kind, tt.ref)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBuildGitLabSubject(t *testing.T) {
	got, err := BuildGitLabSubject("mygroup/myproj", "branch", "main")
	require.NoError(t, err)
	assert.Equal(t, "project_path:mygroup/myproj:ref_type:branch:ref:main", got)

	got, err = BuildGitLabSubject("g/p", "", "v1")
	require.NoError(t, err)
	assert.Equal(t, "project_path:g/p:ref_type:branch:ref:v1", got)

	_, err = BuildGitLabSubject("", "branch", "main")
	require.Error(t, err)
}

func TestNormalizeIssuer(t *testing.T) {
	assert.Equal(t, "https://gitlab.com", normalizeIssuer("https://gitlab.com/"))
}

func TestBuildKubernetesSubject(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "ok", input: "default/my-sa", want: "system:serviceaccount:default:my-sa"},
		{name: "trim", input: " /prod/worker/ ", want: "system:serviceaccount:prod:worker"},
		{name: "empty", input: "", wantErr: true},
		{name: "no slash", input: "default", wantErr: true},
		{name: "empty ns", input: "/sa", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildKubernetesSubject(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
