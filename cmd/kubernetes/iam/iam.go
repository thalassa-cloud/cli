package iam

import (
	"github.com/spf13/cobra"
	"github.com/thalassa-cloud/cli/cmd/kubernetes/iam/roles"
)

// IamCmd manages Kubernetes cluster IAM (RBAC roles and bindings).
var IamCmd = &cobra.Command{
	Use:   "iam",
	Short: "Kubernetes cluster IAM roles and bindings",
	Long: `Manage Kubernetes cluster IAM roles, permission rules, and role bindings.
These roles control access within your organisation's Kubernetes clusters on Thalassa Cloud.`,
}

func init() {
	IamCmd.AddCommand(roles.RolesCmd)
}
