package loadbalancers

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	iaasutil "github.com/thalassa-cloud/cli/internal/iaas"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
)

var (
	createName                  string
	createDescription           string
	createSubnet                string
	createInternal              bool
	createDeleteProtection      bool
	createLabels                []string
	createAnnotations           []string
	createSecurityGroups        []string
	createWait                  bool
)

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a load balancer",
	Long:    "Create a new load balancer in the specified subnet.",
	Example: "tcloud networking loadbalancers create --name web --subnet subnet-123\ntcloud networking loadbalancers create --name internal --subnet subnet-123 --internal --wait",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if createName == "" {
			return fmt.Errorf("name is required")
		}
		if createSubnet == "" {
			return fmt.Errorf("subnet is required")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		subnet, err := iaasutil.GetSubnetByIdentitySlugOrName(cmd.Context(), client.IaaS(), createSubnet)
		if err != nil {
			return fmt.Errorf("failed to get subnet: %w", err)
		}

		req := iaas.CreateLoadbalancer{
			Name:                     createName,
			Description:              createDescription,
			Subnet:                   subnet.Identity,
			InternalLoadbalancer:     createInternal,
			DeleteProtection:         createDeleteProtection,
			Labels:                   parseKeyValueSlice(createLabels),
			Annotations:              parseKeyValueSlice(createAnnotations),
			SecurityGroupAttachments: createSecurityGroups,
		}

		lb, err := client.IaaS().CreateLoadbalancer(cmd.Context(), req)
		if err != nil {
			return err
		}

		if createWait {
			ctxWithTimeout, cancel := context.WithTimeout(cmd.Context(), 10*time.Minute)
			defer cancel()
			fmt.Println("Waiting for load balancer to be ready...")
			if err := client.IaaS().WaitUntilLoadbalancerIsReady(ctxWithTimeout, lb.Identity); err != nil {
				return fmt.Errorf("failed waiting for load balancer: %w", err)
			}
			lb, err = client.IaaS().GetLoadbalancer(ctxWithTimeout, lb.Identity)
			if err != nil {
				return fmt.Errorf("failed to get load balancer: %w", err)
			}
			fmt.Println("Load balancer is ready")
		}

		fmt.Printf("Load balancer created successfully\n")
		fmt.Printf("ID: %s\n", lb.Identity)
		fmt.Printf("Name: %s\n", lb.Name)
		fmt.Printf("Status: %s\n", lb.Status)
		return nil
	},
}

func init() {
	LoadbalancersCmd.AddCommand(createCmd)

	createCmd.Flags().StringVar(&createName, "name", "", "Name of the load balancer")
	createCmd.Flags().StringVar(&createDescription, "description", "", "Description of the load balancer")
	createCmd.Flags().StringVar(&createSubnet, "subnet", "", "Subnet identity, slug, or name")
	createCmd.Flags().BoolVar(&createInternal, "internal", false, "Create an internal load balancer (no public IP)")
	createCmd.Flags().BoolVar(&createDeleteProtection, "delete-protection", false, "Enable delete protection")
	createCmd.Flags().StringSliceVar(&createLabels, "labels", []string{}, "Labels in key=value format")
	createCmd.Flags().StringSliceVar(&createAnnotations, "annotations", []string{}, "Annotations in key=value format")
	createCmd.Flags().StringSliceVar(&createSecurityGroups, "security-groups", []string{}, "Security group identities to attach")
	createCmd.Flags().BoolVar(&createWait, "wait", false, "Wait for the load balancer to be ready")

	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("subnet")

	createCmd.RegisterFlagCompletionFunc("subnet", completeSubnetID)
	createCmd.RegisterFlagCompletionFunc("security-groups", completeSecurityGroupID)
}
