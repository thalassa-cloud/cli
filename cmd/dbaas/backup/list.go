package backup

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/labels"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/dbaas"
	"github.com/thalassa-cloud/client-go/filters"
)

var (
	backupListLabelSelector string
	backupListClusterFilter string
	backupListOlderThan     string
	backupListNewerThan     string
	backupListStatusFilter  []string
	backupListNoHeader      bool
	backupListShowExactTime bool
	backupListShowLabels    bool
)

// parseDuration parses a duration string supporting days (d), weeks (w), months (mo), years (y)
// in addition to standard Go duration units (h, m, s, etc.)
func parseDuration(s string) (time.Duration, error) {
	// First try standard time.ParseDuration
	if d, err := time.ParseDuration(s); err == nil {
		return d, nil
	}

	// Handle custom units: d, w, mo, y
	var total time.Duration
	re := regexp.MustCompile(`(\d+)([dwy]|mo)`)
	matches := re.FindAllStringSubmatch(s, -1)

	if len(matches) == 0 {
		return 0, fmt.Errorf("invalid duration format: %s", s)
	}

	for _, match := range matches {
		if len(match) != 3 {
			continue
		}
		value, err := strconv.Atoi(match[1])
		if err != nil {
			return 0, fmt.Errorf("invalid number in duration: %s", match[1])
		}

		unit := match[2]
		switch unit {
		case "d":
			total += time.Duration(value) * 24 * time.Hour
		case "w":
			total += time.Duration(value) * 7 * 24 * time.Hour
		case "mo":
			// Approximate month as 30 days
			total += time.Duration(value) * 30 * 24 * time.Hour
		case "y":
			// Approximate year as 365 days
			total += time.Duration(value) * 365 * 24 * time.Hour
		}
	}

	// Check if there are remaining characters not matched
	remaining := re.ReplaceAllString(s, "")
	if strings.TrimSpace(remaining) != "" {
		return 0, fmt.Errorf("invalid duration format: %s (unrecognized: %s)", s, remaining)
	}

	return total, nil
}

// backupListCmd represents the backup list command
var backupListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List database backups",
	Long:    "List database backups for a specific cluster or all backups in the organisation",
	Aliases: []string{"ls", "get"},
	Args:    cobra.MaximumNArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return completion.CompleteDbClusterID(cmd, args, toComplete)
		}
		return nil, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		var backups []dbaas.DbClusterBackup

		// Build filters
		f := []filters.Filter{}

		// Add label selector filter if provided
		if backupListLabelSelector != "" {
			f = append(f, &filters.LabelFilter{
				MatchLabels: labels.ParseLabelSelector(backupListLabelSelector),
			})
		}

		// Add cluster filter if provided via flag
		if backupListClusterFilter != "" {
			f = append(f, &filters.FilterKeyValue{
				Key:   "dbCluster",
				Value: backupListClusterFilter,
			})
		}

		if len(backupListStatusFilter) > 0 {
			f = append(f, &filters.FilterKeyValue{
				Key:   "status",
				Value: strings.Join(backupListStatusFilter, ","),
			})
		}

		listRequest := &dbaas.ListDbBackupsRequest{
			Filters: f,
		}

		// If cluster identity is provided as argument, list backups for that cluster
		if len(args) > 0 {
			clusterIdentity := args[0]
			backups, err = client.DBaaS().ListDbBackupsForDbCluster(cmd.Context(), clusterIdentity, listRequest)
			if err != nil {
				return fmt.Errorf("failed to list backups for cluster: %w", err)
			}
		} else {
			// Otherwise list all backups for the organisation
			backups, err = client.DBaaS().ListDbBackupsForOrganisation(cmd.Context(), listRequest)
			if err != nil {
				return fmt.Errorf("failed to list backups: %w", err)
			}
		}

		// Filter by age and status if specified
		now := time.Now()
		if backupListOlderThan != "" || backupListNewerThan != "" {
			filteredBackups := []dbaas.DbClusterBackup{}
			for _, backup := range backups {
				// Filter by age
				backupAge := now.Sub(backup.CreatedAt)

				// Filter by older-than
				if backupListOlderThan != "" {
					duration, err := parseDuration(backupListOlderThan)
					if err != nil {
						return fmt.Errorf("invalid --older-than duration: %w", err)
					}
					if backupAge < duration {
						continue // Skip backups that are not old enough
					}
				}

				// Filter by newer-than
				if backupListNewerThan != "" {
					duration, err := parseDuration(backupListNewerThan)
					if err != nil {
						return fmt.Errorf("invalid --newer-than duration: %w", err)
					}
					if backupAge > duration {
						continue // Skip backups that are too old
					}
				}

				filteredBackups = append(filteredBackups, backup)
			}
			backups = filteredBackups
		}

		if len(backups) == 0 {
			fmt.Println("No backups found")
			return nil
		}

		body := make([][]string, 0, len(backups))
		for _, backup := range backups {
			clusterName := ""
			if backup.DbCluster != nil {
				clusterName = backup.DbCluster.Name
			}

			status := string(backup.Status)
			if backup.DeleteScheduledAt != nil {
				status = fmt.Sprintf("%s (deletion scheduled)", status)
			}

			trigger := string(backup.BackupTrigger)
			if backup.BackupSchedule != nil {
				trigger = fmt.Sprintf("%s (%s)", trigger, backup.BackupSchedule.Name)
			}

			row := []string{
				backup.Identity,
				clusterName,
				string(backup.EngineType),
				backup.EngineVersion,
				backup.BackupType,
				trigger,
				status,
				formattime.FormatTime(backup.CreatedAt.Local(), backupListShowExactTime),
			}

			if backup.StoppedAt != nil {
				row = append(row, formattime.FormatTime(backup.StoppedAt.Local(), backupListShowExactTime))
			} else {
				row = append(row, "-")
			}

			if backupListShowLabels {
				labelStrs := []string{}
				for k, v := range backup.Labels {
					labelStrs = append(labelStrs, k+"="+v)
				}
				sort.Strings(labelStrs)
				if len(labelStrs) == 0 {
					labelStrs = []string{"-"}
				}
				row = append(row, strings.Join(labelStrs, ","))
			}

			body = append(body, row)
		}

		if backupListNoHeader {
			table.Print(nil, body)
		} else {
			headers := []string{"ID", "Cluster", "Engine", "Version", "Type", "Trigger", "Status", "Created", "Completed"}
			if backupListShowLabels {
				headers = append(headers, "Labels")
			}
			table.Print(headers, body)
		}

		return nil
	},
}

func init() {
	BackupCmd.AddCommand(backupListCmd)

	backupListCmd.Flags().BoolVar(&backupListNoHeader, NoHeaderKey, false, "Do not print the header")
	backupListCmd.Flags().BoolVar(&backupListShowExactTime, "exact-time", false, "Show exact time instead of relative time")
	backupListCmd.Flags().BoolVar(&backupListShowLabels, "show-labels", false, "Show labels")
	backupListCmd.Flags().StringVarP(&backupListLabelSelector, "selector", "l", "", "Label selector to filter backups (format: key1=value1,key2=value2)")
	backupListCmd.Flags().StringVar(&backupListClusterFilter, "cluster", "", "Filter by database cluster identity, slug, or name")
	backupListCmd.Flags().StringVar(&backupListOlderThan, "older-than", "", "Filter backups older than the specified duration (e.g., 30d, 1w, 1mo, 1y, 24h)")
	backupListCmd.Flags().StringVar(&backupListNewerThan, "newer-than", "", "Filter backups newer than the specified duration (e.g., 7d, 1w, 1mo, 1y, 24h)")
	backupListCmd.Flags().StringSliceVar(&backupListStatusFilter, "status", []string{}, "Filter by backup status (can be specified multiple times, e.g., --status ready --status failed)")

	// Register completions
	backupListCmd.RegisterFlagCompletionFunc("cluster", completion.CompleteDbClusterID)
}
