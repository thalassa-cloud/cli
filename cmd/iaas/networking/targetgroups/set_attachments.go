package targetgroups

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
)

var (
	setAttachmentServers   []string
	setAttachmentEndpoints []string
)

var setAttachmentsCmd = &cobra.Command{
	Use:     "set-attachments TARGET_GROUP",
	Short:   "Replace target group attachments",
	Long:    "Replace all server and endpoint attachments on a target group. Existing attachments not listed are removed.",
	Example: "tcloud networking target-groups set-attachments tg-123 --server machine-1 --server machine-2",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(setAttachmentServers) == 0 && len(setAttachmentEndpoints) == 0 {
			return fmt.Errorf("at least one --server or --endpoint is required")
		}

		attachments := make([]iaas.AttachTarget, 0, len(setAttachmentServers)+len(setAttachmentEndpoints))
		for _, server := range setAttachmentServers {
			attachments = append(attachments, iaas.AttachTarget{ServerIdentity: server})
		}
		for _, endpoint := range setAttachmentEndpoints {
			attachments = append(attachments, iaas.AttachTarget{EndpointIdentity: endpoint})
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		if err := client.IaaS().SetTargetGroupServerAttachments(cmd.Context(), iaas.TargetGroupAttachmentsBatch{
			TargetGroupID: args[0],
			Attachments:   attachments,
		}); err != nil {
			return err
		}

		fmt.Printf("Target group %s attachments updated successfully\n", args[0])
		return nil
	},
}

func init() {
	TargetGroupsCmd.AddCommand(setAttachmentsCmd)

	setAttachmentsCmd.Flags().StringSliceVar(&setAttachmentServers, "server", []string{}, "Server identities to attach")
	setAttachmentsCmd.Flags().StringSliceVar(&setAttachmentEndpoints, "endpoint", []string{}, "Endpoint identities to attach")

	setAttachmentsCmd.ValidArgsFunction = completeTargetGroupID
	setAttachmentsCmd.RegisterFlagCompletionFunc("server", completeMachineID)
}
