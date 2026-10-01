package accesselevations

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var cancelForce bool

var cancelCmd = &cobra.Command{
	Use:               "cancel <identity>",
	Short:             "Cancel your pending access elevation request",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteIAMMyAccessElevationIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		ok, err := shared.PromptDestructiveUnlessForce(cancelForce, fmt.Sprintf(
			"Are you sure you want to cancel this access elevation request?\n  Request: %s\n", args[0]))
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		req, err := client.IAM().CancelAccessElevation(cmd.Context(), args[0])
		if err != nil {
			return fmt.Errorf("failed to cancel access elevation: %w", err)
		}

		printAccessElevationDetails(req, false)
		return nil
	},
}

func init() {
	AccessElevationsCmd.AddCommand(cancelCmd)
	cancelCmd.Flags().BoolVar(&cancelForce, shared.ForceKey, false, "Skip the confirmation prompt")
}
