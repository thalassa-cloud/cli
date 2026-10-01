package policies

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientiam "github.com/thalassa-cloud/client-go/iam"
)

var (
	createName                string
	createDescription         string
	createLabels              []string
	createAnnotations         []string
	createReplicateToChildren bool
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an IAM policy",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		if createName == "" {
			return fmt.Errorf("--name is required")
		}
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		policy, err := client.IAM().CreateIamPolicy(ctx, clientiam.CreateIamPolicyRequest{
			Name:                createName,
			Description:         createDescription,
			Labels:              shared.KeyValuePairsToMap(createLabels),
			Annotations:         shared.KeyValuePairsToMap(createAnnotations),
			ReplicateToChildren: createReplicateToChildren,
		})
		if err != nil {
			return fmt.Errorf("failed to create IAM policy: %w", err)
		}
		body := [][]string{{policy.Identity, policy.Name, policy.Slug}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Slug"}, body)
		}
		return nil
	},
}

func init() {
	PoliciesCmd.AddCommand(createCmd)
	createCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	createCmd.Flags().StringVar(&createName, "name", "", "Policy name")
	createCmd.Flags().StringVar(&createDescription, "description", "", "Policy description")
	createCmd.Flags().StringSliceVar(&createLabels, "labels", nil, "Labels as key=value (repeatable)")
	createCmd.Flags().StringSliceVar(&createAnnotations, "annotations", nil, "Annotations as key=value (repeatable)")
	createCmd.Flags().BoolVar(&createReplicateToChildren, "replicate-to-children", false, "Replicate this policy to child projects")
	_ = createCmd.MarkFlagRequired("name")
}
