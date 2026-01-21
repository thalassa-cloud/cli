package dbaas

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	viewShowExactTime bool
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:               "view",
	Short:             "View database cluster details",
	Long:              "View detailed information about a database cluster.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteDbClusterID,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		clusterIdentity := args[0]

		cluster, err := client.DBaaS().GetDbCluster(cmd.Context(), clusterIdentity)
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("database cluster not found: %s", clusterIdentity)
			}
			return fmt.Errorf("failed to get cluster: %w", err)
		}

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

		body := [][]string{
			{"ID", cluster.Identity},
			{"Name", cluster.Name},
			{"Status", string(cluster.Status)},
			{"Engine", string(cluster.Engine)},
			{"Engine Version", engineVersion},
			{"Instance Type", instanceType},
			{"Replicas", fmt.Sprintf("%d", cluster.Replicas)},
			{"Storage", fmt.Sprintf("%d GB", cluster.AllocatedStorage)},
			{"VPC", vpcName},
			{"Created", formattime.FormatTime(cluster.CreatedAt.Local(), viewShowExactTime)},
		}

		if cluster.Description != "" {
			body = append(body, []string{"Description", cluster.Description})
		}

		if len(cluster.Labels) > 0 {
			labelStrs := []string{}
			for k, v := range cluster.Labels {
				labelStrs = append(labelStrs, k+"="+v)
			}
			sort.Strings(labelStrs)
			body = append(body, []string{"Labels", strings.Join(labelStrs, ", ")})
		}

		if len(cluster.Annotations) > 0 {
			annotationStrs := []string{}
			for k, v := range cluster.Annotations {
				annotationStrs = append(annotationStrs, k+"="+v)
			}
			sort.Strings(annotationStrs)
			body = append(body, []string{"Annotations", strings.Join(annotationStrs, ", ")})
		}

		if cluster.DeleteProtection {
			body = append(body, []string{"Delete Protection", "enabled"})
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"Field", "Value"}, body)
		}

		return nil
	},
}

func init() {
	DbaasCmd.AddCommand(viewCmd)

	viewCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	viewCmd.Flags().BoolVar(&viewShowExactTime, "exact-time", false, "Show exact time instead of relative time")
}
