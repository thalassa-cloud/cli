package tfs

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/table"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	tcclient "github.com/thalassa-cloud/client-go/pkg/client"
)

var (
	viewShowExactTime bool
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:               "view",
	Short:             "View TFS instance details",
	Long:              "View detailed information about a TFS instance.",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteTfsInstanceID,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		instanceIdentity := args[0]

		instance, err := client.Tfs().GetTfsInstance(cmd.Context(), instanceIdentity)
		if err != nil {
			if tcclient.IsNotFound(err) {
				return fmt.Errorf("TFS instance not found: %s", instanceIdentity)
			}
			return fmt.Errorf("failed to get TFS instance: %w", err)
		}

		regionName := ""
		if instance.Region != nil {
			regionName = instance.Region.Name
		}

		body := [][]string{
			{"ID", instance.Identity},
			{"Name", instance.Name},
			{"Status", string(instance.Status)},
			{"Region", regionName},
			{"Created", formattime.FormatTime(instance.CreatedAt.Local(), viewShowExactTime)},
		}

		if instance.Description != nil && *instance.Description != "" {
			body = append(body, []string{"Description", *instance.Description})
		}

		if len(instance.Labels) > 0 {
			labelStrs := []string{}
			for k, v := range instance.Labels {
				labelStrs = append(labelStrs, k+"="+v)
			}
			sort.Strings(labelStrs)
			body = append(body, []string{"Labels", strings.Join(labelStrs, ", ")})
		}

		if len(instance.Annotations) > 0 {
			annotationStrs := []string{}
			for k, v := range instance.Annotations {
				annotationStrs = append(annotationStrs, k+"="+v)
			}
			sort.Strings(annotationStrs)
			body = append(body, []string{"Annotations", strings.Join(annotationStrs, ", ")})
		}

		if instance.DeleteProtection {
			body = append(body, []string{"Delete Protection", "enabled"})
		}

		if noHeader {
			table.Print(nil, body)
		} else {
			table.Print([]string{"Field", "Value"}, body)
		}

		return nil
	},
}

func init() {
	TfsCmd.AddCommand(viewCmd)

	viewCmd.Flags().BoolVar(&noHeader, NoHeaderKey, false, "Do not print the header")
	viewCmd.Flags().BoolVar(&viewShowExactTime, "exact-time", false, "Show exact time instead of relative time")
}
