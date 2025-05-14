package compute

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"

	"k8s.io/apimachinery/pkg/api/resource"
)

const NoHeaderKey = "no-header"

var noHeader bool

var (
	showLabels     bool
	outputFormat   string
	categoryFilter string
)

// getMachineTypesCmd represents the get command
var getMachineTypesCmd = &cobra.Command{
	Use:     "machine-types",
	Short:   "Get a list of machine types",
	Long:    `Get a list of machine types available in the current organisation`,
	Example: `thalassa compute machine-types`,
	Aliases: []string{"machine-types", "machine-type", "machinetypes", "machinetype", "instancetypes", "instancetype", "types", "type"},
	RunE: func(cmd *cobra.Command, args []string) error {

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		machinetypeCategories, err := client.IaaS().ListMachineTypeCategories(cmd.Context())
		if err != nil {
			return err
		}

		totalEntries := 0
		for _, category := range machinetypeCategories {
			totalEntries += len(category.MachineTypes)
		}

		body := make([][]string, 0, totalEntries)
		for _, category := range machinetypeCategories {
			if categoryFilter != "" && !strings.EqualFold(category.Name, categoryFilter) {
				continue
			}
			for _, machinetype := range category.MachineTypes {
				memory := resource.NewQuantity(int64(machinetype.RamMb*1024*1024), resource.BinarySI).String()
				row := []string{
					machinetype.Name,
					machinetype.Slug,
					category.Name,
					fmt.Sprintf("%d", machinetype.Vcpus),
					memory,
				}
				if outputFormat == "wide" {
					row = append(row, machinetype.Description)
				}
				body = append(body, row)
			}
		}

		headers := []string{"Name", "Slug", "Category", "CPU", "Memory"}
		if outputFormat == "wide" {
			headers = append(headers, "Description")
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
	ComputeCmd.AddCommand(getMachineTypesCmd)

	getMachineTypesCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	getMachineTypesCmd.Flags().BoolVar(&showLabels, "show-labels", false, "Show labels associated with machines")
	getMachineTypesCmd.Flags().StringVarP(&outputFormat, "output", "o", "", "Output format. One of: wide")
	getMachineTypesCmd.Flags().StringVar(&categoryFilter, "category", "", "Filter by category")
}
