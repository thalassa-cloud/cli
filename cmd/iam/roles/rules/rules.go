package rules

import "github.com/spf13/cobra"

// RulesCmd manages permission rules on organisation roles.
var RulesCmd = &cobra.Command{
	Use:   "rules",
	Short: "Permission rules on a role",
}
