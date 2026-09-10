package snapshotpolicies

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	updateName        string
	updateDescription string
	updateTTL         string
	updateKeepCount   int
	updateEnabled     bool
	updateSchedule    string
	updateTimezone    string
	updateTargetType  string
	updateSelectors   []string
	updateVolumes     []string
	updateLabels      []string
	updateAnnotations []string
)

var updateCmd = &cobra.Command{
	Use:               "update",
	Short:             "Update a snapshot policy",
	Long:              "Update an existing snapshot policy. Unspecified fields are preserved from the current resource.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteSnapshotPolicyID,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		current, err := client.IaaS().GetSnapshotPolicy(cmd.Context(), args[0])
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("snapshot policy not found: %s", args[0])
			}
			return fmt.Errorf("failed to get snapshot policy: %w", err)
		}

		req := iaas.UpdateSnapshotPolicyRequest{
			Name:        current.Name,
			Description: current.Description,
			Labels:      current.Labels,
			Annotations: current.Annotations,
			Ttl:         current.Ttl,
			KeepCount:   current.KeepCount,
			Enabled:     current.Enabled,
			Schedule:    current.Schedule,
			Timezone:    current.Timezone,
			Target:      current.Target,
		}

		if cmd.Flags().Changed("name") {
			req.Name = updateName
		}
		if cmd.Flags().Changed("description") {
			req.Description = updateDescription
		}
		if cmd.Flags().Changed("labels") {
			req.Labels = parseKeyValueSlice(updateLabels)
		}
		if cmd.Flags().Changed("annotations") {
			req.Annotations = parseKeyValueSlice(updateAnnotations)
		}
		if cmd.Flags().Changed("ttl") {
			ttl, err := parseTTL(updateTTL)
			if err != nil {
				return err
			}
			req.Ttl = ttl
		}
		if cmd.Flags().Changed("keep-count") {
			keep := updateKeepCount
			req.KeepCount = &keep
		}
		if cmd.Flags().Changed("enabled") {
			req.Enabled = updateEnabled
		}
		if cmd.Flags().Changed("schedule") {
			req.Schedule = updateSchedule
		}
		if cmd.Flags().Changed("timezone") {
			req.Timezone = updateTimezone
		}

		targetChanged := cmd.Flags().Changed("target-type") || cmd.Flags().Changed("selector") || cmd.Flags().Changed("volumes")
		if targetChanged {
			targetType := string(current.Target.Type)
			selectors := updateSelectors
			volumes := updateVolumes
			if cmd.Flags().Changed("target-type") {
				targetType = updateTargetType
			}
			if !cmd.Flags().Changed("selector") {
				selectors = make([]string, 0, len(current.Target.Selector))
				for k, v := range current.Target.Selector {
					selectors = append(selectors, k+"="+v)
				}
			}
			if !cmd.Flags().Changed("volumes") {
				volumes = current.Target.VolumeIdentities
			}
			target, err := buildTarget(targetType, selectors, volumes)
			if err != nil {
				return err
			}
			req.Target = target
		}

		policy, err := client.IaaS().UpdateSnapshotPolicy(cmd.Context(), current.Identity, req)
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
	SnapshotPoliciesCmd.AddCommand(updateCmd)

	updateCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	updateCmd.Flags().StringVar(&updateName, "name", "", "Name of the snapshot policy")
	updateCmd.Flags().StringVar(&updateDescription, "description", "", "Description of the snapshot policy")
	updateCmd.Flags().StringVar(&updateTTL, "ttl", "", "Snapshot retention duration (e.g. 24h, 168h)")
	updateCmd.Flags().IntVar(&updateKeepCount, "keep-count", 0, "Maximum number of snapshots to retain")
	updateCmd.Flags().BoolVar(&updateEnabled, "enabled", true, "Enable or disable the snapshot policy")
	updateCmd.Flags().StringVar(&updateSchedule, "schedule", "", "Cron schedule for snapshot creation")
	updateCmd.Flags().StringVar(&updateTimezone, "timezone", "", "Timezone for the schedule")
	updateCmd.Flags().StringVar(&updateTargetType, "target-type", "", "Target type: selector or explicit")
	updateCmd.Flags().StringSliceVar(&updateSelectors, "selector", []string{}, "Label selectors for target volumes (key=value)")
	updateCmd.Flags().StringSliceVar(&updateVolumes, "volumes", []string{}, "Volume identities when target-type is explicit")
	updateCmd.Flags().StringSliceVar(&updateLabels, "labels", []string{}, "Labels in key=value format")
	updateCmd.Flags().StringSliceVar(&updateAnnotations, "annotations", []string{}, "Annotations in key=value format")

	_ = updateCmd.RegisterFlagCompletionFunc("target-type", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"selector", "explicit"}, cobra.ShellCompDirectiveNoFileComp
	})
	_ = updateCmd.RegisterFlagCompletionFunc("volumes", completion.CompleteVolumeID)
}
