package subnets

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	iaasutil "github.com/thalassa-cloud/cli/internal/iaas"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
)

const (
	CreateFlagName        = "name"
	CreateFlagDescription = "description"
	CreateFlagVpc         = "vpc"
	CreateFlagCIDR        = "cidr"

	CreateFlagLabels      = "labels"
	CreateFlagAnnotations = "annotations"
)

var (
	createSubnetValues = iaas.CreateSubnet{}
	createSubnetWait   bool
)

// getCmd represents the get command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a subnet",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		tcclient, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		if createSubnetValues.Name == "" {
			return fmt.Errorf("name is required")
		}
		if createSubnetValues.VpcIdentity == "" {
			return fmt.Errorf("vpc is required")
		}
		if createSubnetValues.Cidr == "" {
			return fmt.Errorf("cidr is required")
		}

		vpc, err := iaasutil.GetVPCByIdentitySlugOrName(cmd.Context(), tcclient.IaaS(), createSubnetValues.VpcIdentity)
		if err != nil {
			return err
		}
		createSubnetValues.VpcIdentity = vpc.Identity

		subnet, err := tcclient.IaaS().CreateSubnet(cmd.Context(), createSubnetValues)
		if err != nil {
			return err
		}

		if createSubnetWait {
			ctxWithTimeout, cancel := context.WithTimeout(cmd.Context(), 10*time.Minute)
			defer cancel()

			fmt.Println("Waiting for subnet to be ready...")
			for {
				subnet, err = tcclient.IaaS().GetSubnet(ctxWithTimeout, subnet.Identity)
				if err != nil {
					return fmt.Errorf("failed to get subnet: %w", err)
				}
				// Subnet is ready when status is "ready" or "available"
				status := string(subnet.Status)
				if strings.EqualFold(status, "ready") || strings.EqualFold(status, "available") {
					break
				}
				// Check for failed state
				if strings.EqualFold(status, "failed") || strings.EqualFold(status, "error") {
					return fmt.Errorf("subnet creation failed with status: %s", status)
				}
				select {
				case <-ctxWithTimeout.Done():
					return fmt.Errorf("timeout waiting for subnet %s to be ready (current status: %s)", subnet.Identity, status)
				case <-time.After(2 * time.Second):
					// Continue polling
				}
			}
			fmt.Println("Subnet is ready")
		}

		body := make([][]string, 0, 1)
		body = append(body, []string{
			subnet.Identity,
			subnet.Name,
			vpc.Name,
			subnet.Cidr,
			formattime.FormatTime(subnet.CreatedAt.Local(), showExactTime),
		})
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "VPC", "CIDR", "Age"}, body)
		}
		return nil
	},
}

func init() {
	SubnetsCmd.AddCommand(createCmd)
	createCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	createCmd.Flags().StringVar(&createSubnetValues.Name, CreateFlagName, "", "Name of the subnet")
	createCmd.Flags().StringVar(&createSubnetValues.Description, CreateFlagDescription, "", "Description of the subnet")
	createCmd.Flags().StringVar(&createSubnetValues.VpcIdentity, CreateFlagVpc, "", "VPC of the subnet")
	createCmd.Flags().StringVar(&createSubnetValues.Cidr, CreateFlagCIDR, "", "CIDR of the subnet")
	createCmd.Flags().BoolVar(&createSubnetWait, "wait", false, "Wait for the subnet to be ready before returning")

	// Register completions
	createCmd.RegisterFlagCompletionFunc("vpc", completion.CompleteVPCID)
}
