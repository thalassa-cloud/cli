package quicklaunch

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/completion"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	clientql "github.com/thalassa-cloud/client-go/quicklaunch"
)

var (
	deleteForce   bool
	deleteCascade string
)

var deleteCmd = &cobra.Command{
	Use:     "delete <identity>",
	Aliases: []string{"rm", "remove"},
	Short:   "Delete a quick-launch job",
	Long: `Delete a quick-launch job.

Use --cascade Delete to also delete provisioned resources, or --cascade Orphan
to keep resources and only remove the quick-launch record.`,
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: completion.CompleteQuickLaunchIdentity,
	RunE: func(cmd *cobra.Command, args []string) error {
		cascade := clientql.QuickLaunchCascade(strings.TrimSpace(deleteCascade))
		switch cascade {
		case "", clientql.QuickLaunchCascadeDelete, clientql.QuickLaunchCascadeOrphan:
		default:
			return fmt.Errorf("--cascade must be Delete or Orphan (got %q)", deleteCascade)
		}

		client, err := thalassaclient.GetThalassaClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}

		summary := fmt.Sprintf("Are you sure you want to delete this quick-launch job?\n  Job: %s\n", args[0])
		if cascade != "" {
			summary += fmt.Sprintf("  Cascade: %s\n", cascade)
		}
		ok, err := shared.PromptDestructiveUnlessForce(deleteForce, summary)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}

		if err := client.QuickLaunch().DeleteQuickLaunch(cmd.Context(), args[0], cascade); err != nil {
			return fmt.Errorf("failed to delete quick-launch job: %w", err)
		}
		fmt.Printf("Deleted quick-launch job %s\n", args[0])
		return nil
	},
}

func init() {
	QuickLaunchCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVar(&deleteForce, shared.ForceKey, false, "Skip the confirmation prompt and delete")
	deleteCmd.Flags().StringVar(&deleteCascade, "cascade", "", "Child resource handling: Delete or Orphan")
	_ = deleteCmd.RegisterFlagCompletionFunc("cascade", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{
			string(clientql.QuickLaunchCascadeDelete) + "\tDelete provisioned resources",
			string(clientql.QuickLaunchCascadeOrphan) + "\tKeep resources, remove job only",
		}, cobra.ShellCompDirectiveNoFileComp
	})
}
