package records

import (
	"github.com/spf13/cobra"
)

// RecordsCmd manages DNS records.
var RecordsCmd = &cobra.Command{
	Use:     "records",
	Aliases: []string{"record", "rr"},
	Short:   "Manage DNS records",
}

var (
	noHeader      bool
	showExactTime bool
	zoneIdentity  string
)
