package projectresolve

import (
	"context"
	"fmt"
	"strings"

	"github.com/thalassa-cloud/client-go/projects"
)

const (
	// RootRef is the user-facing reference for organisation root (no project scope).
	RootRef = "root"
)

// ProjectsAPI is implemented by *projects.Client.
type ProjectsAPI interface {
	GetProject(ctx context.Context, identity string) (*projects.Project, error)
}

// NormalizeProjectRef maps user-facing root references to an empty project identity.
func NormalizeProjectRef(ref string) string {
	ref = strings.TrimSpace(ref)
	if strings.EqualFold(ref, RootRef) {
		return ""
	}
	return ref
}

// ResolveProjectIdentity resolves a user-supplied project identity or slug to the project identity.
func ResolveProjectIdentity(ctx context.Context, api ProjectsAPI, ref string) (string, error) {
	ref = NormalizeProjectRef(ref)
	if ref == "" {
		return "", nil
	}

	project, err := api.GetProject(ctx, ref)
	if err != nil {
		return "", fmt.Errorf("project not found: %q", ref)
	}
	return project.Identity, nil
}
