package securitygroups

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	updateName           string
	updateDescription    string
	updateLabels         []string
	updateAnnotations    []string
	updateAllowSameGroup bool
)

var updateCmd = &cobra.Command{
	Use:               "update SECURITY_GROUP",
	Short:             "Update security group metadata",
	Long:              "Update name, description, labels, annotations, or allow-same-group. Rules are left unchanged.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completeSecurityGroupID,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		current, err := client.IaaS().GetSecurityGroup(cmd.Context(), args[0])
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("security group not found: %s", args[0])
			}
			return fmt.Errorf("failed to get security group: %w", err)
		}

		req := iaas.UpdateSecurityGroupRequest{
			Name:                  current.Name,
			Description:           current.Description,
			Labels:                current.Labels,
			Annotations:           current.Annotations,
			ObjectVersion:         current.ObjectVersion,
			AllowSameGroupTraffic: current.AllowSameGroupTraffic,
			SkipRulesUpdate:       true,
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
		if cmd.Flags().Changed("allow-same-group") {
			req.AllowSameGroupTraffic = updateAllowSameGroup
		}

		sg, err := client.IaaS().UpdateSecurityGroup(cmd.Context(), current.Identity, req)
		if err != nil {
			return fmt.Errorf("failed to update security group: %w", err)
		}
		fmt.Printf("Security group %s updated\n", sg.Identity)
		return nil
	},
}

func init() {
	SecurityGroupsCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVar(&updateName, "name", "", "Name")
	updateCmd.Flags().StringVar(&updateDescription, "description", "", "Description")
	updateCmd.Flags().StringSliceVar(&updateLabels, "labels", nil, "Labels as key=value (repeatable)")
	updateCmd.Flags().StringSliceVar(&updateAnnotations, "annotations", nil, "Annotations as key=value (repeatable)")
	updateCmd.Flags().BoolVar(&updateAllowSameGroup, "allow-same-group", false, "Allow traffic between instances in the same security group")
}
