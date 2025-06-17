package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/cmd/context"
	"github.com/thalassa-cloud/cli/cmd/iaas/compute"
	"github.com/thalassa-cloud/cli/cmd/iaas/networking"
	"github.com/thalassa-cloud/cli/cmd/iaas/regions"
	"github.com/thalassa-cloud/cli/cmd/iaas/storage"
	"github.com/thalassa-cloud/cli/cmd/kubernetes"
	"github.com/thalassa-cloud/cli/cmd/me"
	"github.com/thalassa-cloud/cli/cmd/version"
	"github.com/thalassa-cloud/cli/internal/config/contextstate"
)

var RootCmd = &cobra.Command{
	Use:   "tcloud",
	Short: "A CLI for working with the Thalassa Cloud Platform",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		handleExecutionError(err)
		os.Exit(1)
	}
}

func handleExecutionError(err error) {
	switch err {
	default:
		_, _ = fmt.Fprintf(os.Stderr, "failed: %v\n", err)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&contextstate.OrganisationFlag, "organisation", "O", "", "Organisation slug or identity (overrides context)")
	RootCmd.PersistentFlags().StringVarP(&contextstate.ContextFlag, "context", "c", "", "Context name")
	RootCmd.PersistentFlags().StringVar(&contextstate.EndpointFlag, "api", "", "API endpoint (overrides context)")
	RootCmd.PersistentFlags().StringVar(&contextstate.PersonalAccessTokenFlag, "token", "", "Personal access token (overrides context)")
	RootCmd.PersistentFlags().StringVar(&contextstate.OidcClientIDFlag, "client-id", "", "OIDC client ID for OIDC authentication (overrides context)")
	RootCmd.PersistentFlags().StringVar(&contextstate.OidcClientSecretFlag, "client-secret", "", "OIDC client secret for OIDC authentication (overrides context)")
	RootCmd.PersistentFlags().BoolVar(&contextstate.DebugFlag, "debug", false, "Debug mode")

	RootCmd.AddCommand(context.ContextCmd)
	RootCmd.AddCommand(version.VersionCmd)

	RootCmd.AddCommand(regions.RegionsCmd)

	RootCmd.AddCommand(networking.NetworkingCmd)
	RootCmd.AddCommand(storage.StorageCmd)
	RootCmd.AddCommand(compute.ComputeCmd)

	RootCmd.AddCommand(kubernetes.KubernetesCmd)
	RootCmd.AddCommand(me.MeCmd)

	cobra.OnInitialize(contextstate.Init)
}
