package listeners

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/labels"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/iaas"
)

var (
	noHeader          bool
	showExactTime     bool
	showLabels        bool
	listLabelSelector string
	loadbalancer      string
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List load balancer listeners",
	Long:    "List listeners for a load balancer.",
	Example: "tcloud networking loadbalancers listeners list --loadbalancer lb-123",
	Aliases: []string{"g", "get", "ls"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if loadbalancer == "" {
			return fmt.Errorf("--loadbalancer is required")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		f := filters.Filters{}
		if listLabelSelector != "" {
			f = append(f, &filters.LabelFilter{
				MatchLabels: labels.ParseLabelSelector(listLabelSelector),
			})
		}

		listeners, err := client.IaaS().ListListeners(cmd.Context(), &iaas.ListLoadbalancerListenersRequest{
			Loadbalancer: loadbalancer,
			Filters:      f,
		})
		if err != nil {
			return err
		}

		body := make([][]string, 0, len(listeners))
		for _, listener := range listeners {
			targetGroup := "-"
			if listener.TargetGroup != nil {
				targetGroup = listener.TargetGroup.Name
			}

			row := []string{
				listener.Identity,
				listener.Name,
				fmt.Sprintf("%d", listener.Port),
				string(listener.Protocol),
				targetGroup,
				formattime.FormatTime(listener.CreatedAt.Local(), showExactTime),
			}

			if showLabels {
				labelPairs := make([]string, 0, len(listener.Labels))
				for k, v := range listener.Labels {
					labelPairs = append(labelPairs, k+"="+v)
				}
				sort.Strings(labelPairs)
				if len(labelPairs) == 0 {
					labelPairs = []string{"-"}
				}
				row = append(row, strings.Join(labelPairs, ","))
			}

			body = append(body, row)
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			headers := []string{"ID", "Name", "Port", "Protocol", "Target Group", "Age"}
			if showLabels {
				headers = append(headers, "Labels")
			}
			table.Print(headers, body)
		}
		return nil
	},
}

func init() {
	ListenersCmd.AddCommand(listCmd)

	listCmd.Flags().StringVar(&loadbalancer, LoadbalancerFlag, "", "Load balancer identity")
	listCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	listCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show exact time instead of relative time")
	listCmd.Flags().BoolVar(&showLabels, "show-labels", false, "Show labels")
	listCmd.Flags().StringVarP(&listLabelSelector, "selector", "l", "", "Label selector (format: key1=value1,key2=value2)")

	listCmd.MarkFlagRequired(LoadbalancerFlag)
	listCmd.RegisterFlagCompletionFunc(LoadbalancerFlag, completeLoadbalancerID)
}
