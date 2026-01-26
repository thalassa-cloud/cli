package backup

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	backupViewShowExactTime bool
	backupViewNoHeader      bool
)

// backupViewCmd represents the backup view command
var backupViewCmd = &cobra.Command{
	Use:   "view",
	Short: "View backup details",
	Long:  "View detailed information about a database backup",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		backupIdentity := args[0]

		backup, err := client.DBaaS().GetDbBackup(cmd.Context(), backupIdentity)
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("backup not found: %s", backupIdentity)
			}
			return fmt.Errorf("failed to get backup: %w", err)
		}

		clusterName := ""
		if backup.DbCluster != nil {
			clusterName = backup.DbCluster.Name
		}

		trigger := string(backup.BackupTrigger)
		if backup.BackupSchedule != nil {
			trigger = fmt.Sprintf("%s (%s)", trigger, backup.BackupSchedule.Name)
		}

		body := [][]string{
			{"ID", backup.Identity},
			{"Cluster", clusterName},
			{"Status", string(backup.Status)},
			{"Engine", string(backup.EngineType)},
			{"Engine Version", backup.EngineVersion},
			{"Backup Type", backup.BackupType},
			{"Trigger", trigger},
			{"Online", fmt.Sprintf("%v", backup.Online)},
			{"Delete Protection", fmt.Sprintf("%v", backup.DeleteProtection)},
			{"Created", formattime.FormatTime(backup.CreatedAt.Local(), backupViewShowExactTime)},
		}

		if backup.StartedAt != nil {
			body = append(body, []string{"Started", formattime.FormatTime(backup.StartedAt.Local(), backupViewShowExactTime)})
		}

		if backup.StoppedAt != nil {
			body = append(body, []string{"Stopped", formattime.FormatTime(backup.StoppedAt.Local(), backupViewShowExactTime)})
		}

		if backup.StatusMessage != "" {
			body = append(body, []string{"Status Message", backup.StatusMessage})
		}

		if backup.BeginLSN != "" {
			body = append(body, []string{"Begin LSN", backup.BeginLSN})
		}

		if backup.EndLSN != "" {
			body = append(body, []string{"End LSN", backup.EndLSN})
		}

		if backup.BeginWAL != "" {
			body = append(body, []string{"Begin WAL", backup.BeginWAL})
		}

		if backup.EndWAL != "" {
			body = append(body, []string{"End WAL", backup.EndWAL})
		}

		if backup.DeleteScheduledAt != nil {
			body = append(body, []string{"Delete Scheduled At", formattime.FormatTime(backup.DeleteScheduledAt.Local(), backupViewShowExactTime)})
		}

		if len(backup.Labels) > 0 {
			labelStrs := []string{}
			for k, v := range backup.Labels {
				labelStrs = append(labelStrs, k+"="+v)
			}
			sort.Strings(labelStrs)
			body = append(body, []string{"Labels", strings.Join(labelStrs, ", ")})
		}

		if len(backup.Annotations) > 0 {
			annotationStrs := []string{}
			for k, v := range backup.Annotations {
				annotationStrs = append(annotationStrs, k+"="+v)
			}
			sort.Strings(annotationStrs)
			body = append(body, []string{"Annotations", strings.Join(annotationStrs, ", ")})
		}

		if backupViewNoHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"Field", "Value"}, body)
		}

		return nil
	},
}

func init() {
	BackupCmd.AddCommand(backupViewCmd)

	backupViewCmd.Flags().BoolVar(&backupViewNoHeader, NoHeaderKey, false, "Do not print the header")
	backupViewCmd.Flags().BoolVar(&backupViewShowExactTime, "exact-time", false, "Show exact time instead of relative time")
}
