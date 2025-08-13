package dbaas

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/dbaas/dbaasalphav1"
)

const NoHeaderKey = "no-header"

var noHeader bool

var (
	showExactTime bool
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "Get a list of database clusters",
	Long:    "Get a list of database clusters within your organisation",
	Example: "tcloud dbaas list\ntcloud dbaas list --no-header\ntcloud dbaas list --exact-time",
	Aliases: []string{"g", "get", "ls", "clusters", "cluster"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		clusters, err := client.DbaaSAlphaV1().ListDbClusters(cmd.Context(), &dbaasalphav1.ListDbClustersRequest{})
		if err != nil {
			return err
		}
		body := make([][]string, 0, len(clusters))
		for _, cluster := range clusters {

			vpcName := ""
			if cluster.Vpc != nil {
				vpcName = cluster.Vpc.Name
			}

			engineVersion := cluster.EngineVersion
			if cluster.DatabaseEngineVersion != nil {
				engineVersion = cluster.DatabaseEngineVersion.EngineVersion
			}

			instanceType := ""
			if cluster.DatabaseInstanceType != nil {
				instanceType = cluster.DatabaseInstanceType.Name
			}

			body = append(body, []string{
				cluster.Identity,
				cluster.Name,
				vpcName,
				string(cluster.Engine),
				engineVersion,
				instanceType,
				fmt.Sprintf("%d", cluster.Replicas),
				fmt.Sprintf("%d GB", cluster.AllocatedStorage),
				string(cluster.Status),
				formattime.FormatTime(cluster.CreatedAt.Local(), showExactTime),
			})
		}
		if len(body) == 0 {
			fmt.Println("No database clusters found")
			return nil
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "VPC", "Engine", "Version", "Instance Type", "Replicas", "Storage", "Status", "Age"}, body)
		}
		return nil
	},
}

func init() {
	DbaasCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	listCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show exact time instead of relative time")
}
