package quotas

import "github.com/spf13/cobra"

// QuotasCmd manages organisation resource quotas.
var QuotasCmd = &cobra.Command{
	Use:     "quotas",
	Aliases: []string{"quota"},
	Short:   "View and request changes to organisation resource quotas",
	Long: `List organisation quotas, inspect current usage, and submit increase requests.
Quotas apply per organisation in your current context.`,
	Example: `  # List all quotas for the current organisation
  tcloud quotas list

  # Show a single quota including pending increase requests
  tcloud quotas get machines

  # Request a higher limit
  tcloud quotas request-increase machines --new-max-usage 50 --reason "Growing production fleet"`,
}
