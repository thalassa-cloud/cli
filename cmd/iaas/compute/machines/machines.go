package machines

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/fzf"
)

// MachinesCmd represents the machines command
var MachinesCmd = &cobra.Command{
	Use:     "machines",
	Aliases: []string{"machine", "vm", "vms", "virtualmachines", "instances", "instance"},
	Short:   "Manage virtual machine instances",
	Long:    "Manage virtual machine instances in the Thalassa Cloud Platform",
}

func init() {
}

func getSelectedMachine(args []string) (string, error) {
	if len(args) == 0 && fzf.IsInteractiveMode(os.Stdout) {
		command := fmt.Sprintf("%s compute machines ls --no-header", os.Args[0])
		return fzf.InteractiveChoice(command)
	} else if len(args) == 1 {
		return args[0], nil
	} else {
		return "", errors.New("invalid machine")
	}
}
