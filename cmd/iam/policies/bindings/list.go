package bindings

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/iamresolve"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientiam "github.com/thalassa-cloud/client-go/iam"
)

var (
	noHeader      bool
	showExactTime bool
)

var listCmd = &cobra.Command{
	Use:               "list <policy>",
	Short:             "List bindings for a policy",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteIAMPolicyIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		policy, err := iamresolve.ResolveIamPolicyRef(ctx, client.IAM(), args[0])
		if err != nil {
			return err
		}

		bindings, err := client.IAM().ListIamPolicyBindings(ctx, policy.Identity, &clientiam.ListIamPolicyBindingsRequest{})
		if err != nil {
			return fmt.Errorf("failed to list bindings: %w", err)
		}
		body := make([][]string, 0, len(bindings))
		for _, b := range bindings {
			subject := ""
			switch {
			case b.User != nil:
				subject = "user:" + shared.UserPtrDisplay(b.User)
			case b.ServiceAccount != nil:
				subject = "service_account:" + b.ServiceAccount.Slug
			}
			expires := ""
			if b.ExpiresAt != nil {
				expires = formattime.FormatTime(b.ExpiresAt.Local(), showExactTime)
			}
			body = append(body, []string{b.Identity, b.Name, subject, expires})
		}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Subject", "Expires"}, body)
		}
		return nil
	},
}

func init() {
	BindingsCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	listCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
}
