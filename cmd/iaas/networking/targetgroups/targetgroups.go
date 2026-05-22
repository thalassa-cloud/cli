package targetgroups

import "github.com/spf13/cobra"

const NoHeaderKey = "no-header"

// TargetGroupsCmd represents the target groups command.
var TargetGroupsCmd = &cobra.Command{
	Use:     "target-groups",
	Aliases: []string{"targetgroups", "targetgroup", "target-group", "tg"},
	Short:   "Manage load balancer target groups",
	Long:    "Manage load balancer target groups, attachments, and backend configuration.",
	Example: "tcloud networking target-groups list\ntcloud networking target-groups create --name web --vpc vpc-123 --port 8080 --protocol http\ntcloud networking target-groups attach tg-123 --server machine-123",
}
