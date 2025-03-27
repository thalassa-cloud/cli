package routetables

import (
	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/config/contextstate"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/pkg/thalassa"
)

const NoHeaderKey = "no-header"

var noHeader bool

var (
	showExactTime bool
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:     "list",
	Short:   "Get a list of routetables",
	Aliases: []string{"g", "get", "ls"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		client, err := thalassa.NewClient(
			client.WithBaseURL(contextstate.Server()),
			client.WithOrganisation(contextstate.Organisation()),
			client.WithAuthPersonalToken(contextstate.Token()),
		)
		if err != nil {
			return err
		}
		routetables, err := client.IaaS().ListRouteTables(cmd.Context())
		if err != nil {
			return err
		}
		body := make([][]string, 0, len(routetables))
		for _, rt := range routetables {
			body = append(body, []string{
				rt.Identity,
				rt.Name,
				rt.Vpc.Name,
				formattime.FormatTime(rt.CreatedAt.Local(), showExactTime),
			})
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "VPC", "Age"}, body)
		}
		return nil
	},
}

func init() {
	RouteTablesCmd.AddCommand(getCmd)

	getCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
}
