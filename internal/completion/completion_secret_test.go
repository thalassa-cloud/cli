package completion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecretBrowseParent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		toComplete string
		want       string
	}{
		{name: "empty", toComplete: "", want: "/"},
		{name: "root", toComplete: "/", want: "/"},
		{name: "partial root segment", toComplete: "/ap", want: "/"},
		{name: "prefix with slash", toComplete: "/app/", want: "/app/"},
		{name: "nested partial", toComplete: "/app/prod/db", want: "/app/prod/"},
		{name: "missing leading slash", toComplete: "app/prod", want: "/app/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, secretBrowseParent(tt.toComplete))
		})
	}
}
