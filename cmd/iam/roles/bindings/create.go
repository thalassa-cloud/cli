package bindings

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/iam/internal/shared"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/iamresolve"
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
	createTeamIdentity           string
	createServiceAccountIdentity string
	createScopes                 []string
)

var createCmd = &cobra.Command{
	Use:   "create <role>",
	Short: "Create a binding to a user, team, or service account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		n := 0
		if createUserIdentity != "" {
			n++
		}
		if createTeamIdentity != "" {
			n++
		}
		if createServiceAccountIdentity != "" {
			n++
		}
		if n != 1 {
			return fmt.Errorf("specify exactly one of --user-identity, --team-identity, or --service-account-identity")
		}
		if createName == "" {
			return fmt.Errorf("--name is required")
		}
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		create := clientiam.CreateRoleBinding{
			Name:        createName,
			Description: createDescription,
			Labels:      shared.KeyValuePairsToMap(createLabels),
			Annotations: shared.KeyValuePairsToMap(createAnnotations),
			Scopes:      createScopes,
		}
		if createUserIdentity != "" {
			create.UserIdentity = &createUserIdentity
		}
		if createTeamIdentity != "" {
			create.TeamIdentity = &createTeamIdentity
		}
		if createServiceAccountIdentity != "" {
			create.ServiceAccountIdentity = &createServiceAccountIdentity
		}
		role, err := iamresolve.ResolveOrganisationRoleRef(ctx, client.IAM(), args[0])
		if err != nil {
			return err
		}
		binding, err := client.IAM().CreateRoleBinding(ctx, role.Identity, create)
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
	createCmd.Flags().StringVar(&createTeamIdentity, "team-identity", "", "Bind to this team identity")
	createCmd.Flags().StringVar(&createServiceAccountIdentity, "service-account-identity", "", "Bind to this service account identity")
	_ = createCmd.RegisterFlagCompletionFunc("user-identity", completion.CompleteIAMAppUserSubject)
	_ = createCmd.RegisterFlagCompletionFunc("team-identity", completion.CompleteIAMTeamIdentityFlag)
	_ = createCmd.RegisterFlagCompletionFunc("service-account-identity", completion.CompleteIAMServiceAccountIdentityFlag)
	createCmd.Flags().StringSliceVar(&createScopes, "scope", nil, "Scopes for the binding (repeatable)")
	_ = createCmd.MarkFlagRequired("name")
}
