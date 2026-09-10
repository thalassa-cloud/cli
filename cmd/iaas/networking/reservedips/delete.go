package reservedips

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/labels"
	"github.com/thalassa-cloud/cli/internal/shared"
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
	Use:     "delete",
	Short:   "Delete reserved IP address(es)",
	Long:    "Delete reserved IP address(es) by identity or label selector. Attached reserved IPs are disassociated before deletion.",
	Example: "tcloud networking reserved-ips delete rip-123 --force\ntcloud networking reserved-ips delete --selector env=test --force",
	Aliases: []string{"d", "del", "remove", "rm"},
	Args:    cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && deleteLabelSelector == "" {
			return fmt.Errorf("either reserved IP identity(ies) or --selector must be provided")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		toDelete := []iaas.ReservedIP{}
		if deleteLabelSelector != "" {
			all, err := client.IaaS().ListReservedIPs(cmd.Context(), &iaas.ListReservedIPsRequest{
				Filters: []filters.Filter{
					&filters.LabelFilter{MatchLabels: labels.ParseLabelSelector(deleteLabelSelector)},
				},
			})
			if err != nil {
				return fmt.Errorf("failed to list reserved IPs: %w", err)
			}
			if len(all) == 0 {
				fmt.Println("No reserved IPs found matching the label selector")
				return nil
			}
			toDelete = append(toDelete, all...)
		} else {
			for _, identity := range args {
				rip, err := client.IaaS().GetReservedIP(cmd.Context(), identity)
				if err != nil {
					if tcclient.IsNotFound(err) {
						fmt.Printf("Reserved IP %s not found\n", identity)
						continue
					}
					return fmt.Errorf("failed to get reserved IP: %w", err)
				}
				toDelete = append(toDelete, *rip)
			}
		}

		if len(toDelete) == 0 {
			fmt.Println("No reserved IPs to delete")
			return nil
		}

		var summary strings.Builder
		fmt.Fprintf(&summary, "Are you sure you want to delete the following reserved IP(s)?\n")
		for _, rip := range toDelete {
			fmt.Fprintf(&summary, "  %s (%s)\n", rip.Name, rip.Identity)
		}
		proceed, err := shared.PromptDestructiveUnlessForce(deleteForce, summary.String())
		if err != nil {
			return err
		}
		if !proceed {
			return nil
		}

		for _, rip := range toDelete {
			fmt.Printf("Deleting reserved IP: %s (%s)\n", rip.Name, rip.Identity)
			if err := client.IaaS().DeleteReservedIP(cmd.Context(), rip.Identity); err != nil {
				return fmt.Errorf("failed to delete reserved IP: %w", err)
			}
			fmt.Printf("Reserved IP %s deleted successfully\n", rip.Identity)
		}
		return nil
	},
}

func init() {
	ReservedIPsCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVar(&deleteForce, "force", false, "Force the deletion and skip the confirmation")
	deleteCmd.Flags().StringVarP(&deleteLabelSelector, "selector", "l", "", "Label selector to filter reserved IPs (format: key1=value1,key2=value2)")
	deleteCmd.ValidArgsFunction = completeReservedIPID
}
