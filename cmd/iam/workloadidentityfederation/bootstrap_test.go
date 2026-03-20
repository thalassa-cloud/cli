package workloadidentityfederation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseGitHubRefKind(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    RefKind
		wantErr bool
	}{
		{name: "empty defaults branch", input: "", want: RefKindBranch},
		{name: "branch", input: "branch", want: RefKindBranch},
		{name: "pull_request", input: "pull_request", want: RefKindPullRequest},
		{name: "invalid", input: "merge_request", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseGitHubRefKind(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
