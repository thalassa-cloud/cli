package backupstores

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	dbaasutil "github.com/thalassa-cloud/cli/internal/dbaas"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/dbaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	backupStoreUpdateName             string
	backupStoreUpdateDescription      string
	backupStoreUpdateRetentionPolicy  string
	backupStoreUpdateRetentionMode    string
	backupStoreUpdateDeleteProtection bool
	backupStoreUpdateLabels           []string
	backupStoreUpdateAnnotations      []string
	backupStoreUpdateNoHeader         bool
)

var backupStoreUpdateCmd = &cobra.Command{
	Use:               "update",
	Short:             "Update a database backup store",
	Long:              "Update properties of an existing backup store",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteDbBackupStoreID,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		storeIdentity := args[0]
		current, err := client.DBaaS().GetDbObjectStore(cmd.Context(), storeIdentity)
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("backup store not found: %s", storeIdentity)
			}
			return fmt.Errorf("failed to get backup store: %w", err)
		}

		req := dbaas.UpdateDbObjectStoreRequest{
			Name:             current.Name,
			Description:      current.Description,
			RetentionPolicy:  current.RetentionPolicy,
			DeleteProtection: current.DeleteProtection,
			Labels:           current.Labels,
			Annotations:      current.Annotations,
		}

		if cmd.Flags().Changed("name") {
			req.Name = backupStoreUpdateName
		}
		if cmd.Flags().Changed("description") {
			req.Description = backupStoreUpdateDescription
		}
		if cmd.Flags().Changed("retention-policy") {
			req.RetentionPolicy = backupStoreUpdateRetentionPolicy
		}
		if cmd.Flags().Changed("retention-mode") {
			mode, err := dbaasutil.ParseRetentionMode(backupStoreUpdateRetentionMode)
			if err != nil {
				return err
			}
			req.RetentionMode = &mode
		}
		if cmd.Flags().Changed("delete-protection") {
			req.DeleteProtection = backupStoreUpdateDeleteProtection
		}
		if cmd.Flags().Changed("labels") {
			req.Labels = shared.KeyValuePairsToMap(backupStoreUpdateLabels)
		}
		if cmd.Flags().Changed("annotations") {
			req.Annotations = shared.KeyValuePairsToMap(backupStoreUpdateAnnotations)
		}

		store, err := client.DBaaS().UpdateDbObjectStore(cmd.Context(), storeIdentity, req)
		if err != nil {
			return fmt.Errorf("failed to update backup store: %w", err)
		}

		regionName := ""
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
		}}

		if backupStoreUpdateNoHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Region", "Status", "Retention", "Retention Mode"}, body)
		}

		return nil
	},
}

func init() {
	BackupStoresCmd.AddCommand(backupStoreUpdateCmd)

	backupStoreUpdateCmd.Flags().BoolVar(&backupStoreUpdateNoHeader, NoHeaderKey, false, "Do not print the header")
	backupStoreUpdateCmd.Flags().StringVar(&backupStoreUpdateName, "name", "", "Name of the backup store")
	backupStoreUpdateCmd.Flags().StringVar(&backupStoreUpdateDescription, "description", "", "Description of the backup store")
	backupStoreUpdateCmd.Flags().StringVar(&backupStoreUpdateRetentionPolicy, "retention-policy", "", "Retention policy in days, for example 30d")
	backupStoreUpdateCmd.Flags().StringVar(&backupStoreUpdateRetentionMode, "retention-mode", "", "Retention mode: retainForPointInTime or forceCleanupAfterExpiry")
	backupStoreUpdateCmd.Flags().BoolVar(&backupStoreUpdateDeleteProtection, "delete-protection", false, "Enable or disable delete protection")
	backupStoreUpdateCmd.Flags().StringSliceVar(&backupStoreUpdateLabels, "labels", []string{}, "Labels in key=value format (can be specified multiple times)")
	backupStoreUpdateCmd.Flags().StringSliceVar(&backupStoreUpdateAnnotations, "annotations", []string{}, "Annotations in key=value format (can be specified multiple times)")

	_ = backupStoreUpdateCmd.RegisterFlagCompletionFunc("retention-mode", completeRetentionMode)
}
