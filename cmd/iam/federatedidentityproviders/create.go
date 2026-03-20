package federatedidentityproviders

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientiam "github.com/thalassa-cloud/client-go/iam"
)

var (
	createName        string
	createDescription string
	createLabels      []string
	createAnnotations []string
	createIssuer      string
	createJwksURI     string
	createStatus      string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Register a federated identity provider",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		if createName == "" || createIssuer == "" {
			return fmt.Errorf("--name and --issuer are required")
		}
		var jwks *string
		if createJwksURI != "" {
			jwks = &createJwksURI
		}
		status := clientiam.FederatedIdentityProviderStatus(createStatus)
		if createStatus == "" {
			status = clientiam.FederatedIdentityProviderStatusActive
		} else if status != clientiam.FederatedIdentityProviderStatusActive && status != clientiam.FederatedIdentityProviderStatusInactive {
			return fmt.Errorf("invalid --status (use active or inactive)")
		}
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		p, err := client.IAM().CreateFederatedIdentityProvider(ctx, clientiam.CreateFederatedIdentityProviderRequest{
			Name:            createName,
			Description:     createDescription,
			Labels:          shared.KeyValuePairsToMap(createLabels),
			Annotations:     shared.KeyValuePairsToMap(createAnnotations),
			ProviderIssuer:  createIssuer,
			ProviderJwksURI: jwks,
			Status:          status,
		})
		if err != nil {
			return fmt.Errorf("failed to create provider: %w", err)
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
	FederatedIdentityProvidersCmd.AddCommand(createCmd)
	createCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	createCmd.Flags().StringVar(&createName, "name", "", "Provider name")
	createCmd.Flags().StringVar(&createDescription, "description", "", "Description")
	createCmd.Flags().StringSliceVar(&createLabels, "labels", nil, "Labels as key=value (repeatable)")
	createCmd.Flags().StringSliceVar(&createAnnotations, "annotations", nil, "Annotations as key=value (repeatable)")
	createCmd.Flags().StringVar(&createIssuer, "issuer", "", "OIDC issuer URL (unique per organisation)")
	createCmd.Flags().StringVar(&createJwksURI, "jwks-uri", "", "Optional JWKS URI override")
	createCmd.Flags().StringVar(&createStatus, "status", "", "active (default) or inactive")
}
