package shared

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	clientiam "github.com/thalassa-cloud/client-go/iam"
)

func ParseAccessCredentialScopes(ss []string) ([]clientiam.AccessCredentialsScope, error) {
	if len(ss) == 0 {
		return nil, nil
	}
	out := make([]clientiam.AccessCredentialsScope, 0, len(ss))
	for _, s := range ss {
		s = strings.TrimSpace(s)
		switch clientiam.AccessCredentialsScope(s) {
		case clientiam.AccessCredentialsScopeAPIRead,
			clientiam.AccessCredentialsScopeAPIWrite,
			clientiam.AccessCredentialsScopeKubernetes,
			clientiam.AccessCredentialsScopeObjectStorage:
			out = append(out, clientiam.AccessCredentialsScope(s))
		default:
			return nil, fmt.Errorf("invalid scope %q (allowed: api:read, api:write, kubernetes, objectStorage)", s)
		}
	}
	return out, nil
}

func ParseConditionsJSON(s, path string) (map[string]interface{}, error) {
	raw := s
	if path != "" {
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read conditions file: %w", err)
		}
		raw = string(b)
	}
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return nil, fmt.Errorf("parse conditions JSON: %w", err)
	}
	return m, nil
}

func ParseOptionalRFC3339(s string) (*time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, fmt.Errorf("parse time as RFC3339: %w", err)
	}
	return &t, nil
}
