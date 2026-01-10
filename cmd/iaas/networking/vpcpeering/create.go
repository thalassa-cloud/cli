package vpcpeering

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

const (
	CreateFlagName                 = "name"
	CreateFlagDescription          = "description"
	CreateFlagRequesterVpc         = "requester-vpc"
	CreateFlagAccepterVpc          = "accepter-vpc"
	CreateFlagAccepterOrganisation = "accepter-organisation"
	CreateFlagAutoAccept           = "auto-accept"
	CreateFlagLabels               = "labels"
	CreateFlagAnnotations          = "annotations"
)

var (
	createName                         string
	createDescription                  string
	createRequesterVpcIdentity         string
	createAccepterVpcIdentity          string
	createAccepterOrganisationIdentity string
	createAutoAccept                   bool
	createLabels                       []string
	createAnnotations                  []string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a VPC peering connection",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		if createName == "" {
			return fmt.Errorf("name is required")
		}
		if createRequesterVpcIdentity == "" {
			return fmt.Errorf("requester-vpc is required")
		}
		if createAccepterVpcIdentity == "" {
			return fmt.Errorf("accepter-vpc is required")
		}
		if createAccepterOrganisationIdentity == "" {
			return fmt.Errorf("accepter-organisation is required")
		}

		// Resolve requester VPC
		requesterVpc, err := client.IaaS().GetVpc(cmd.Context(), createRequesterVpcIdentity)
		if err != nil {
			if tcclient.IsNotFound(err) {
				vpcs, err := client.IaaS().ListVpcs(cmd.Context(), &iaas.ListVpcsRequest{})
				if err != nil {
					return err
				}
				for _, v := range vpcs {
					if v.Slug == createRequesterVpcIdentity || v.Name == createRequesterVpcIdentity {
						requesterVpc = &v
						break
					}
				}
				if requesterVpc == nil {
					return fmt.Errorf("requester VPC not found: %s", createRequesterVpcIdentity)
				}
			} else {
				return fmt.Errorf("failed to get requester VPC: %w", err)
			}
		}

		// Resolve accepter VPC
		accepterVpc, err := client.IaaS().GetVpc(cmd.Context(), createAccepterVpcIdentity)
		if err != nil {
			if tcclient.IsNotFound(err) {
				vpcs, err := client.IaaS().ListVpcs(cmd.Context(), &iaas.ListVpcsRequest{})
				if err != nil {
					return err
				}
				for _, v := range vpcs {
					if v.Slug == createAccepterVpcIdentity || v.Name == createAccepterVpcIdentity {
						accepterVpc = &v
						break
					}
				}
				if accepterVpc == nil {
					return fmt.Errorf("accepter VPC not found: %s", createAccepterVpcIdentity)
				}
			} else {
				return fmt.Errorf("failed to get accepter VPC: %w", err)
			}
		}

		// Parse labels from key=value format
		labels := make(map[string]string)
		for _, label := range createLabels {
			parts := strings.SplitN(label, "=", 2)
			if len(parts) == 2 {
				labels[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}

		// Parse annotations from key=value format
		annotations := make(map[string]string)
		for _, annotation := range createAnnotations {
			parts := strings.SplitN(annotation, "=", 2)
			if len(parts) == 2 {
				annotations[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}

		req := iaas.CreateVpcPeeringConnectionRequest{
			Name:                         createName,
			Description:                  createDescription,
			RequesterVpcIdentity:         requesterVpc.Identity,
			AccepterVpcIdentity:          accepterVpc.Identity,
			AccepterOrganisationIdentity: createAccepterOrganisationIdentity,
			AutoAccept:                   createAutoAccept,
			Labels:                       labels,
			Annotations:                  annotations,
		}

		connection, err := client.IaaS().CreateVpcPeeringConnection(cmd.Context(), req)
		if err != nil {
			return fmt.Errorf("failed to create VPC peering connection: %w", err)
		}

		body := make([][]string, 0, 1)
		body = append(body, []string{
			connection.Identity,
			connection.Name,
			string(connection.Status),
			formattime.FormatTime(connection.CreatedAt.Local(), false),
		})
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Status", "Age"}, body)
		}
		return nil
	},
}

func init() {
	VpcPeeringCmd.AddCommand(createCmd)

	createCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	createCmd.Flags().StringVar(&createName, CreateFlagName, "", "Name of the VPC peering connection")
	createCmd.Flags().StringVar(&createDescription, CreateFlagDescription, "", "Description of the VPC peering connection")
	createCmd.Flags().StringVar(&createRequesterVpcIdentity, CreateFlagRequesterVpc, "", "Identity of the requester VPC")
	createCmd.Flags().StringVar(&createAccepterVpcIdentity, CreateFlagAccepterVpc, "", "Identity of the accepter VPC")
	createCmd.Flags().StringVar(&createAccepterOrganisationIdentity, CreateFlagAccepterOrganisation, "", "Identity of the accepter organisation")
	createCmd.Flags().BoolVar(&createAutoAccept, CreateFlagAutoAccept, false, "Automatically accept the peering connection (only if requester and accepter are in same region and organisation)")
	createCmd.Flags().StringSliceVar(&createLabels, CreateFlagLabels, []string{}, "Labels in key=value format")
	createCmd.Flags().StringSliceVar(&createAnnotations, CreateFlagAnnotations, []string{}, "Annotations in key=value format")

	createCmd.MarkFlagRequired(CreateFlagName)
	createCmd.MarkFlagRequired(CreateFlagRequesterVpc)
	createCmd.MarkFlagRequired(CreateFlagAccepterVpc)
	createCmd.MarkFlagRequired(CreateFlagAccepterOrganisation)
}
