package fzf

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mattn/go-isatty"
)

// ErrSelectionCancelled is returned when the user aborts fzf (e.g. Esc / Ctrl-C).
var ErrSelectionCancelled = errors.New("selection cancelled")

// InteractiveChoiceOptions provides configuration for the interactive choice functionality.
type InteractiveChoiceOptions struct {
	// FzfArgs are additional arguments to pass to the fzf command
	FzfArgs []string
	// EnablePreview determines if previews should be enabled
	EnablePreview bool
}

// DefaultOptions returns the default options for interactive choice
func DefaultOptions() InteractiveChoiceOptions {
	return InteractiveChoiceOptions{
		FzfArgs:       []string{"--ansi"},
		EnablePreview: false,
	}
}

// InteractiveChoice runs fzf with the given command and returns the selected item.
// It uses default options for fzf configuration.
func InteractiveChoice(command string) (string, error) {
	return InteractiveChoiceWithOptions(command, DefaultOptions())
}

// InteractiveChoiceWithOptions runs fzf with the given command and options, returning the selected item.
func InteractiveChoiceWithOptions(command string, opts InteractiveChoiceOptions) (string, error) {
	// Check if fzf is available first to provide a better error message
	if !fzfInstalled() {
		return "", errors.New("fzf command not found, please install it first")
	}

	args := opts.FzfArgs
	if !opts.EnablePreview {
		args = append(args, "--no-preview")
	}

	cmd := exec.Command("fzf", args...)
	var out bytes.Buffer
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = &out

	cmd.Env = append(os.Environ(), fmt.Sprintf("FZF_DEFAULT_COMMAND=%s", command))
	if err := cmd.Run(); err != nil {
		if isFzfCancel(err) {
			return "", ErrSelectionCancelled
		}
		return "", fmt.Errorf("error running fzf: %w", err)
	}

	choice := strings.TrimSpace(out.String())
	if choice == "" {
		return "", ErrSelectionCancelled
	}

	return firstColumn(choice), nil
}

// InteractiveChoiceFromLines runs fzf against the provided lines and returns the selected item.
// The first tab-separated column of the chosen line is returned.
func InteractiveChoiceFromLines(lines []string) (string, error) {
	return InteractiveChoiceFromLinesWithOptions(lines, DefaultOptions())
}

// InteractiveChoiceFromLinesWithOptions runs fzf against the provided lines with the given options.
func InteractiveChoiceFromLinesWithOptions(lines []string, opts InteractiveChoiceOptions) (string, error) {
	if !fzfInstalled() {
		return "", errors.New("fzf command not found, please install it first")
	}
	if len(lines) == 0 {
		return "", errors.New("no options available")
	}

	args := opts.FzfArgs
	if !opts.EnablePreview {
		args = append(args, "--no-preview")
	}

	cmd := exec.Command("fzf", args...)
	var out bytes.Buffer
	cmd.Stdin = strings.NewReader(strings.Join(lines, "\n") + "\n")
	cmd.Stderr = os.Stderr
	cmd.Stdout = &out
	cmd.Env = withoutEnv(os.Environ(), "FZF_DEFAULT_COMMAND")

	if err := cmd.Run(); err != nil {
		if isFzfCancel(err) {
			return "", ErrSelectionCancelled
		}
		return "", fmt.Errorf("error running fzf: %w", err)
	}

	choice := strings.TrimSpace(out.String())
	if choice == "" {
		return "", ErrSelectionCancelled
	}

	return firstColumn(choice), nil
}

// isFzfCancel reports whether fzf exited because the user aborted selection.
// Exit 1 is abort/Esc; 130 is Ctrl-C.
func isFzfCancel(err error) bool {
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		return false
	}
	code := exitErr.ExitCode()
	return code == 1 || code == 130
}

func firstColumn(choice string) string {
	parts := strings.SplitN(choice, "\t", 2)
	return strings.TrimSpace(parts[0])
}

func withoutEnv(env []string, key string) []string {
	prefix := key + "="
	filtered := make([]string, 0, len(env))
	for _, item := range env {
		if strings.HasPrefix(item, prefix) {
			continue
		}
		filtered = append(filtered, item)
	}
	return filtered
}

// IsInteractiveMode determines if we can use fzf for interactive selection.
// It checks if:
// 1. TC_IGNORE_FZF environment variable is not set
// 2. stdout is a terminal
// 3. fzf is installed
func IsInteractiveMode(stdout *os.File) bool {
	value := os.Getenv("TC_IGNORE_FZF")
	return value == "" && isTerminal(stdout) && fzfInstalled()
}

// isTerminal checks if the given file descriptor is a terminal.
func isTerminal(fd *os.File) bool {
	return isatty.IsTerminal(fd.Fd())
}

// fzfInstalled checks if the fzf command is available in the PATH.
// It returns true if fzf is installed, false otherwise.
func fzfInstalled() bool {
	_, err := exec.LookPath("fzf")
	return err == nil
}
