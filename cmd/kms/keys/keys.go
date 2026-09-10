package keys

import (
	"github.com/spf13/cobra"
)

// KeysCmd manages KMS keys.
var KeysCmd = &cobra.Command{
	Use:     "keys",
	Aliases: []string{"key", "k"},
	Short:   "Manage KMS keys",
}

var (
	noHeader      bool
	showExactTime bool
	region        string
)
