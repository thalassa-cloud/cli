package configuration

import "github.com/spf13/cobra"

const NamespaceFlag = "namespace"

// ConfigurationCmd manages namespace configuration.
var ConfigurationCmd = &cobra.Command{
	Use:     "configuration",
	Aliases: []string{"config", "cfg"},
	Short:   "Manage namespace configuration",
	Long:    "Manage visibility and retention policy configuration for a container registry namespace.",
	Example: "tcloud registry namespaces configuration view --namespace crns-123\ntcloud registry namespaces configuration create --namespace crns-123 --visibility private",
}
