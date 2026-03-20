package members

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/internal/shared"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientiam "github.com/thalassa-cloud/client-go/iam"
)

var updateRole string

var updateCmd = &cobra.Command{
	Use:               "update <member>",
	Short:             "Change an organisation member's role (OWNER or MEMBER)",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteIAMOrganisationMemberIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if updateRole == "" {
			return fmt.Errorf("--role is required (OWNER or MEMBER)")
		}
		role := clientiam.OrganisationMemberType(updateRole)
		if role != clientiam.OrganisationMemberTypeOwner && role != clientiam.OrganisationMemberTypeMember {
			return fmt.Errorf("role must be OWNER or MEMBER")
		}
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		if err := client.IAM().UpdateOrganisationMember(ctx, args[0], clientiam.UpdateOrganisationMemberRequest{
			MemberType: role,
		}); err != nil {
			return fmt.Errorf("failed to update member: %w", err)
		}
		fmt.Printf("Updated member %s to role %s\n", args[0], role)
		return nil
	},
}

func init() {
	MembersCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	updateCmd.Flags().StringVar(&updateRole, "role", "", "Organisation role: OWNER or MEMBER")
	_ = updateCmd.RegisterFlagCompletionFunc("role", completion.CompleteIAMOrganisationMemberType)
	_ = updateCmd.MarkFlagRequired("role")
}
