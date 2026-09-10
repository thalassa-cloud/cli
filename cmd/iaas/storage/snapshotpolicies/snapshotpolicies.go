package snapshotpolicies

import (
	"github.com/spf13/cobra"
)

const NoHeaderKey = "no-header"

var noHeader bool

// SnapshotPoliciesCmd manages snapshot policies.
var SnapshotPoliciesCmd = &cobra.Command{
	Use:     "snapshot-policies",
	Aliases: []string{"snapshot-policy", "sp"},
	Short:   "Manage snapshot policies",
	Long:    "Manage automated volume snapshot policies within the Thalassa Cloud Platform.",
	Example: "tcloud storage snapshot-policies list\ntcloud storage snapshot-policies create --name daily --region nl-ams --schedule '0 2 * * *' --ttl 168h --target-type selector --selector backup=true",
}

func init() {
}
