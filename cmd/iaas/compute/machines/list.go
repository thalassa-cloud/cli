package machines

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/config/contextstate"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"

	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/thalassa-cloud/client-go/pkg/client"
	"github.com/thalassa-cloud/client-go/thalassa"
)

const NoHeaderKey = "no-header"

var noHeader bool

var (
	showExactTime bool
	showLabels    bool
	outputFormat  string
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:     "list",
	Short:   "Get a list of machines",
	Aliases: []string{"g", "get", "ls"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassa.NewClient(
			client.WithBaseURL(contextstate.Server()),
			client.WithOrganisation(contextstate.Organisation()),
			client.WithAuthPersonalToken(contextstate.Token()),
		)
		if err != nil {
			return err
		}
		machines, err := client.IaaS().ListMachines(cmd.Context())
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
						volumeAttachments = append(volumeAttachments, attachment.DeviceName)
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
	getCmd.Flags().StringVarP(&outputFormat, "output", "o", "", "Output format. One of: wide")
}
