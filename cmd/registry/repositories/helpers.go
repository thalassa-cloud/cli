package repositories

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/thalassa"
)

type repositoryAction func(ctx context.Context, client thalassa.Client, namespace, repoID string) error

func runRepositoryActions(
	cmd *cobra.Command,
	args []string,
	force bool,
	summary string,
	startMessage func(repoID string),
	notFoundMessage func(repoID string),
	successMessage func(repoID string),
	actionErrPrefix string,
	action repositoryAction,
) error {
	if namespace == "" {
		return fmt.Errorf("--namespace is required")
	}

	proceed, err := shared.PromptDestructiveUnlessForce(force, summary)
	if err != nil {
		return err
	}
	if !proceed {
		return nil
	}

	client, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	for _, repoID := range args {
		startMessage(repoID)
		if err := action(cmd.Context(), client, namespace, repoID); err != nil {
			if tcclient.IsNotFound(err) {
				notFoundMessage(repoID)
				continue
			}
			return fmt.Errorf("%s: %w", actionErrPrefix, err)
		}
		successMessage(repoID)
	}
	return nil
}
