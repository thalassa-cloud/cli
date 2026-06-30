package registry

import "github.com/thalassa-cloud/cli/internal/completion"

var (
	CompleteNamespaceID   = completion.CompleteContainerRegistryNamespaceID
	CompleteRepositoryID  = completion.CompleteContainerRegistryRepositoryID
	CompleteRegion        = completion.CompleteRegion
	CompleteOutputFormat  = completion.CompleteOutputFormat
)
