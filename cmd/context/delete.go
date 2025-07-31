package context

import (
	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/config/contextstate"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a context",
	Long:    "Delete a context from the config",
	Example: "tcloud context delete <context>",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return contextstate.GlobalConfigManager().RemoveContext(args[0])
	},
}

func init() {
	ContextCmd.AddCommand(deleteCmd)
	ContextCmd.AddCommand(deleteUserCmd)
	ContextCmd.AddCommand(deleteServerCmd)
}

// deleteUser & deleteServer are helpers for the delete command
var deleteUserCmd = &cobra.Command{
	Use:     "delete-user",
	Short:   "Delete a user",
	Long:    "Delete a user from the config",
	Example: "tcloud context delete-user <user>",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return contextstate.GlobalConfigManager().RemoveContextUser(args[0])
	},
}

var deleteServerCmd = &cobra.Command{
	Use:     "delete-server",
	Short:   "Delete a server",
	Long:    "Delete a server from the config",
	Example: "tcloud context delete-server <server>",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return contextstate.GlobalConfigManager().RemoveContextServer(args[0])
	},
}
