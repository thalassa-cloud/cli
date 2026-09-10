package secrets

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	clientsecrets "github.com/thalassa-cloud/client-go/secrets"

	"github.com/thalassa-cloud/cli/internal/formattime"
	"github.com/thalassa-cloud/cli/internal/fzf"
	"github.com/thalassa-cloud/cli/internal/shared"
	"github.com/thalassa-cloud/cli/internal/table"
)

const (
	browseSelectUp = ".."

	secretActionBack     = "back"
	secretActionVersions = "versions"
	secretActionReveal   = "reveal"
	secretActionQuit     = "quit"
)

type secretsBrowser interface {
	BrowseSecrets(ctx context.Context, region, path string) (*clientsecrets.BrowseSecretsResponse, error)
	GetSecret(ctx context.Context, region, path string, includeVersions bool) (*clientsecrets.Secret, error)
	GetSecretValue(ctx context.Context, region, path string, version *int) (*clientsecrets.GetSecretValueResponse, error)
}

func runInteractiveBrowse(ctx context.Context, out io.Writer, api secretsBrowser, region, startPath string) error {
	current, err := clientsecrets.NormalizePath(startPath)
	if err != nil {
		return err
	}

	for {
		result, err := api.BrowseSecrets(ctx, region, current)
		if err != nil {
			return fmt.Errorf("failed to browse secrets: %w", err)
		}
		if result.Path != "" {
			current = result.Path
		}

		lines := buildBrowseSelectionLines(current, result)
		if len(lines) == 0 {
			_, _ = fmt.Fprintf(out, "No prefixes or secrets at %s\n", current)
			return nil
		}

		choice, err := fzf.InteractiveChoiceFromLinesWithOptions(lines, fzf.InteractiveChoiceOptions{
			FzfArgs: []string{
				"--ansi",
				"--header=region " + region + "  path " + current + "  (Esc to quit)",
				"--prompt=secrets> ",
			},
		})
		if err != nil {
			if errors.Is(err, fzf.ErrSelectionCancelled) {
				return nil
			}
			return err
		}

		switch {
		case choice == browseSelectUp:
			current = parentBrowsePath(current)
		case isBrowsePrefix(choice, result.Prefixes):
			current = choice
		default:
			quit, openErr := openSecretInteractive(ctx, out, api, region, choice)
			if openErr != nil {
				return openErr
			}
			if quit {
				return nil
			}
		}
	}
}

func buildBrowseSelectionLines(current string, result *clientsecrets.BrowseSecretsResponse) []string {
	if result == nil {
		return nil
	}

	lines := make([]string, 0, 1+len(result.Prefixes)+len(result.Secrets))
	if current != "/" {
		lines = append(lines, browseSelectUp+"\t..\tgo up")
	}

	for _, prefix := range result.Prefixes {
		label := path.Base(strings.TrimSuffix(prefix, "/"))
		if label == "" || label == "." {
			label = prefix
		}
		lines = append(lines, prefix+"\t"+label+"/\tprefix")
	}

	for _, secret := range result.Secrets {
		label := path.Base(secret.Path)
		if label == "" || label == "/" || label == "." {
			label = secret.Path
		}
		lines = append(lines, fmt.Sprintf("%s\t%s\tsecret · v%d · %s",
			secret.Path,
			label,
			secret.CurrentVersion,
			formattime.FormatTime(secret.UpdatedAt.Local(), showExactTime),
		))
	}
	return lines
}

func isBrowsePrefix(choice string, prefixes []string) bool {
	if strings.HasSuffix(choice, "/") {
		return true
	}
	for _, prefix := range prefixes {
		if prefix == choice {
			return true
		}
	}
	return false
}

func parentBrowsePath(p string) string {
	normalized := strings.TrimSpace(p)
	if normalized == "" || normalized == "/" {
		return "/"
	}
	normalized = strings.TrimSuffix(normalized, "/")
	idx := strings.LastIndex(normalized, "/")
	if idx <= 0 {
		return "/"
	}
	return normalized[:idx+1]
}

func openSecretInteractive(ctx context.Context, out io.Writer, api secretsBrowser, region, secretPath string) (quit bool, err error) {
	secret, err := api.GetSecret(ctx, region, secretPath, false)
	if err != nil {
		return false, fmt.Errorf("failed to get secret: %w", err)
	}
	printSecretMetadata(out, secret)
	if err := pauseForInteractiveView(out); err != nil {
		return false, err
	}

	for {
		action, err := fzf.InteractiveChoiceFromLinesWithOptions([]string{
			secretActionBack + "\tBack to browse",
			secretActionVersions + "\tShow version history",
			secretActionReveal + "\tReveal secret value (sensitive)",
			secretActionQuit + "\tQuit",
		}, fzf.InteractiveChoiceOptions{
			FzfArgs: []string{
				"--ansi",
				"--header=secret " + secretPath,
				"--prompt=secret> ",
			},
		})
		if err != nil {
			if errors.Is(err, fzf.ErrSelectionCancelled) {
				return false, nil
			}
			return false, err
		}

		switch action {
		case secretActionBack:
			return false, nil
		case secretActionQuit:
			return true, nil
		case secretActionVersions:
			if err := printSecretVersions(ctx, out, api, region, secretPath); err != nil {
				return false, err
			}
			if err := pauseForInteractiveView(out); err != nil {
				return false, err
			}
		case secretActionReveal:
			if err := revealSecretValue(ctx, out, api, region, secretPath); err != nil {
				return false, err
			}
			if err := pauseForInteractiveView(out); err != nil {
				return false, err
			}
		default:
			return false, nil
		}
	}
}

// pauseForInteractiveView keeps printed output visible until the user continues,
// because the next fzf screen would otherwise hide it.
func pauseForInteractiveView(out io.Writer) error {
	_, _ = fmt.Fprint(out, "\nPress Enter to continue...")
	_, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
	_, _ = fmt.Fprintln(out)
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("wait for continue: %w", err)
	}
	return nil
}

func printSecretMetadata(out io.Writer, secret *clientsecrets.Secret) {
	if secret == nil {
		return
	}

	kmsKey := "-"
	if secret.KmsKey != nil {
		kmsKey = secret.KmsKey.Identity
		if kmsKey == "" {
			kmsKey = secret.KmsKey.Name
		}
	}

	body := [][]string{
		{"Path", secret.Path},
		{"Description", secret.Description},
		{"Current Version", fmt.Sprintf("%d", secret.CurrentVersion)},
		{"KMS Key", kmsKey},
		{"Created", formattime.FormatTime(secret.CreatedAt.Local(), showExactTime)},
		{"Updated", formattime.FormatTime(secret.UpdatedAt.Local(), showExactTime)},
	}
	_, _ = fmt.Fprintln(out)
	if noHeader {
		table.PrintWithWriter(out, nil, body)
	} else {
		table.PrintWithWriter(out, []string{"Field", "Value"}, body)
	}
}

func printSecretVersions(ctx context.Context, out io.Writer, api secretsBrowser, region, secretPath string) error {
	secret, err := api.GetSecret(ctx, region, secretPath, true)
	if err != nil {
		return fmt.Errorf("failed to get secret versions: %w", err)
	}
	if len(secret.Versions) == 0 {
		_, _ = fmt.Fprintln(out, "No versions returned.")
		return nil
	}
	printVersionTable(out, secret.Versions)
	return nil
}

func printVersionTable(out io.Writer, versions []clientsecrets.SecretVersion) {
	_, _ = fmt.Fprintln(out)
	versionBody := make([][]string, 0, len(versions))
	for _, v := range versions {
		destroyed := "-"
		if v.DestroyedAt != nil {
			destroyed = formattime.FormatTime(v.DestroyedAt.Local(), showExactTime)
		}
		versionBody = append(versionBody, []string{
			fmt.Sprintf("%d", v.Version),
			v.Status,
			formattime.FormatTime(v.CreatedAt.Local(), showExactTime),
			destroyed,
		})
	}
	if noHeader {
		table.PrintWithWriter(out, nil, versionBody)
	} else {
		table.PrintWithWriter(out, []string{"Version", "Status", "Created", "Destroyed"}, versionBody)
	}
}

func revealSecretValue(ctx context.Context, out io.Writer, api secretsBrowser, region, secretPath string) error {
	proceed, err := shared.PromptDestructiveUnlessForce(false, fmt.Sprintf(
		"Reveal secret material for %q to this terminal?\nThis prints sensitive data. Prefer piping get-value to a file when possible.",
		secretPath,
	))
	if err != nil {
		return err
	}
	if !proceed {
		return nil
	}

	result, err := api.GetSecretValue(ctx, region, secretPath, nil)
	if err != nil {
		return fmt.Errorf("failed to get secret value: %w", err)
	}

	_, _ = fmt.Fprintln(out)
	switch {
	case result.SecretString != "":
		plaintext, decodeErr := clientsecrets.DecodeBytes("secretString", result.SecretString)
		if decodeErr != nil {
			if _, err := fmt.Fprint(out, result.SecretString); err != nil {
				return err
			}
			if !strings.HasSuffix(result.SecretString, "\n") {
				_, _ = fmt.Fprintln(out)
			}
			return nil
		}
		if _, err := fmt.Fprint(out, string(plaintext)); err != nil {
			return err
		}
		if len(plaintext) == 0 || plaintext[len(plaintext)-1] != '\n' {
			_, _ = fmt.Fprintln(out)
		}
	case len(result.SecretKeyValues) > 0:
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		if err := enc.Encode(result.SecretKeyValues); err != nil {
			return err
		}
	default:
		_, _ = fmt.Fprintln(out, "{}")
	}
	return nil
}
