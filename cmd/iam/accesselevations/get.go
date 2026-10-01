package accesselevations

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientiam "github.com/thalassa-cloud/client-go/iam"
)

var getCmd = &cobra.Command{
	Use:               "get <identity>",
	Aliases:           []string{"view", "show", "describe"},
	Short:             "Show an access elevation request",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteIAMAccessElevationIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		var req *clientiam.IamAccessElevationRequest
		if mineFlag {
			req, err = client.IAM().GetMyAccessElevation(cmd.Context(), args[0])
		} else {
			req, err = client.IAM().GetAccessElevation(cmd.Context(), args[0])
		}
		if err != nil {
			return fmt.Errorf("failed to get access elevation: %w", err)
		}

		printAccessElevationDetails(req, showExactTime)
		return nil
	},
}

func init() {
	AccessElevationsCmd.AddCommand(getCmd)
	getCmd.Flags().BoolVar(&mineFlag, "mine", false, "Fetch via your own requests API")
	getCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
}
