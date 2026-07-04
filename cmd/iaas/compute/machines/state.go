package machines

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/thalassa-cloud/cli/internal/thalassaclient"
	"github.com/thalassa-cloud/client-go/iaas"
)

type machineStateChange struct {
	alreadyState    iaas.MachineState
	alreadyMessage  string
	change          func(context.Context, *iaas.Client, string) error
	progressMessage string
	targetState     iaas.MachineState
	doneMessage     string
}

func runMachineStateChange(cmd *cobra.Command, args []string, change machineStateChange) error {
	thalassaClient, err := thalassaclient.GetThalassaClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	machineIdentity, err := getSelectedMachine(args)
	if err != nil {
		return err
	}

	iaasClient := thalassaClient.IaaS()
	machine, err := iaasClient.GetMachine(cmd.Context(), machineIdentity)
	if err != nil {
		return err
	}

	if machine.Status.Status == string(change.alreadyState) {
		fmt.Println(change.alreadyMessage)
		return nil
	}

	if err := change.change(cmd.Context(), iaasClient, machine.Identity); err != nil {
		return err
	}
	fmt.Println(change.progressMessage)

	if !wait {
		return nil
	}

	for {
		machine, err = iaasClient.GetMachine(cmd.Context(), machine.Identity)
		if err != nil {
			return err
		}
		if machine.Status.Status == string(change.targetState) {
			break
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Println(change.doneMessage)
	return nil
}
