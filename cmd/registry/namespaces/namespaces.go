package namespaces

import (
	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/cmd/registry/namespaces/configuration"
	"github.com/thalassa-cloud/cli/cmd/registry/namespaces/retention"
)

const NoHeaderKey = "no-header"

// NamespacesCmd manages container registry namespaces.
var NamespacesCmd = &cobra.Command{
	Use:     "namespaces",
	Aliases: []string{"namespace", "ns"},
	Short:   "Manage container registry namespaces",
	Example: "tcloud registry namespaces list\ntcloud registry namespaces create --namespace my-app --region eu-west-1",
}

func init() {
	NamespacesCmd.AddCommand(configuration.ConfigurationCmd)
	NamespacesCmd.AddCommand(retention.RetentionCmd)
}
