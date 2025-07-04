package me

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

const (
	NoHeaderKey = "no-header"
)

var (
	noHeader bool
	slugOnly bool
)

var organisationsCmd = &cobra.Command{
	Use:   "organisations",
	Short: "Get information about the current user's organisations",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		organisations, err := client.Me().ListMyOrganisations(ctx)
		if err != nil {
			return fmt.Errorf("failed to get organisations: %w", err)
		}
		body := make([][]string, 0)
		for _, organisation := range organisations {
			if slugOnly {
				body = append(body, []string{
					organisation.Slug,
				})
			} else {
				body = append(body, []string{
					organisation.Identity,
					organisation.Name,
					organisation.Slug,
				})
			}
		}

		headers := []string{"ID", "Name", "Slug"}
		if slugOnly {
			headers = []string{"Slug"}
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
	MeCmd.AddCommand(organisationsCmd)
	organisationsCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "do not print headers")
	organisationsCmd.Flags().BoolVar(&slugOnly, "slug-only", false, "only print the slug")
}
