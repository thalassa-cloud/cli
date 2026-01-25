package machines

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/labels"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/filters"
	"github.com/thalassa-cloud/client-go/iaas"

	"k8s.io/apimachinery/pkg/api/resource"
)

const NoHeaderKey = "no-header"

var noHeader bool

var (
	showExactTime     bool
	showLabels        bool
	listLabelSelector string
	outputFormat      string
	listRegionFilter  string
	listVpcFilter     string
	listStatusFilter  string
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:     "list",
	Short:   "Get a list of machines",
	Aliases: []string{"g", "get", "ls"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		f := []filters.Filter{}
		if listLabelSelector != "" {
			f = append(f, &filters.LabelFilter{
				MatchLabels: labels.ParseLabelSelector(listLabelSelector),
			})
		}
		if listRegionFilter != "" {
			f = append(f, &filters.FilterKeyValue{
				Key:   "region",
				Value: listRegionFilter,
			})
		}
		if listVpcFilter != "" {
			f = append(f, &filters.FilterKeyValue{
				Key:   "vpc",
				Value: listVpcFilter,
			})
		}
		if listStatusFilter != "" {
			f = append(f, &filters.FilterKeyValue{
				Key:   "status",
				Value: listStatusFilter,
			})
		}

		machines, err := client.IaaS().ListMachines(cmd.Context(), &iaas.ListMachinesRequest{
			Filters: f,
		})
		if err != nil {
			return err
		}
		body := make([][]string, 0, len(machines))
		for _, machine := range machines {
			ips := []string{}
			for _, ip := range machine.Interfaces {
				ips = append(ips, ip.IPAddresses...)
			}

			regionName := ""
			if machine.Vpc.CloudRegion != nil {
				regionName = machine.Vpc.CloudRegion.Name
				if regionName == "" {
					regionName = machine.Vpc.CloudRegion.Identity
				}
				if regionName == "" {
					regionName = machine.Vpc.CloudRegion.Slug
				}
				if regionName == "" {
					regionName = machine.Vpc.CloudRegion.Identity
				}
			}

			row := []string{
				machine.Identity,
				machine.Name,
				machine.Status.Status,
				machine.Vpc.Name,
				regionName,
				machine.MachineType.Name,
				strings.Join(ips, ", "),
				formattime.FormatTime(machine.CreatedAt.Local(), showExactTime),
			}

			if outputFormat == "wide" {
				subnetName := "-"
				if machine.Subnet != nil {
					subnetName = machine.Subnet.Name
				}

				volumeAttachments := []string{}
				for _, attachment := range machine.VolumeAttachments {
					if attachment.PersistentVolume != nil {
						volumeAttachments = append(volumeAttachments, attachment.PersistentVolume.Name)
					} else {
						volumeAttachments = append(volumeAttachments, attachment.Serial)
					}
				}
				sort.Strings(volumeAttachments)
				row = append(row, subnetName, strings.Join(volumeAttachments, ","))

				cpu := machine.MachineType.Vcpus
				memory := resource.NewQuantity(int64(machine.MachineType.RamMb*1024*1024), resource.BinarySI).String()
				row = append(row, fmt.Sprintf("%d", cpu), memory)
			}

			if showLabels {
				labels := []string{}
				for k, v := range machine.Labels {
					labels = append(labels, k+"="+v)
				}
				sort.Strings(labels)
				if len(labels) == 0 {
					labels = []string{"-"}
				}
				row = append(row, strings.Join(labels, ","))
			}

			body = append(body, row)
		}

		headers := []string{"ID", "Name", "Status", "VPC", "Region", "Type", "IP", "Age"}
		if outputFormat == "wide" {
			headers = append(headers, "Subnet", "Volumes", "CPU", "Memory")
		}

		if showLabels {
			headers = append(headers, "Labels")
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
	MachinesCmd.AddCommand(getCmd)

	getCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	getCmd.Flags().BoolVar(&showExactTime, "show-exact-time", false, "Show exact time instead of relative time")
	getCmd.Flags().BoolVar(&showLabels, "show-labels", false, "Show labels associated with machines")
	getCmd.Flags().StringVarP(&listLabelSelector, "selector", "l", "", "Label selector to filter machines (format: key1=value1,key2=value2)")
	getCmd.Flags().StringVarP(&outputFormat, "output", "o", "", "Output format. One of: wide")
	getCmd.Flags().StringVar(&listRegionFilter, "region", "", "Region of the machine")
	getCmd.Flags().StringVar(&listVpcFilter, "vpc", "", "VPC of the machine")
	getCmd.Flags().StringVar(&listStatusFilter, "status", "", "Status of the machine")

	// Register completions
	getCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
	getCmd.RegisterFlagCompletionFunc("vpc", completion.CompleteVPCID)
}
