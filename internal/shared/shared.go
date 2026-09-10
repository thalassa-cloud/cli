package shared

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/thalassa-cloud/client-go/pkg/base"
)

const (
	NoHeaderKey = "no-header"
	// ForceKey is the flag name for skipping destructive confirmation prompts.
	ForceKey = "force"
)

// PromptDestructiveUnlessForce prompts for typing "yes" unless force is true.
// On abort (anything other than "yes"), returns proceed=false and err=nil.
func PromptDestructiveUnlessForce(force bool, summary string) (proceed bool, err error) {
	if force {
		return true, nil
	}
	fmt.Print(summary)
	if summary != "" && !strings.HasSuffix(summary, "\n") {
		fmt.Println()
	}
	fmt.Print("Enter 'yes' to confirm: ")
	var input string
	if _, scanErr := fmt.Scanln(&input); scanErr != nil {
		return false, fmt.Errorf("read confirmation: %w", scanErr)
	}
	if strings.TrimSpace(input) != "yes" {
		fmt.Println("Aborted")
		return false, nil
	}
	return true, nil
}

// PromptYesNoUnlessForce prompts for y/yes unless force is true.
// On decline, returns proceed=false and err=nil without printing "Aborted".
func PromptYesNoUnlessForce(force bool, question string) (proceed bool, err error) {
	if force {
		return true, nil
	}
	fmt.Print(question)
	var input string
	if _, scanErr := fmt.Scanln(&input); scanErr != nil {
		return false, fmt.Errorf("read confirmation: %w", scanErr)
	}
	switch strings.ToLower(strings.TrimSpace(input)) {
	case "y", "yes":
		return true, nil
	default:
		return false, nil
	}
}

// PromptString reads a line from stdin. Empty input returns defaultValue.
func PromptString(prompt, defaultValue string) (string, error) {
	return PromptStringFrom(os.Stdin, prompt, defaultValue)
}

// PromptStringFrom reads a line from r. Empty input returns defaultValue.
func PromptStringFrom(r io.Reader, prompt, defaultValue string) (string, error) {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", prompt, defaultValue)
	} else {
		fmt.Printf("%s: ", prompt)
	}
	scanner := bufio.NewScanner(r)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("read input: %w", err)
		}
		return "", fmt.Errorf("read input: unexpected EOF")
	}
	value := strings.TrimSpace(scanner.Text())
	if value == "" {
		return defaultValue, nil
	}
	return value, nil
}

func KeyValuePairsToMap(pairs []string) map[string]string {
	out := make(map[string]string)
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			out[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return out
}

func UserDisplay(u base.AppUser) string {
	if u.Email != "" {
		return u.Email
	}
	if u.Name != "" {
		return u.Name
	}
	return u.Subject
}

func UserPtrDisplay(u *base.AppUser) string {
	if u == nil {
		return ""
	}
	return UserDisplay(*u)
}
