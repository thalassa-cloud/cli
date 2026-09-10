package zones

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
)

var viewCmd = &cobra.Command{
	Use:               "view <zone>",
	Aliases:           []string{"get", "show"},
	Short:             "View a DNS zone",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteDnsZoneIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		zone, err := client.DNS().GetZone(cmd.Context(), args[0])
		if err != nil {
			return fmt.Errorf("failed to get DNS zone: %w", err)
		}

		body := [][]string{{
			zone.Identity,
			zone.Name,
			zone.Slug,
			zone.Description,
			formatMap(zone.Labels),
			formatMap(zone.Annotations),
			formattime.FormatTime(zone.CreatedAt.Local(), showExactTime),
		}}
		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"ID", "Name", "Slug", "Description", "Labels", "Annotations", "Created"}, body)
		}
		return nil
	},
}

func formatMap(m map[string]string) string {
	if len(m) == 0 {
		return "-"
	}
	parts := make([]string, 0, len(m))
	for k, v := range m {
		parts = append(parts, k+"="+v)
	}
	sort.Strings(parts)
	return strings.Join(parts, ",")
}

func init() {
	ZonesCmd.AddCommand(viewCmd)
	viewCmd.Flags().BoolVar(&noHeader, shared.NoHeaderKey, false, "Do not print table headers")
	viewCmd.Flags().BoolVar(&showExactTime, "exact-time", false, "Show full timestamps instead of relative time")
}
