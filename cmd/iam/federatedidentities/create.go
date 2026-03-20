package federatedidentities

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
	createName           string
	createDescription    string
	createLabels         []string
	createAnnotations    []string
	createSA             string
	createProvider       string
	createSubject        string
	createAudiences      []string
	createAudienceMode   string
	createScopes         []string
	createExpiresAt      string
	createConditions     string
	createConditionsFile string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a federated identity",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		if createName == "" || createSA == "" || createProvider == "" || createSubject == "" {
			return fmt.Errorf("--name, --service-account, --provider, and --subject are required")
		}
		scopes, err := shared.ParseAccessCredentialScopes(createScopes)
		if err != nil {
			return err
		}
		mode := clientiam.AudienceMatchMode(createAudienceMode)
		if createAudienceMode == "" {
			mode = clientiam.AudienceMatchModeAny
		} else if mode != clientiam.AudienceMatchModeExact && mode != clientiam.AudienceMatchModeAny && mode != clientiam.AudienceMatchModeAll {
			return fmt.Errorf("invalid --audience-match-mode (use exact, any, or all)")
		}
		expires, err := shared.ParseOptionalRFC3339(createExpiresAt)
		if err != nil {
			return err
		}
		conds, err := shared.ParseConditionsJSON(createConditions, createConditionsFile)
		if err != nil {
			return err
		}
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		fi, err := client.IAM().CreateFederatedIdentity(ctx, clientiam.CreateFederatedIdentityRequest{
			Name:                   createName,
			Description:            createDescription,
			Labels:                 shared.KeyValuePairsToMap(createLabels),
			Annotations:            shared.KeyValuePairsToMap(createAnnotations),
			ServiceAccountIdentity: createSA,
			ProviderIdentity:       createProvider,
			ProviderSubject:        createSubject,
			TrustedAudiences:       createAudiences,
			AudienceMatchMode:      mode,
			AllowedScopes:          scopes,
			ExpiresAt:              expires,
			Conditions:             conds,
		})
		if err != nil {
			return fmt.Errorf("failed to create federated identity: %w", err)
		}
		if fi == nil {
			return nil
		}
		body := [][]string{{fi.Identity, fi.Name, string(fi.Status)}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Status"}, body)
		}
		return nil
	},
}

func init() {
	FederatedIdentitiesCmd.AddCommand(createCmd)
	createCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	createCmd.Flags().StringVar(&createName, "name", "", "Display name")
	createCmd.Flags().StringVar(&createDescription, "description", "", "Description")
	createCmd.Flags().StringSliceVar(&createLabels, "labels", nil, "Labels as key=value (repeatable)")
	createCmd.Flags().StringSliceVar(&createAnnotations, "annotations", nil, "Annotations as key=value (repeatable)")
	createCmd.Flags().StringVar(&createSA, "service-account", "", "Service account identity to bind")
	createCmd.Flags().StringVar(&createProvider, "provider", "", "Federated identity provider identity")
	createCmd.Flags().StringVar(&createSubject, "subject", "", "OIDC sub claim for this identity")
	createCmd.Flags().StringSliceVar(&createAudiences, "trusted-audience", nil, "Trusted JWT audiences (repeatable)")
	createCmd.Flags().StringVar(&createAudienceMode, "audience-match-mode", "", "exact, any (default), or all")
	createCmd.Flags().StringSliceVar(&createScopes, "scope", nil, "Allowed scopes: api:read, api:write, kubernetes, objectStorage (repeatable)")
	createCmd.Flags().StringVar(&createExpiresAt, "expires-at", "", "RFC3339 expiry time")
	createCmd.Flags().StringVar(&createConditions, "conditions", "", "Conditions as JSON object")
	createCmd.Flags().StringVar(&createConditionsFile, "conditions-file", "", "Path to JSON file for conditions")
	_ = createCmd.RegisterFlagCompletionFunc("service-account", completion.CompleteIAMServiceAccountIdentityFlag)
	_ = createCmd.RegisterFlagCompletionFunc("provider", completion.CompleteIAMFederatedIdentityProviderIdentityFlag)
}
