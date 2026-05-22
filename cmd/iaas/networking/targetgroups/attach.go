package targetgroups

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
)

var (
	attachServer   string
	attachEndpoint string
)

var attachCmd = &cobra.Command{
	Use:     "attach TARGET_GROUP",
	Short:   "Attach a target to a target group",
	Long:    "Attach a server or endpoint to a target group.",
	Example: "tcloud networking target-groups attach tg-123 --server machine-123\ntcloud networking target-groups attach tg-123 --endpoint endpoint-123",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if attachServer == "" && attachEndpoint == "" {
			return fmt.Errorf("either --server or --endpoint is required")
		}
		if attachServer != "" && attachEndpoint != "" {
			return fmt.Errorf("only one of --server or --endpoint may be set")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		attach := iaas.AttachTarget{}
		if attachServer != "" {
			attach.ServerIdentity = attachServer
		}
		if attachEndpoint != "" {
			attach.EndpointIdentity = attachEndpoint
		}

		attachment, err := client.IaaS().AttachServerToTargetGroup(cmd.Context(), iaas.AttachTargetGroupRequest{
			TargetGroupID: args[0],
			AttachTarget:  attach,
		})
		if err != nil {
			return err
		}

		fmt.Printf("Target attached successfully\n")
		fmt.Printf("Attachment ID: %s\n", attachment.Identity)
		return nil
	},
}

func init() {
	TargetGroupsCmd.AddCommand(attachCmd)

	attachCmd.Flags().StringVar(&attachServer, "server", "", "Server identity to attach")
	attachCmd.Flags().StringVar(&attachEndpoint, "endpoint", "", "Endpoint identity to attach")

	attachCmd.ValidArgsFunction = completeTargetGroupID
	attachCmd.RegisterFlagCompletionFunc("server", completeMachineID)
}
