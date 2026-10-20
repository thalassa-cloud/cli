package context

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalassa-cloud/cli/internal/config/contextstate"
)

func TestApplyOrganisationToContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		ctx                *contextstate.Context
		organisation       string
		wantClearedProject string
		wantOrganisation   string
		wantProject        string
	}{
		{
			name: "same organisation keeps project",
			ctx: &contextstate.Context{
				Organisation: "acme",
				Project:      "proj-1",
			},
			organisation:       "acme",
			wantClearedProject: "",
			wantOrganisation:   "acme",
			wantProject:        "proj-1",
		},
		{
			name: "different organisation clears project",
			ctx: &contextstate.Context{
				Organisation: "acme",
				Project:      "proj-1",
			},
			organisation:       "other",
			wantClearedProject: "proj-1",
			wantOrganisation:   "other",
			wantProject:        "",
		},
		{
			name: "different organisation with empty project",
			ctx: &contextstate.Context{
				Organisation: "acme",
				Project:      "",
			},
			organisation:       "other",
			wantClearedProject: "",
			wantOrganisation:   "other",
			wantProject:        "",
		},
		{
			name: "empty organisation leaves context unchanged",
			ctx: &contextstate.Context{
				Organisation: "acme",
				Project:      "proj-1",
			},
			organisation:       "",
			wantClearedProject: "",
			wantOrganisation:   "acme",
			wantProject:        "proj-1",
		},
		{
			name:               "nil context is a no-op",
			ctx:                nil,
			organisation:       "other",
			wantClearedProject: "",
		},
		{
			name: "sets organisation when previously empty",
			ctx: &contextstate.Context{
				Organisation: "",
				Project:      "proj-1",
			},
			organisation:       "acme",
			wantClearedProject: "proj-1",
			wantOrganisation:   "acme",
			wantProject:        "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cleared := applyOrganisationToContext(tt.ctx, tt.organisation)
			assert.Equal(t, tt.wantClearedProject, cleared)
			if tt.ctx == nil {
				return
			}
			require.NotNil(t, tt.ctx)
			assert.Equal(t, tt.wantOrganisation, tt.ctx.Organisation)
			assert.Equal(t, tt.wantProject, tt.ctx.Project)
		})
	}
}
