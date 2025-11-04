package audit

import (
	"github.com/spf13/cobra"
)

var AuditCmd = &cobra.Command{
	Use:     "audit",
	Aliases: []string{"audit-logs", "auditlogs"},
	Short:   "Manage organisation audit logs",
	Long:    "Audit commands to manage and export organisation audit logs for compliance",
}

func init() {
}
