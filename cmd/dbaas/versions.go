package dbaas

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/dbaas"
)

var engineType string

// versionsCmd represents the versions command
var versionsCmd = &cobra.Command{
	Use:     "versions",
	Short:   "Get a list of database engine versions",
	Long:    "Get a list of available database engine versions for a specific engine",
	Example: "tcloud dbaas versions --engine postgres\ntcloud dbaas versions --engine postgres --no-header",
	Aliases: []string{"version", "v"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if engineType == "" {
			return fmt.Errorf("engine type is required (use --engine flag)")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Parse the engine type
		engine := dbaas.DbClusterDatabaseEngine(engineType)

		versions, err := client.DBaaS().ListEngineVersions(cmd.Context(), engine, &dbaas.ListEngineVersionsRequest{})
		if err != nil {
			return err
		}
		body := make([][]string, 0, len(versions))
		for _, version := range versions {
			body = append(body, []string{
				version.Identity,
				version.EngineVersion,
				string(version.Engine),
				fmt.Sprintf("%d.%d", version.MajorVersion, version.MinorVersion),
				fmt.Sprintf("%d", version.MajorVersion),
				fmt.Sprintf("%d", version.MinorVersion),
			})
		}
		if len(body) == 0 {
			fmt.Printf("No engine versions found for engine: %s\n", engineType)
			return nil
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Version", "Engine", "Full Version", "Major", "Minor"}, body)
		}
		return nil
	},
}

func init() {
	DbaasCmd.AddCommand(versionsCmd)
	versionsCmd.Flags().StringVar(&engineType, "engine", "", "Database engine type (e.g., postgres)")
	versionsCmd.MarkFlagRequired("engine")
}
