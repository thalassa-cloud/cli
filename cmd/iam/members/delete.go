package members

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/internal/shared"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var deleteForce bool

var deleteCmd = &cobra.Command{
	Use:               "delete <member>",
	Short:             "Remove a member from the organisation",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteIAMOrganisationMemberIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		ok, err := shared.PromptDestructiveUnlessForce(deleteForce, fmt.Sprintf("Are you sure you want to remove this member from the organisation?\n  Member: %s\n", args[0]))
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		if err := client.IAM().DeleteOrganisationMember(ctx, args[0]); err != nil {
			return fmt.Errorf("failed to remove member: %w", err)
		}
		fmt.Printf("Removed organisation member %s\n", args[0])
		return nil
	},
}

func init() {
	MembersCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVar(&deleteForce, shared.ForceKey, false, "Skip the confirmation prompt and remove the member")
	deleteCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
}
