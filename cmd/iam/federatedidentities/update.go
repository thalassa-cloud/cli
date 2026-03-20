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
	updateDescription    string
	updateLabels         []string
	updateAnnotations    []string
	updateAudiences      []string
	updateAudienceMode   string
	updateScopes         []string
	updateStatus         string
	updateExpiresAt      string
	updateConditions     string
	updateConditionsFile string
)

var updateCmd = &cobra.Command{
	Use:               "update <identity>",
	Short:             "Update a federated identity (only set flags are sent)",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteIAMFederatedIdentityIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		req := clientiam.UpdateFederatedIdentityRequest{}
		if cmd.Flags().Changed("description") {
			req.Description = updateDescription
		}
		if cmd.Flags().Changed("labels") {
			req.Labels = shared.KeyValuePairsToMap(updateLabels)
		}
		if cmd.Flags().Changed("annotations") {
			req.Annotations = shared.KeyValuePairsToMap(updateAnnotations)
		}
		if cmd.Flags().Changed("trusted-audience") {
			req.TrustedAudiences = updateAudiences
		}
		if cmd.Flags().Changed("audience-match-mode") {
			m := clientiam.AudienceMatchMode(updateAudienceMode)
			if m != clientiam.AudienceMatchModeExact && m != clientiam.AudienceMatchModeAny && m != clientiam.AudienceMatchModeAll {
				return fmt.Errorf("invalid --audience-match-mode")
			}
			req.AudienceMatchMode = m
		}
		if cmd.Flags().Changed("scope") {
			scopes, err := shared.ParseAccessCredentialScopes(updateScopes)
			if err != nil {
				return err
			}
			req.AllowedScopes = scopes
		}
		if cmd.Flags().Changed("status") {
			s := clientiam.FederatedIdentityStatus(updateStatus)
			if s != clientiam.FederatedIdentityStatusActive && s != clientiam.FederatedIdentityStatusInactive &&
				s != clientiam.FederatedIdentityStatusExpired && s != clientiam.FederatedIdentityStatusRevoked {
				return fmt.Errorf("invalid --status")
			}
			req.Status = s
		}
		if cmd.Flags().Changed("expires-at") {
			t, err := shared.ParseOptionalRFC3339(updateExpiresAt)
			if err != nil {
				return err
			}
			req.ExpiresAt = t
		}
		if cmd.Flags().Changed("conditions") || cmd.Flags().Changed("conditions-file") {
			conds, err := shared.ParseConditionsJSON(updateConditions, updateConditionsFile)
			if err != nil {
				return err
			}
			req.Conditions = conds
		}
		fi, err := client.IAM().UpdateFederatedIdentity(ctx, args[0], req)
		if err != nil {
			return fmt.Errorf("failed to update federated identity: %w", err)
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
	FederatedIdentitiesCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	updateCmd.Flags().StringVar(&updateDescription, "description", "", "Description")
	updateCmd.Flags().StringSliceVar(&updateLabels, "labels", nil, "Replace labels (key=value, repeatable)")
	updateCmd.Flags().StringSliceVar(&updateAnnotations, "annotations", nil, "Replace annotations (key=value, repeatable)")
	updateCmd.Flags().StringSliceVar(&updateAudiences, "trusted-audience", nil, "Replace trusted audiences (repeatable)")
	updateCmd.Flags().StringVar(&updateAudienceMode, "audience-match-mode", "", "exact, any, or all")
	updateCmd.Flags().StringSliceVar(&updateScopes, "scope", nil, "Replace allowed scopes (repeatable)")
	updateCmd.Flags().StringVar(&updateStatus, "status", "", "active, inactive, expired, or revoked")
	updateCmd.Flags().StringVar(&updateExpiresAt, "expires-at", "", "RFC3339 expiry (empty to clear not supported by all APIs)")
	updateCmd.Flags().StringVar(&updateConditions, "conditions", "", "Conditions JSON object")
	updateCmd.Flags().StringVar(&updateConditionsFile, "conditions-file", "", "Path to JSON file for conditions")
}
