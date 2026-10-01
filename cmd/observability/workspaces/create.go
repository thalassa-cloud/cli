package workspaces

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	iaasutil "github.com/thalassa-cloud/cli/internal/iaas"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
	clientobs "github.com/thalassa-cloud/client-go/observability"
)

var (
	createName        string
	createDescription string
	createRegion      string
	createRetention   int
	createLabels      []string
	createAnnotations []string
	createWait        bool
	createWaitTimeout time.Duration
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an observability workspace",
	Example: `  tcloud observability workspaces create --name prod-metrics --region nl-ams
  tcloud observability workspaces create --name prod-metrics --region nl-ams --retention-days 30 --wait`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		if createName == "" {
			return fmt.Errorf("--name is required")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		req := clientobs.CreateObservabilityWorkspaceRequest{
			Name:        createName,
			Description: createDescription,
			Labels:      shared.KeyValuePairsToMap(createLabels),
			Annotations: shared.KeyValuePairsToMap(createAnnotations),
		}
		if cmd.Flags().Changed("retention-days") {
			days := createRetention
			req.RetentionDays = &days
		}
		if createRegion != "" {
			regions, err := client.IaaS().ListRegions(cmd.Context(), &iaas.ListRegionsRequest{})
			if err != nil {
				return fmt.Errorf("failed to list regions: %w", err)
			}
			region, err := iaasutil.FindRegionByIdentitySlugOrNameWithError(regions, createRegion)
			if err != nil {
				return err
			}
			req.RegionIdentity = region.Identity
		}

		ws, err := client.Observability().CreateObservabilityWorkspace(cmd.Context(), req)
		if err != nil {
			return fmt.Errorf("failed to create observability workspace: %w", err)
		}

		if createWait {
			waitCtx := cmd.Context()
			var cancel context.CancelFunc
			if createWaitTimeout > 0 {
				waitCtx, cancel = context.WithTimeout(cmd.Context(), createWaitTimeout)
				defer cancel()
			}
			ws, err = client.Observability().WaitUntilObservabilityWorkspaceReady(waitCtx, ws.Identity)
			if err != nil {
				return fmt.Errorf("workspace created but wait for ready failed: %w", err)
			}
		}

		printWorkspaceDetails(ws)
		return nil
	},
}

func init() {
	WorkspacesCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&createName, "name", "", "Workspace name (required)")
	createCmd.Flags().StringVar(&createDescription, "description", "", "Workspace description")
	createCmd.Flags().StringVar(&createRegion, "region", "", "Region identity, slug, or name")
	createCmd.Flags().IntVar(&createRetention, "retention-days", 0, "Retention in days for metrics and logs (1-1095)")
	createCmd.Flags().StringSliceVar(&createLabels, "labels", nil, "Labels as key=value (repeatable)")
	createCmd.Flags().StringSliceVar(&createAnnotations, "annotations", nil, "Annotations as key=value (repeatable)")
	createCmd.Flags().BoolVar(&createWait, "wait", false, "Wait until the workspace is ready")
	createCmd.Flags().DurationVar(&createWaitTimeout, "wait-timeout", 10*time.Minute, "Timeout when waiting with --wait")
	_ = createCmd.MarkFlagRequired("name")
	_ = createCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
}
