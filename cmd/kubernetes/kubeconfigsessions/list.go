package kubeconfigsessions

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/kuberesolve"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/kubernetes"
)

const (
	NoHeaderKey = "no-header"
	ClusterFlag = "cluster"
)

var (
	noHeader      bool
	showExactTime bool
	listCluster   string
)

var listCmd = &cobra.Command{
	Use:               "list [cluster]",
	Aliases:           []string{"ls", "l"},
	Short:             "List kubeconfig sessions for a Kubernetes cluster",
	Args:              cobra.MaximumNArgs(1),
	ValidArgsFunction: completion.CompleteKubernetesCluster,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		clusterRef := listCluster
		if len(args) > 0 {
			if clusterRef != "" && clusterRef != args[0] {
				return fmt.Errorf("provide the cluster as an argument or --cluster, not both with different values")
			}
			clusterRef = args[0]
		}
		if clusterRef == "" {
			return fmt.Errorf("cluster is required (argument or --cluster)")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		cluster, err := kuberesolve.ResolveKubernetesClusterRef(ctx, client.Kubernetes(), clusterRef)
		if err != nil {
			return err
		}

		sessions, err := client.Kubernetes().ListKubeconfigSessions(ctx, cluster.Identity)
		if err != nil {
			return fmt.Errorf("failed to list kubeconfig sessions: %w", err)
		}

		if len(sessions) == 0 {
			fmt.Println("No kubeconfig sessions found")
			return nil
		}

		body := make([][]string, 0, len(sessions))
		for _, session := range sessions {
			body = append(body, []string{
				session.Identity,
				sessionSubject(session),
				formattime.FormatTime(session.CreatedAt.Local(), showExactTime),
				formatOptionalTime(session.LastUsedAt, showExactTime),
				formattime.FormatTime(session.ExpiresAt.Local(), showExactTime),
			})
		}

		headers := []string{"Identity", "User/ServiceAccount", "CreatedAt", "LastUsedAt", "ExpiresAt"}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print(headers, body)
		}
		return nil
	},
}

func sessionSubject(session kubernetes.KubernetesClusterSession) string {
	if session.User != nil {
		return "user:" + shared.UserPtrDisplay(session.User)
	}
	if session.ServiceAccount != nil {
		name := session.ServiceAccount.Name
		if name == "" {
			name = session.ServiceAccount.Identity
		}
		return "serviceaccount:" + name
	}
	return "-"
}

func formatOptionalTime(t *time.Time, exact bool) string {
	if t == nil {
		return "-"
	}
	return formattime.FormatTime(t.Local(), exact)
}

func init() {
	KubeconfigSessionsCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	listCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
	listCmd.Flags().StringVar(&listCluster, ClusterFlag, "", "Cluster identity, name, or slug")
	_ = listCmd.RegisterFlagCompletionFunc(ClusterFlag, completion.CompleteKubernetesCluster)
}
