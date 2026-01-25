package dbaas

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	iaasutil "github.com/thalassa-cloud/cli/internal/iaas"
	"github.com/thalassa-cloud/cli/internal/labels"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/dbaas"
	"github.com/thalassa-cloud/client-go/filters"
)

const NoHeaderKey = "no-header"

var noHeader bool

var (
	showExactTime     bool
	showLabels        bool
	listLabelSelector string
	listEngineFilter  string
	listVpcFilter     string
	listSubnetFilter  string
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

		f := []filters.Filter{}

		// Resolve VPC filter if provided
		if listVpcFilter != "" {
			vpc, err := iaasutil.GetVPCByIdentitySlugOrName(cmd.Context(), client.IaaS(), listVpcFilter)
			if err != nil {
				return fmt.Errorf("failed to get vpc: %w", err)
			}
			f = append(f, &filters.FilterKeyValue{
				Key:   "vpc",
				Value: vpc.Identity,
			})
		}

		// Resolve subnet filter if provided
		if listSubnetFilter != "" {
			subnet, err := iaasutil.GetSubnetByIdentitySlugOrName(cmd.Context(), client.IaaS(), listSubnetFilter)
			if err != nil {
				return fmt.Errorf("failed to get subnet: %w", err)
			}
			f = append(f, &filters.FilterKeyValue{
				Key:   "subnet",
				Value: subnet.Identity,
			})
		}

		// Add engine filter if provided
		if listEngineFilter != "" {
			f = append(f, &filters.FilterKeyValue{
				Key:   "engine",
				Value: listEngineFilter,
			})
		}

		// Add label selector filter if provided
		if listLabelSelector != "" {
			f = append(f, &filters.LabelFilter{
				MatchLabels: labels.ParseLabelSelector(listLabelSelector),
			})
		}

		clusters, err := client.DBaaS().ListDbClusters(cmd.Context(), &dbaas.ListDbClustersRequest{
			Filters: f,
		})
		if err != nil {
			return err
		}
		body := make([][]string, 0, len(clusters))
		for _, cluster := range clusters {

			vpcName := ""
			if cluster.Vpc != nil {
				vpcName = cluster.Vpc.Name
			}

			subnetName := ""
			if cluster.Subnet != nil {
				subnetName = cluster.Subnet.Name
			}

			engineVersion := cluster.EngineVersion
			if cluster.DatabaseEngineVersion != nil {
				engineVersion = cluster.DatabaseEngineVersion.EngineVersion
			}

			instanceType := ""
			if cluster.DatabaseInstanceType != nil {
				instanceType = cluster.DatabaseInstanceType.Name
			}

			row := []string{
				cluster.Identity,
				cluster.Name,
				vpcName,
				subnetName,
				string(cluster.Engine),
				engineVersion,
				instanceType,
				fmt.Sprintf("%d", cluster.Replicas),
				fmt.Sprintf("%d GB", cluster.AllocatedStorage),
				string(cluster.Status),
				formattime.FormatTime(cluster.CreatedAt.Local(), showExactTime),
			}

			if showLabels {
				labelStrs := []string{}
				for k, v := range cluster.Labels {
					labelStrs = append(labelStrs, k+"="+v)
				}
				sort.Strings(labelStrs)
				if len(labelStrs) == 0 {
					labelStrs = []string{"-"}
				}
				row = append(row, strings.Join(labelStrs, ","))
			}

			body = append(body, row)
		}
		if len(body) == 0 {
			fmt.Println("No database clusters found")
			return nil
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			headers := []string{"ID", "Name", "VPC", "Subnet", "Engine", "Version", "Instance Type", "Replicas", "Storage", "Status", "Age"}
			if showLabels {
				headers = append(headers, "Labels")
			}
			table.Print(headers, body)
		}
		return nil
	},
}

func init() {
	DbaasCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	listCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show exact time instead of relative time")
	listCmd.Flags().BoolVar(&showLabels, "show-labels", false, "Show labels")
	listCmd.Flags().StringVarP(&listLabelSelector, "selector", "l", "", "Label selector to filter clusters (format: key1=value1,key2=value2)")
	listCmd.Flags().StringVar(&listEngineFilter, "engine", "", "Filter by database engine (e.g., postgres)")
	listCmd.Flags().StringVar(&listVpcFilter, "vpc", "", "Filter by VPC identity, slug, or name")
	listCmd.Flags().StringVar(&listSubnetFilter, "subnet", "", "Filter by subnet identity, slug, or name")

	// Register completions
	listCmd.RegisterFlagCompletionFunc("vpc", completion.CompleteVPCID)
	listCmd.RegisterFlagCompletionFunc("subnet", completion.CompleteSubnetEnhanced)
}
