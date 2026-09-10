package routetables

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
)

var (
	createName        string
	createDescription string
	createVPC         string
	createLabels      []string
	createAnnotations []string
)

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a route table",
	Example: "tcloud networking routetables create --name private --vpc vpc-123",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		if createName == "" {
			return fmt.Errorf("name is required")
		}
		if createVPC == "" {
			return fmt.Errorf("vpc is required")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		req := iaas.CreateRouteTable{
			Name:        createName,
			Labels:      shared.KeyValuePairsToMap(createLabels),
			Annotations: shared.KeyValuePairsToMap(createAnnotations),
			VpcIdentity: createVPC,
		}
		if createDescription != "" {
			req.Description = &createDescription
		}

		rt, err := client.IaaS().CreateRouteTable(cmd.Context(), req)
		if err != nil {
			return fmt.Errorf("failed to create route table: %w", err)
		}

		fmt.Printf("Route table created successfully\n")
		fmt.Printf("ID: %s\n", rt.Identity)
		fmt.Printf("Name: %s\n", rt.Name)
		return nil
	},
}

func init() {
	RouteTablesCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&createName, "name", "", "Name of the route table")
	createCmd.Flags().StringVar(&createDescription, "description", "", "Description")
	createCmd.Flags().StringVar(&createVPC, "vpc", "", "VPC identity")
	createCmd.Flags().StringSliceVar(&createLabels, "labels", nil, "Labels as key=value (repeatable)")
	createCmd.Flags().StringSliceVar(&createAnnotations, "annotations", nil, "Annotations as key=value (repeatable)")
	_ = createCmd.MarkFlagRequired("name")
	_ = createCmd.MarkFlagRequired("vpc")
	_ = createCmd.RegisterFlagCompletionFunc("vpc", completeVPCID)
}
