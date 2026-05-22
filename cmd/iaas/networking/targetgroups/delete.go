package targetgroups

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/labels"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/iaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	deleteForce         bool
	deleteLabelSelector string
)

var deleteCmd = &cobra.Command{
	Use:     "delete [TARGET_GROUP...]",
	Short:   "Delete target group(s)",
	Long:    "Delete target group(s) by identity or label selector.",
	Example: "tcloud networking target-groups delete tg-123\ntcloud networking target-groups delete --selector env=test --force",
	Aliases: []string{"d", "del", "remove", "rm"},
	Args:    cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && deleteLabelSelector == "" {
			return fmt.Errorf("either target group identity(ies) or --selector must be provided")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		targetGroupsToDelete := []string{}

		if deleteLabelSelector != "" {
			all, err := client.IaaS().ListTargetGroups(cmd.Context(), &iaas.ListTargetGroupsRequest{
				Filters: []filters.Filter{
					&filters.LabelFilter{MatchLabels: labels.ParseLabelSelector(deleteLabelSelector)},
				},
			})
			if err != nil {
				return fmt.Errorf("failed to list target groups: %w", err)
			}
			if len(all) == 0 {
				fmt.Println("No target groups found matching the label selector")
				return nil
			}
			for _, tg := range all {
				targetGroupsToDelete = append(targetGroupsToDelete, tg.Identity)
			}
		} else {
			targetGroupsToDelete = append(targetGroupsToDelete, args...)
		}

		if !deleteForce {
			fmt.Printf("Are you sure you want to delete %d target group(s)?\n", len(targetGroupsToDelete))
			var confirm string
			fmt.Printf("Enter 'yes' to confirm: ")
			fmt.Scanln(&confirm)
			if confirm != "yes" {
				fmt.Println("Aborted")
				return nil
			}
		}

		for _, tgID := range targetGroupsToDelete {
			fmt.Printf("Deleting target group: %s\n", tgID)
			if err := client.IaaS().DeleteTargetGroup(cmd.Context(), iaas.DeleteTargetGroupRequest{
				Identity: tgID,
			}); err != nil {
				if tcclient.IsNotFound(err) {
					fmt.Printf("Target group %s not found\n", tgID)
					continue
				}
				return fmt.Errorf("failed to delete target group: %w", err)
			}
			fmt.Printf("Target group %s deleted successfully\n", tgID)
		}

		return nil
	},
}

func init() {
	TargetGroupsCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVar(&deleteForce, "force", false, "Skip confirmation")
	deleteCmd.Flags().StringVarP(&deleteLabelSelector, "selector", "l", "", "Label selector (format: key1=value1,key2=value2)")

	deleteCmd.ValidArgsFunction = completeTargetGroupID
}
