package kubernetesversions

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/internal/config/contextstate"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/thalassa"
)

const NoHeaderKey = "no-header"

var noHeader bool

var (
	showExactTime bool
)

var KubernetesVersionsCmd = &cobra.Command{
	Use:     "versions",
	Aliases: []string{"version", "v"},
	Short:   "Kubernetes Versions management",

	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassa.NewClient(
			client.WithBaseURL(contextstate.Server()),
			client.WithOrganisation(contextstate.Organisation()),
			client.WithAuthPersonalToken(contextstate.Token()),
		)
		if err != nil {
			return err
		}
		versions, err := client.Kubernetes().ListKubernetesVersions(cmd.Context())
		if err != nil {
			return err
		}
		body := make([][]string, 0, len(versions))
		for _, version := range versions {

			body = append(body, []string{
				version.Identity,
				version.Name,

				version.KubernetesVersion,
				version.ContainerdVersion,
				// version.CNIPluginsVersion,
				// version.CrictlVersion,
				version.RuncVersion,
				// version.CiliumVersion,
				// version.CloudControllerManagerVersion,
				// version.IstioVersion,
				formattime.FormatTime(version.CreatedAt.Local(), showExactTime),
			})
		}
		if len(body) == 0 {
			fmt.Println("No Kubernetes Verions found")
			return nil
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Kubernetes", "Containerd", "Runc", "Age"}, body)
		}
		return nil
	},
}

func init() {
	KubernetesVersionsCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
}
