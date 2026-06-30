package shared

import (
	"fmt"
	"strings"

	"github.com/thalassa-cloud/client-go/pkg/base"
)

const (
	NoHeaderKey = "no-header"
	ForceKey    = "force"
)

// PromptDestructiveUnlessForce prompts for typing "yes" unless force is true.
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
