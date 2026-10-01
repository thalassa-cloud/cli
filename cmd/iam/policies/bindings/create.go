package bindings

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/iamresolve"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientiam "github.com/thalassa-cloud/client-go/iam"
)

var (
	createName                   string
	createDescription            string
	createLabels                 []string
	createAnnotations            []string
	createUserIdentity           string
	createServiceAccountIdentity string
	createExpiresAt              string
	createDuration               string
)

var createCmd = &cobra.Command{
	Use:               "create <policy>",
	Short:             "Create a binding to a user or service account",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteIAMPolicyIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		n := 0
		if createUserIdentity != "" {
			n++
		}
		if createServiceAccountIdentity != "" {
			n++
		}
		if n != 1 {
			return fmt.Errorf("specify exactly one of --user-identity or --service-account-identity")
		}
		if createName == "" {
			return fmt.Errorf("--name is required")
		}
		if createExpiresAt != "" && createDuration != "" {
			return fmt.Errorf("specify at most one of --expires-at or --duration")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		create := clientiam.CreateIamPolicyBindingRequest{
			Name:        createName,
			Description: createDescription,
			Labels:      shared.KeyValuePairsToMap(createLabels),
			Annotations: shared.KeyValuePairsToMap(createAnnotations),
			Duration:    createDuration,
		}
		if create.Labels == nil {
			create.Labels = map[string]string{}
		}
		if create.Annotations == nil {
			create.Annotations = map[string]string{}
		}
		if createUserIdentity != "" {
			create.UserIdentity = &createUserIdentity
		}
		if createServiceAccountIdentity != "" {
			create.ServiceAccountIdentity = &createServiceAccountIdentity
		}
		if createExpiresAt != "" {
			t, err := time.Parse(time.RFC3339, createExpiresAt)
			if err != nil {
				return fmt.Errorf("--expires-at must be RFC3339 (e.g. 2026-12-31T23:59:59Z): %w", err)
			}
			create.ExpiresAt = &t
		}

		policy, err := iamresolve.ResolveIamPolicyRef(ctx, client.IAM(), args[0])
		if err != nil {
			return err
		}
		binding, err := client.IAM().CreateIamPolicyBinding(ctx, policy.Identity, create)
		if err != nil {
			return fmt.Errorf("failed to create binding: %w", err)
		}
		if binding == nil {
			return nil
		}
		body := [][]string{{binding.Identity, binding.Name, binding.Slug}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Slug"}, body)
		}
		return nil
	},
}

func init() {
	BindingsCmd.AddCommand(createCmd)
	createCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	createCmd.Flags().StringVar(&createName, "name", "", "Binding name")
	createCmd.Flags().StringVar(&createDescription, "description", "", "Binding description")
	createCmd.Flags().StringSliceVar(&createLabels, "labels", nil, "Labels as key=value (repeatable)")
	createCmd.Flags().StringSliceVar(&createAnnotations, "annotations", nil, "Annotations as key=value (repeatable)")
	createCmd.Flags().StringVar(&createUserIdentity, "user-identity", "", "Bind to this user identity")
	createCmd.Flags().StringVar(&createServiceAccountIdentity, "service-account-identity", "", "Bind to this service account identity")
	createCmd.Flags().StringVar(&createExpiresAt, "expires-at", "", "Absolute expiry time (RFC3339)")
	createCmd.Flags().StringVar(&createDuration, "duration", "", "Relative TTL (e.g. 24h, 7d) instead of --expires-at")
	_ = createCmd.RegisterFlagCompletionFunc("user-identity", completion.CompleteIAMAppUserSubject)
	_ = createCmd.RegisterFlagCompletionFunc("service-account-identity", completion.CompleteIAMServiceAccountIdentityFlag)
	_ = createCmd.MarkFlagRequired("name")
}
