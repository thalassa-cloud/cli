package secrets

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	clientsecrets "github.com/thalassa-cloud/client-go/secrets"
)

func TestParentBrowsePath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path string
		want string
	}{
		{name: "root", path: "/", want: "/"},
		{name: "empty", path: "", want: "/"},
		{name: "one level with slash", path: "/app/", want: "/"},
		{name: "one level without slash", path: "/app", want: "/"},
		{name: "nested", path: "/app/prod/", want: "/app/"},
		{name: "nested without slash", path: "/app/prod/db", want: "/app/prod/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, parentBrowsePath(tt.path))
		})
	}
}

func TestBuildBrowseSelectionLines(t *testing.T) {
	t.Parallel()

	updated := time.Date(2026, 9, 10, 12, 0, 0, 0, time.UTC)
	lines := buildBrowseSelectionLines("/app/", &clientsecrets.BrowseSecretsResponse{
		Path:     "/app/",
		Prefixes: []string{"/app/prod/"},
		Secrets: []clientsecrets.Secret{
			{Path: "/app/config", CurrentVersion: 3, UpdatedAt: updated},
		},
	})

	require.Len(t, lines, 3)
	assert.Equal(t, "..\t..\tgo up", lines[0])
	assert.Equal(t, "/app/prod/\tprod/\tprefix", lines[1])
	assert.Contains(t, lines[2], "/app/config\tconfig\tsecret · v3 ·")
}

func TestBuildBrowseSelectionLinesRootOmitsUp(t *testing.T) {
	t.Parallel()

	lines := buildBrowseSelectionLines("/", &clientsecrets.BrowseSecretsResponse{
		Path:     "/",
		Prefixes: []string{"/app/"},
	})
	require.Len(t, lines, 1)
	assert.Equal(t, "/app/\tapp/\tprefix", lines[0])
}

func TestIsBrowsePrefix(t *testing.T) {
	t.Parallel()

	assert.True(t, isBrowsePrefix("/app/prod/", nil))
	assert.True(t, isBrowsePrefix("/app/prod", []string{"/app/prod"}))
	assert.False(t, isBrowsePrefix("/app/config", []string{"/app/prod/"}))
}
