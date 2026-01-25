package dbaas

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/dbaas"
)

// instanceTypesCmd represents the instance-types command
var instanceTypesCmd = &cobra.Command{
	Use:     "instance-types",
	Short:   "Get a list of database instance types",
	Long:    "Get a list of available database instance types within your organisation",
	Example: "tcloud dbaas instance-types\ntcloud dbaas instance-types --no-header",
	Aliases: []string{"instancetypes", "instance-type", "it"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		instanceTypes, err := client.DBaaS().ListDatabaseInstanceTypes(cmd.Context(), &dbaas.ListDatabaseInstanceTypesRequest{})
		if err != nil {
			return err
		}

		sort.Slice(instanceTypes, func(i, j int) bool {
			// sort by memory
			if instanceTypes[i].Memory == instanceTypes[j].Memory {
				return instanceTypes[i].Name < instanceTypes[j].Name
			}
			return instanceTypes[i].Memory < instanceTypes[j].Memory
		})

		body := make([][]string, 0, len(instanceTypes))
		for _, instanceType := range instanceTypes {
			body = append(body, []string{
				instanceType.Identity,
				instanceType.Name,
				instanceType.CategorySlug,
				fmt.Sprintf("%d vCPU", instanceType.Cpus),
				fmt.Sprintf("%d GB", instanceType.Memory),
				instanceType.Architecture,
			})
		}
		if len(body) == 0 {
			fmt.Println("No database instance types found")
			return nil
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Category", "vCPU", "Memory", "Architecture"}, body)
		}
		return nil
	},
}

func init() {
	DbaasCmd.AddCommand(instanceTypesCmd)
}
