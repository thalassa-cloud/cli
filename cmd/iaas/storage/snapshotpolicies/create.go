package snapshotpolicies

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	iaasutil "github.com/thalassa-cloud/cli/internal/iaas"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
)

var (
	createName        string
	createDescription string
	createRegion      string
	createTTL         string
	createKeepCount   int
	createEnabled     bool
	createSchedule    string
	createTimezone    string
	createTargetType  string
	createSelectors   []string
	createVolumes     []string
	createLabels      []string
	createAnnotations []string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a snapshot policy",
	Long:  "Create a new automated snapshot policy for volumes in a region.",
	Example: `tcloud storage snapshot-policies create --name daily --region nl-ams --schedule "0 2 * * *" --ttl 168h --timezone UTC --target-type selector --selector backup=true
tcloud storage snapshot-policies create --name weekly --region nl-ams --schedule "0 3 * * 0" --ttl 720h --target-type explicit --volumes vol-1,vol-2`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if createName == "" {
			return fmt.Errorf("name is required")
		}
		if createRegion == "" {
			return fmt.Errorf("region is required")
		}
		if createSchedule == "" {
			return fmt.Errorf("schedule is required")
		}

		ttl, err := parseTTL(createTTL)
		if err != nil {
			return err
		}
		target, err := buildTarget(createTargetType, createSelectors, createVolumes)
		if err != nil {
			return err
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		regions, err := client.IaaS().ListRegions(cmd.Context(), &iaas.ListRegionsRequest{})
		if err != nil {
			return fmt.Errorf("failed to list regions: %w", err)
		}
		region, err := iaasutil.FindRegionByIdentitySlugOrNameWithError(regions, createRegion)
		if err != nil {
			return err
		}

		req := iaas.CreateSnapshotPolicyRequest{
			Name:        createName,
			Description: createDescription,
			Labels:      parseKeyValueSlice(createLabels),
			Annotations: parseKeyValueSlice(createAnnotations),
			Region:      region.Identity,
			Ttl:         ttl,
			Enabled:     createEnabled,
			Schedule:    createSchedule,
			Timezone:    createTimezone,
			Target:      target,
		}
		if cmd.Flags().Changed("keep-count") {
			keep := createKeepCount
			req.KeepCount = &keep
		}

		policy, err := client.IaaS().CreateSnapshotPolicy(cmd.Context(), req)
		if err != nil {
			return err
		}

		body := [][]string{{
			policy.Identity,
			policy.Name,
			regionName(*policy),
			fmt.Sprintf("%t", policy.Enabled),
			policy.Schedule,
			policy.Ttl.String(),
			keepCountString(policy.KeepCount),
			targetSummary(policy.Target),
		}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Region", "Enabled", "Schedule", "TTL", "Keep", "Target"}, body)
		}
		return nil
	},
}

func init() {
	SnapshotPoliciesCmd.AddCommand(createCmd)

	createCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	createCmd.Flags().StringVar(&createName, "name", "", "Name of the snapshot policy")
	createCmd.Flags().StringVar(&createDescription, "description", "", "Description of the snapshot policy")
	createCmd.Flags().StringVar(&createRegion, "region", "", "Region of the snapshot policy")
	createCmd.Flags().StringVar(&createTTL, "ttl", "", "Snapshot retention duration (e.g. 24h, 168h)")
	createCmd.Flags().IntVar(&createKeepCount, "keep-count", 0, "Maximum number of snapshots to retain")
	createCmd.Flags().BoolVar(&createEnabled, "enabled", true, "Enable the snapshot policy")
	createCmd.Flags().StringVar(&createSchedule, "schedule", "", "Cron schedule for snapshot creation")
	createCmd.Flags().StringVar(&createTimezone, "timezone", "UTC", "Timezone for the schedule")
	createCmd.Flags().StringVar(&createTargetType, "target-type", "", "Target type: selector or explicit")
	createCmd.Flags().StringSliceVar(&createSelectors, "selector", []string{}, "Label selectors for target volumes (key=value)")
	createCmd.Flags().StringSliceVar(&createVolumes, "volumes", []string{}, "Volume identities when target-type is explicit")
	createCmd.Flags().StringSliceVar(&createLabels, "labels", []string{}, "Labels in key=value format")
	createCmd.Flags().StringSliceVar(&createAnnotations, "annotations", []string{}, "Annotations in key=value format")

	_ = createCmd.MarkFlagRequired("name")
	_ = createCmd.MarkFlagRequired("region")
	_ = createCmd.MarkFlagRequired("ttl")
	_ = createCmd.MarkFlagRequired("schedule")
	_ = createCmd.MarkFlagRequired("target-type")

	_ = createCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
	_ = createCmd.RegisterFlagCompletionFunc("target-type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"selector", "explicit"}, cobra.ShellCompDirectiveNoFileComp
	})
	_ = createCmd.RegisterFlagCompletionFunc("volumes", completion.CompleteVolumeID)
}
