package machines

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/kuberesolve"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/kubernetes"
)

const (
	NoHeaderKey  = "no-header"
	ClusterFlag  = "cluster"
	NodePoolFlag = "nodepool"
)

var (
	noHeader      bool
	showExactTime bool
	cluster       string
	nodePool      string
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"get", "g", "ls", "l"},
	Short:   "List machines in a Kubernetes cluster",
	Long:    "Lists worker machines (nodes) across node pools in the given cluster.",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		if cluster == "" {
			return fmt.Errorf("--cluster is required")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		cl, err := kuberesolve.ResolveKubernetesClusterRef(ctx, client.Kubernetes(), cluster)
		if err != nil {
			return err
		}

		nodePools, err := client.Kubernetes().ListKubernetesNodePools(ctx, cl.Identity, &kubernetes.ListKubernetesNodePoolsRequest{})
		if err != nil {
			return fmt.Errorf("failed to list node pools: %w", err)
		}

		body := make([][]string, 0)
		for _, np := range nodePools {
			if nodePool != "" && !matchesNodePoolRef(&np, nodePool) {
				continue
			}

			machines, err := client.Kubernetes().ListNodePoolMachines(ctx, cl.Identity, np.Identity)
			if err != nil {
				return fmt.Errorf("failed to list machines for node pool %s: %w", np.Name, err)
			}

			for _, m := range machines {
				body = append(body, []string{
					m.Identity,
					m.MachineName,
					np.Name,
					cl.Name,
					nodeReadyStatus(m.SystemInfo.Conditions),
					m.SystemInfo.KubeletVersion,
					formatNodeInternalIP(m.SystemInfo.Addresses),
					formattime.FormatTime(m.CreatedAt.Local(), showExactTime),
				})
			}
		}

		if len(body) == 0 {
			fmt.Println("No Kubernetes machines found")
			return nil
		}

		headers := []string{"ID", "Name", "Node pool", "Cluster", "Ready", "Kubelet", "IP", "Age"}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print(headers, body)
		}
		return nil
	},
}

func matchesNodePoolRef(np *kubernetes.KubernetesNodePool, ref string) bool {
	return strings.EqualFold(np.Identity, ref) ||
		strings.EqualFold(np.Name, ref) ||
		strings.EqualFold(np.Slug, ref)
}

func nodeReadyStatus(conditions []kubernetes.NodeCondition) string {
	for _, c := range conditions {
		if c.Type == "Ready" {
			return c.Status
		}
	}
	return "Unknown"
}

func formatNodeInternalIP(addresses []kubernetes.NodeAddress) string {
	ips := make([]string, 0, len(addresses))
	for _, a := range addresses {
		if a.Type != "InternalIP" {
			continue
		}
		ips = append(ips, a.Address)
	}
	return strings.Join(ips, ", ")
}

func init() {
	MachinesCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	listCmd.Flags().BoolVar(&showExactTime, "show-exact-time", false, "Show exact time instead of relative time")
	listCmd.Flags().StringVar(&cluster, ClusterFlag, "", "Cluster identity, name, or slug")
	listCmd.Flags().StringVar(&nodePool, NodePoolFlag, "", "Filter by node pool identity, name, or slug")
	_ = listCmd.MarkFlagRequired(ClusterFlag)
	listCmd.RegisterFlagCompletionFunc(ClusterFlag, completion.CompleteKubernetesCluster)
	listCmd.RegisterFlagCompletionFunc(NodePoolFlag, completion.CompleteKubernetesNodePool)
}
