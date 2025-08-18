package securitygroups

import (
	"github.com/spf13/cobra"
)

// SecurityGroupsCmd represents the security groups command
var SecurityGroupsCmd = &cobra.Command{
	Use:     "security-groups",
	Aliases: []string{"securitygroups", "securitygroup", "sg"},
	Short:   "Manage security groups",
	Long:    "Manage security groups and their rules within the Thalassa Cloud Platform",
	Example: "tcloud networking security-groups list\ntcloud networking security-groups create --name my-sg --vpc vpc-123\ntcloud networking security-groups delete sg-456",
}

func init() {
	// Add subcommands here as they are implemented
}
