package quicklaunch

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	iaasutil "github.com/thalassa-cloud/cli/internal/iaas"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
	clientql "github.com/thalassa-cloud/client-go/quicklaunch"
)

var (
	createName        string
	createDescription string
	createTemplate    string
	createRegion      string
	createVpcCidr     string
	createSubnetCidrs []string
	createMachineType string
	createLabels      []string
	createAnnotations []string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Start a quick-launch job",
	Long: `Start asynchronous provisioning from a template.

Templates:
  vpc          VPC with subnets and NAT gateway
  kubernetes  VPC stack plus a Kubernetes cluster

Provisioning continues after create returns; use get/logs to follow progress.`,
	Example: `  tcloud quick-launch create --name demo --region nl-ams --template vpc
  tcloud quick-launch create --name demo-k8s --region nl-ams --template kubernetes --machine-type gp.medium`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, _ []string) error {
		if createName == "" {
			return fmt.Errorf("--name is required")
		}
		if createRegion == "" {
			return fmt.Errorf("--region is required")
		}

		template := clientql.QuickLaunchTemplateType(createTemplate)
		switch template {
		case "", clientql.QuickLaunchTemplateVPC, clientql.QuickLaunchTemplateKubernetes:
		default:
			return fmt.Errorf("--template must be vpc or kubernetes (got %q)", createTemplate)
		}
		if template == "" {
			template = clientql.QuickLaunchTemplateVPC
		}
		if createMachineType != "" && template != clientql.QuickLaunchTemplateKubernetes {
			return fmt.Errorf("--machine-type is only valid with --template kubernetes")
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		regions, err := client.IaaS().ListRegions(cmd.Context(), &iaas.ListRegionsRequest{})
		if err != nil {
			return fmt.Errorf("failed to list regions: %w", err)
		}
		region, err := iaasutil.FindRegionByIdentitySlugOrNameWithError(regions, createRegion)
		if err != nil {
			return err
		}

		req := clientql.QuickLaunchRequest{
			Template:            template,
			Name:                createName,
			Description:         createDescription,
			CloudRegionIdentity: region.Identity,
			VpcCidr:             createVpcCidr,
			SubnetCidrs:         append([]string(nil), createSubnetCidrs...),
			MachineType:         createMachineType,
			Labels:              shared.KeyValuePairsToMap(createLabels),
			Annotations:         shared.KeyValuePairsToMap(createAnnotations),
		}

		ql, err := client.QuickLaunch().CreateQuickLaunch(cmd.Context(), req)
		if err != nil {
			return fmt.Errorf("failed to create quick-launch job: %w", err)
		}

		printQuickLaunchDetails(ql, false)
		return nil
	},
}

func init() {
	QuickLaunchCmd.AddCommand(createCmd)
	createCmd.Flags().StringVar(&createName, "name", "", "Base name prefix for created resources (required)")
	createCmd.Flags().StringVar(&createDescription, "description", "", "Optional description")
	createCmd.Flags().StringVar(&createTemplate, "template", string(clientql.QuickLaunchTemplateVPC), "Template: vpc or kubernetes")
	createCmd.Flags().StringVar(&createRegion, "region", "", "Region identity, slug, or name (required)")
	createCmd.Flags().StringVar(&createVpcCidr, "vpc-cidr", "", "Optional VPC CIDR (default assigned by the API)")
	createCmd.Flags().StringSliceVar(&createSubnetCidrs, "subnet-cidr", nil, "Optional subnet CIDRs (repeatable)")
	createCmd.Flags().StringVar(&createMachineType, "machine-type", "", "Node pool machine type (kubernetes template only)")
	createCmd.Flags().StringSliceVar(&createLabels, "labels", nil, "Labels as key=value (repeatable)")
	createCmd.Flags().StringSliceVar(&createAnnotations, "annotations", nil, "Annotations as key=value (repeatable)")
	_ = createCmd.MarkFlagRequired("name")
	_ = createCmd.MarkFlagRequired("region")
	_ = createCmd.RegisterFlagCompletionFunc("region", completion.CompleteRegion)
	_ = createCmd.RegisterFlagCompletionFunc("machine-type", completion.CompleteMachineType)
	_ = createCmd.RegisterFlagCompletionFunc("template", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{
			string(clientql.QuickLaunchTemplateVPC) + "\tVPC with subnets and NAT",
			string(clientql.QuickLaunchTemplateKubernetes) + "\tVPC stack plus Kubernetes cluster",
		}, cobra.ShellCompDirectiveNoFileComp
	})
}
