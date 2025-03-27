package compute

import (
	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/internal/config/contextstate"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/pkg/thalassa"
)

var getMachineImagesCmd = &cobra.Command{
	Use:     "machine-images",
	Aliases: []string{"machine-images", "machine-image", "images", "image"},
	Short:   "Get a list of machine images",
	Long:    `Get a list of machine images available in the current organisation`,
	Example: `thalassa compute machine-images`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassa.NewClient(
			client.WithBaseURL(contextstate.Server()),
			client.WithOrganisation(contextstate.Organisation()),
			client.WithAuthPersonalToken(contextstate.Token()),
		)
		if err != nil {
			return err
		}
		images, err := client.IaaS().ListMachineImages(cmd.Context())
		if err != nil {
			return err
		}
		body := make([][]string, 0, len(images))
		for _, image := range images {

			row := []string{
				image.Identity,
				image.Name,
				image.Slug,
				image.Architecture,
			}
			body = append(body, row)
		}

		headers := []string{"ID", "Name", "Slug", "Architecture"}

		if showLabels {
			headers = append(headers, "Labels")
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print(headers, body)
		}
		return nil
	},
}

func init() {
	ComputeCmd.AddCommand(getMachineImagesCmd)

	getMachineImagesCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	getMachineImagesCmd.Flags().BoolVar(&showLabels, "show-labels", false, "Show labels associated with machines")
	getMachineImagesCmd.Flags().StringVarP(&outputFormat, "output", "o", "", "Output format. One of: wide")
}
