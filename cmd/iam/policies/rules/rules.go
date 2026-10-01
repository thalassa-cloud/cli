package rules

import "github.com/spf13/cobra"

// RulesCmd manages permission rules on IAM policies.
var RulesCmd = &cobra.Command{
	Use:   "rules",
	Short: "Permission rules on a policy",
}
