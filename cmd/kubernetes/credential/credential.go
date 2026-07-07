package credential

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/config/contextstate"
	"github.com/thalassa-cloud/cli/internal/kubernetes/auth"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var (
	clusterIdentity string
	organisation    string
	project         string
	contextName     string
)

var CredentialCmd = &cobra.Command{
	Use:    "credential",
	Short:  "Provide Kubernetes credentials to kubectl",
	Hidden: true,
	Args:   cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if clusterIdentity == "" {
			return fmt.Errorf("cluster identity is required")
		}
		if organisation == "" {
			return fmt.Errorf("organisation is required")
		}
		if contextName != "" {
			if err := contextstate.Set(contextName); err != nil {
				return err
			}
		}

		client, err := thalassaclient.GetThalassaClientWithScope(thalassaclient.Scope{
			Organisation: organisation,
			Project:      project,
		})
		if err != nil {
			return err
		}

		return auth.WriteExecCredential(cmd.Context(), client, clusterIdentity, os.Stdout)
	},
}

func init() {
	CredentialCmd.Flags().StringVar(&clusterIdentity, "cluster", "", "Kubernetes cluster identity")
	CredentialCmd.Flags().StringVar(&organisation, "organisation", "", "Organisation slug or identity")
	CredentialCmd.Flags().StringVar(&project, "project", "", "Project identity or slug")
	CredentialCmd.Flags().StringVar(&contextName, "context", "", "CLI context name")
	_ = CredentialCmd.MarkFlagRequired("cluster")
	_ = CredentialCmd.MarkFlagRequired("organisation")
	CredentialCmd.SilenceUsage = true
}
