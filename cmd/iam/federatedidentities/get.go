package federatedidentities

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/internal/shared"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var getCmd = &cobra.Command{
	Use:               "get <identity>",
	Short:             "Show a federated identity",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteIAMFederatedIdentityIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		fi, err := client.IAM().GetFederatedIdentity(ctx, args[0])
		if err != nil {
			return fmt.Errorf("failed to get federated identity: %w", err)
		}
		fmt.Printf("Identity:          %s\n", fi.Identity)
		fmt.Printf("Name:              %s\n", fi.Name)
		fmt.Printf("Description:       %s\n", fi.Description)
		fmt.Printf("Provider subject:  %s\n", fi.ProviderSubject)
		fmt.Printf("Status:            %s\n", fi.Status)
		fmt.Printf("Audience mode:     %s\n", fi.AudienceMatchMode)
		fmt.Printf("Trusted audiences: %s\n", strings.Join(fi.TrustedAudiences, ","))
		scopes := make([]string, 0, len(fi.AllowedScopes))
		for _, s := range fi.AllowedScopes {
			scopes = append(scopes, string(s))
		}
		fmt.Printf("Allowed scopes:    %s\n", strings.Join(scopes, ","))
		if fi.Provider != nil {
			fmt.Printf("Provider:          %s (%s)\n", fi.Provider.Name, fi.Provider.Identity)
		}
		if fi.ServiceAccount != nil {
			fmt.Printf("Service account:   %s (%s)\n", fi.ServiceAccount.Name, fi.ServiceAccount.Identity)
		}
		if fi.ExpiresAt != nil {
			fmt.Printf("Expires at:        %s\n", fi.ExpiresAt.Format(time.RFC3339))
		}
		fmt.Printf("Created:           %s\n", formattime.FormatTime(fi.CreatedAt.Local(), showExactTime))
		return nil
	},
}

func init() {
	FederatedIdentitiesCmd.AddCommand(getCmd)
	getCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	getCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
}
