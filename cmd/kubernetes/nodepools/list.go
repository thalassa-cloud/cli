package nodepools

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
	"github.com/thalassa-cloud/client-go/kubernetes"
	"github.com/thalassa-cloud/client-go/thalassa"
)

const (
	NoHeaderKey = "no-header"

	ClusterFlag = "cluster"
	VpcFlag     = "vpc"
)

var (
	noHeader      bool
	showExactTime bool
	cluster       string
	vpc           string
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"get", "g", "ls", "l"},
	Short:   "Kubernetes Cluster NodePool list",

	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		// Initialize client
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		// Get VPC identity if specified
		var vpcIdentity string
		if vpc != "" {
			vpcIdentity, err = getVpcIdentity(ctx, client, vpc)
			if err != nil {
				return fmt.Errorf("failed to get VPC identity: %w", err)
			}
		}

		// Get clusters
		clusters, err := client.Kubernetes().ListKubernetesClusters(ctx, &kubernetes.ListKubernetesClustersRequest{})
		if err != nil {
			return fmt.Errorf("failed to list clusters: %w", err)
		}

		// Collect node pools data
		body := make([][]string, 0)
		for _, c := range clusters {
			// Skip clusters that don't match VPC filter
			if vpcIdentity != "" && (c.VPC == nil || c.VPC.Identity != vpcIdentity) {
				continue
			}

			// Skip clusters that don't match cluster filter
			if cluster != "" && c.Identity != cluster {
				continue
			}

			// Get node pools for the cluster
			nodePools, err := client.Kubernetes().ListKubernetesNodePools(ctx, c.Identity, &kubernetes.ListKubernetesNodePoolsRequest{})
			if err != nil {
				return fmt.Errorf("failed to list node pools for cluster %s: %w", c.Name, err)
			}

			// Add node pools to the result
			for _, np := range nodePools {
				replicas := formatReplicas(&np)
				body = append(body, []string{
					np.Identity,
					np.Name,
					c.Name,
					replicas,
					np.MachineType.Name,
					string(np.Status),
					formattime.FormatTime(np.CreatedAt.Local(), showExactTime),
				})
			}
		}

		// Print results
		if len(body) == 0 {
			fmt.Println("No Kubernetes Node Pools found")
			return nil
		}

		headers := []string{"ID", "Name", "Cluster", "Replicas", "Type", "Status", "Age"}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print(headers, body)
		}
		return nil
	},
}

// getVpcIdentity retrieves the VPC identity by name, identity, or slug
func getVpcIdentity(ctx context.Context, client thalassa.Client, vpcIdentifier string) (string, error) {
	vpcs, err := client.IaaS().ListVpcs(ctx, &iaas.ListVpcsRequest{})
	if err != nil {
		return "", err
	}

	for _, v := range vpcs {
		if v.Name == vpcIdentifier || v.Identity == vpcIdentifier || v.Slug == vpcIdentifier {
			return v.Identity, nil
		}
	}
	return "", fmt.Errorf("VPC not found: %s", vpcIdentifier)
}

// formatReplicas formats the replicas string based on autoscaling settings
func formatReplicas(np *kubernetes.KubernetesNodePool) string {
	if np.EnableAutoscaling {
		return fmt.Sprintf("%d-%d", np.MinReplicas, np.MaxReplicas)
	}
	return fmt.Sprintf("%d", np.Replicas)
}

func init() {
	listCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	listCmd.Flags().StringVar(&cluster, ClusterFlag, "", "Cluster ID")
	listCmd.Flags().StringVar(&vpc, VpcFlag, "", "VPC ID")
	KubernetesNodePoolsCmd.AddCommand(listCmd)
}
