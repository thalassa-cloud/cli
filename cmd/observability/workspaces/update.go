package workspaces

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientobs "github.com/thalassa-cloud/client-go/observability"
)

var (
	updateName        string
	updateDescription string
	updateRetention   int
	updateLabels      []string
	updateAnnotations []string
	updateWait        bool
	updateWaitTimeout time.Duration
)

var updateCmd = &cobra.Command{
	Use:               "update <workspace>",
	Short:             "Update an observability workspace (only set flags are changed)",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteObservabilityWorkspaceID,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		current, err := client.Observability().GetObservabilityWorkspace(cmd.Context(), args[0])
		if err != nil {
			return fmt.Errorf("failed to get observability workspace: %w", err)
		}

		req := clientobs.UpdateObservabilityWorkspaceRequest{
			Name:          current.Name,
			Description:   current.Description,
			Labels:        current.Labels,
			Annotations:   current.Annotations,
			RetentionDays: current.RetentionDays,
		}
		if cmd.Flags().Changed("name") {
			req.Name = updateName
		}
		if cmd.Flags().Changed("description") {
			req.Description = updateDescription
		}
		if cmd.Flags().Changed("labels") {
			req.Labels = shared.KeyValuePairsToMap(updateLabels)
		}
		if cmd.Flags().Changed("annotations") {
			req.Annotations = shared.KeyValuePairsToMap(updateAnnotations)
		}
		if cmd.Flags().Changed("retention-days") {
			days := updateRetention
			req.RetentionDays = &days
		}
		if req.Labels == nil {
			req.Labels = map[string]string{}
		}
		if req.Annotations == nil {
			req.Annotations = map[string]string{}
		}

		ws, err := client.Observability().UpdateObservabilityWorkspace(cmd.Context(), current.Identity, req)
		if err != nil {
			return fmt.Errorf("failed to update observability workspace: %w", err)
		}

		if updateWait {
			waitCtx := cmd.Context()
			var cancel context.CancelFunc
			if updateWaitTimeout > 0 {
				waitCtx, cancel = context.WithTimeout(cmd.Context(), updateWaitTimeout)
				defer cancel()
			}
			ws, err = client.Observability().WaitUntilObservabilityWorkspaceReady(waitCtx, ws.Identity)
			if err != nil {
				return fmt.Errorf("workspace updated but wait for ready failed: %w", err)
			}
		}

		printWorkspaceDetails(ws)
		return nil
	},
}

func init() {
	WorkspacesCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVar(&updateName, "name", "", "Workspace name")
	updateCmd.Flags().StringVar(&updateDescription, "description", "", "Workspace description")
	updateCmd.Flags().IntVar(&updateRetention, "retention-days", 0, "Retention in days for metrics and logs (1-1095)")
	updateCmd.Flags().StringSliceVar(&updateLabels, "labels", nil, "Replace labels (key=value, repeatable)")
	updateCmd.Flags().StringSliceVar(&updateAnnotations, "annotations", nil, "Replace annotations (key=value, repeatable)")
	updateCmd.Flags().BoolVar(&updateWait, "wait", false, "Wait until the workspace is ready")
	updateCmd.Flags().DurationVar(&updateWaitTimeout, "wait-timeout", 10*time.Minute, "Timeout when waiting with --wait")
}
