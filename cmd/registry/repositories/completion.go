package repositories

import "github.com/thalassa-cloud/cli/internal/completion"

var (
	completeNamespaceID  = completion.CompleteContainerRegistryNamespaceID
	completeRepositoryID = completion.CompleteContainerRegistryRepositoryID
	completeOutputFormat = completion.CompleteOutputFormat
)
