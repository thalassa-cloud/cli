package federatedidentityproviders

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/internal/shared"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientiam "github.com/thalassa-cloud/client-go/iam"
)

var (
	updateName        string
	updateDescription string
	updateLabels      []string
	updateAnnotations []string
	updateJwksURI     string
	updateStatus      string
)

var updateCmd = &cobra.Command{
	Use:               "update <identity>",
	Short:             "Update a federated identity provider",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteIAMFederatedIdentityProviderIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		req := clientiam.UpdateFederatedIdentityProviderRequest{}
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
		if cmd.Flags().Changed("jwks-uri") {
			if updateJwksURI == "" {
				empty := ""
				req.ProviderJwksURI = &empty
			} else {
				req.ProviderJwksURI = &updateJwksURI
			}
		}
		if cmd.Flags().Changed("status") {
			s := clientiam.FederatedIdentityProviderStatus(updateStatus)
			if s != clientiam.FederatedIdentityProviderStatusActive && s != clientiam.FederatedIdentityProviderStatusInactive {
				return fmt.Errorf("invalid --status")
			}
			req.Status = s
		}
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		p, err := client.IAM().UpdateFederatedIdentityProvider(ctx, args[0], req)
		if err != nil {
			return fmt.Errorf("failed to update provider: %w", err)
		}
		if p == nil {
			return nil
		}
		body := [][]string{{p.Identity, p.Name, string(p.Status)}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Status"}, body)
		}
		return nil
	},
}

func init() {
	FederatedIdentityProvidersCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	updateCmd.Flags().StringVar(&updateName, "name", "", "Provider display name")
	updateCmd.Flags().StringVar(&updateDescription, "description", "", "Description")
	updateCmd.Flags().StringSliceVar(&updateLabels, "labels", nil, "Replace labels (key=value, repeatable)")
	updateCmd.Flags().StringSliceVar(&updateAnnotations, "annotations", nil, "Replace annotations (key=value, repeatable)")
	updateCmd.Flags().StringVar(&updateJwksURI, "jwks-uri", "", "JWKS URI (set to empty string to clear)")
	updateCmd.Flags().StringVar(&updateStatus, "status", "", "active or inactive")
}
