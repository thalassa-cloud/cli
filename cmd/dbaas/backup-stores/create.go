package backupstores

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	dbaasutil "github.com/thalassa-cloud/cli/internal/dbaas"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/dbaas"
)

var (
	backupStoreCreateName             string
	backupStoreCreateDescription      string
	backupStoreCreateRegion           string
	backupStoreCreateRetentionPolicy  string
	backupStoreCreateRetentionMode    string
	backupStoreCreateDeleteProtection bool
	backupStoreCreateLabels           []string
	backupStoreCreateAnnotations      []string
	backupStoreCreateNoHeader         bool
)

var backupStoreCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a database backup store",
	Long:  "Create a backup store for Barman backups and point-in-time recovery",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		if backupStoreCreateName == "" {
			return fmt.Errorf("name is required")
		}
		if backupStoreCreateRegion == "" {
			return fmt.Errorf("region is required")
		}

		retentionMode, err := dbaasutil.ParseRetentionMode(backupStoreCreateRetentionMode)
		if err != nil {
			return err
		}

		req := dbaas.CreateDbObjectStoreRequest{
			Name:             backupStoreCreateName,
			Description:      backupStoreCreateDescription,
			Region:           backupStoreCreateRegion,
			RetentionPolicy:  backupStoreCreateRetentionPolicy,
			RetentionMode:    retentionMode,
			DeleteProtection: backupStoreCreateDeleteProtection,
			Labels:           shared.KeyValuePairsToMap(backupStoreCreateLabels),
			Annotations:      shared.KeyValuePairsToMap(backupStoreCreateAnnotations),
		}

		store, err := client.DBaaS().CreateDbObjectStore(cmd.Context(), req)
		if err != nil {
			return fmt.Errorf("failed to create backup store: %w", err)
		}

		regionName := backupStoreCreateRegion
		if store.Region != nil {
			regionName = store.Region.Name
		}

		body := [][]string{{
			store.Identity,
			store.Name,
			regionName,
			string(store.Status),
			store.RetentionPolicy,
			dbaasutil.FormatRetentionMode(store.RetentionMode),
			formattime.FormatTime(store.CreatedAt.Local(), false),
		}}

		if backupStoreCreateNoHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Region", "Status", "Retention", "Retention Mode", "Created"}, body)
		}

		return nil
	},
}

func init() {
	BackupStoresCmd.AddCommand(backupStoreCreateCmd)

	backupStoreCreateCmd.Flags().BoolVar(&backupStoreCreateNoHeader, NoHeaderKey, false, "Do not print the header")
	backupStoreCreateCmd.Flags().StringVar(&backupStoreCreateName, "name", "", "Name of the backup store (required)")
	backupStoreCreateCmd.Flags().StringVar(&backupStoreCreateDescription, "description", "", "Description of the backup store")
	backupStoreCreateCmd.Flags().StringVar(&backupStoreCreateRegion, "region", "", "Region identity or slug (required)")
	backupStoreCreateCmd.Flags().StringVar(&backupStoreCreateRetentionPolicy, "retention-policy", "", "Retention policy in days, for example 30d")
	backupStoreCreateCmd.Flags().StringVar(&backupStoreCreateRetentionMode, "retention-mode", string(dbaas.DbObjectStoreRetentionModeRetainForPointInTime), "Retention mode: retainForPointInTime or forceCleanupAfterExpiry")
	backupStoreCreateCmd.Flags().BoolVar(&backupStoreCreateDeleteProtection, "delete-protection", false, "Enable delete protection")
	backupStoreCreateCmd.Flags().StringSliceVar(&backupStoreCreateLabels, "labels", []string{}, "Labels in key=value format (can be specified multiple times)")
	backupStoreCreateCmd.Flags().StringSliceVar(&backupStoreCreateAnnotations, "annotations", []string{}, "Annotations in key=value format (can be specified multiple times)")

	_ = backupStoreCreateCmd.MarkFlagRequired("name")
	_ = backupStoreCreateCmd.MarkFlagRequired("region")
	_ = backupStoreCreateCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
	_ = backupStoreCreateCmd.RegisterFlagCompletionFunc("retention-mode", completeRetentionMode)
}

func completeRetentionMode(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{
		string(dbaas.DbObjectStoreRetentionModeRetainForPointInTime),
		string(dbaas.DbObjectStoreRetentionModeForceCleanupAfterExpiry),
	}, cobra.ShellCompDirectiveNoFileComp
}
