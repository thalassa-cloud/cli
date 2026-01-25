package kubernetes

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/kubernetes"
)

const NoHeaderKey = "no-header"

const (
	VpcFlag = "vpc"
)

var noHeader bool

var (
	showExactTime bool
	vpc           string
)

// listCmd represents the get command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "Get a list of Kubernetes clusters",
	Long:    "Get a list of Kubernetes clusters within your organisation",
	Aliases: []string{"g", "get", "ls", "clusters", "cluster"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		f := []filters.Filter{}
		if vpc != "" {
			f = append(f, &filters.FilterKeyValue{
				Key:   "vpc",
				Value: vpc,
			})
		}
		clusters, err := client.Kubernetes().ListKubernetesClusters(cmd.Context(), &kubernetes.ListKubernetesClustersRequest{
			Filters: f,
		})
		if err != nil {
			return err
		}
		body := make([][]string, 0, len(clusters))
		for _, cluster := range clusters {

			vpcName := ""
			if cluster.VPC != nil {
				vpcName = cluster.VPC.Name
			}
			body = append(body, []string{
				cluster.Identity,
				cluster.Name,
				vpcName,
				cluster.ClusterVersion.Name,
				string(cluster.ClusterType),
				cluster.Status,
				formattime.FormatTime(cluster.CreatedAt.Local(), showExactTime),
			})
		}
		if len(body) == 0 {
			fmt.Println("No Clusters found")
			return nil
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Vpc", "Version", "Type", "Status", "Age"}, body)
		}
		return nil
	},
}

func init() {
	KubernetesCmd.AddCommand(listCmd)

	// flags
	listCmd.Flags().StringVar(&vpc, VpcFlag, "", "VPC ID")
	listCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	listCmd.RegisterFlagCompletionFunc(VpcFlag, completion.CompleteVPCID)
}
