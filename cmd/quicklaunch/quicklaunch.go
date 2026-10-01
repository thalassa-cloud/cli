package quicklaunch

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/formattime"
	clientql "github.com/thalassa-cloud/client-go/quicklaunch"
)

// QuickLaunchCmd manages quick-launch jobs.
var QuickLaunchCmd = &cobra.Command{
	Use:     "quick-launch",
	Aliases: []string{"quicklaunch", "ql"},
	Short:   "Provision stacks from quick-launch templates",
	Long: `Start and manage quick-launch jobs that provision common stacks asynchronously
(VPC networking, or VPC + Kubernetes) in the current organisation/project scope.`,
}

func printQuickLaunchDetails(ql *clientql.QuickLaunch, exactTime bool) {
	fmt.Printf("ID:          %s\n", ql.Identity)
	fmt.Printf("Name:        %s\n", ql.Name)
	if ql.Slug != "" {
		fmt.Printf("Slug:        %s\n", ql.Slug)
	}
	if ql.Description != "" {
		fmt.Printf("Description: %s\n", ql.Description)
	}
	fmt.Printf("Template:    %s\n", ql.Template)
	fmt.Printf("Status:      %s\n", ql.Status)
	if ql.StatusMessage != "" {
		fmt.Printf("Status msg:  %s\n", ql.StatusMessage)
	}
	if ql.CloudRegionIdentity != "" {
		fmt.Printf("Region:      %s\n", ql.CloudRegionIdentity)
	}
	if ql.VpcCidr != "" {
		fmt.Printf("VPC CIDR:    %s\n", ql.VpcCidr)
	}
	if len(ql.SubnetCidrs) > 0 {
		fmt.Printf("Subnet CIDRs: %s\n", strings.Join(ql.SubnetCidrs, ", "))
	}
	if ql.MachineType != "" {
		fmt.Printf("Machine type: %s\n", ql.MachineType)
	}
	fmt.Printf("Created:     %s\n", formattime.FormatTime(ql.CreatedAt.Local(), exactTime))
	if ql.UpdatedAt != nil {
		fmt.Printf("Updated:     %s\n", formattime.FormatTime(ql.UpdatedAt.Local(), exactTime))
	}

	if len(ql.Resources) > 0 {
		fmt.Println("\nResources:")
		for _, r := range ql.Resources {
			status := r.LastStatus
			if status == "" {
				status = "-"
			}
			fmt.Printf("  - %-22s %-32s %s (%s)\n", r.Type, r.Name, r.Identity, status)
		}
	}
}
