package context

import "github.com/thalassa-cloud/cli/internal/config/contextstate"

// applyOrganisationToContext sets the organisation on ctx. When the organisation
// changes to a non-empty value, the project is cleared because projects are
// organisation-scoped. Returns the previous project when it was cleared.
func applyOrganisationToContext(ctx *contextstate.Context, organisation string) (clearedProject string) {
	if ctx == nil || organisation == "" {
		return ""
	}
	if ctx.Organisation == organisation {
		return ""
	}
	clearedProject = ctx.Project
	ctx.Organisation = organisation
	ctx.Project = ""
	return clearedProject
}
