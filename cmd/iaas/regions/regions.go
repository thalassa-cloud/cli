package regions

import (
	"github.com/spf13/cobra"
)

var RegionsCmd = &cobra.Command{
	Use:     "regions",
	Aliases: []string{"region"},
	Short:   "Thalassa Cloud Platform Regions",
}

func init() {
}
