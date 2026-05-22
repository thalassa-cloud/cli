package targetgroups

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var detachCmd = &cobra.Command{
	Use:     "detach TARGET_GROUP ATTACHMENT",
	Short:   "Detach a target from a target group",
	Long:    "Detach a target attachment from a target group.",
	Example: "tcloud networking target-groups detach tg-123 attachment-456",
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		if err := client.IaaS().DetachServerFromTargetGroup(cmd.Context(), iaas.DetachTargetRequest{
			TargetGroupID: args[0],
			AttachmentID:  args[1],
		}); err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("attachment not found: %s", args[1])
			}
			return err
		}

		fmt.Printf("Attachment %s detached successfully\n", args[1])
		return nil
	},
}

func init() {
	TargetGroupsCmd.AddCommand(detachCmd)
	detachCmd.ValidArgsFunction = completeTargetGroupID
}
