package configuration

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/thalassa-cloud/client-go/containerregistry"
)

func loadRetentionPolicyFromFile(path string) (*containerregistry.RetentionPolicy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read retention policy file: %w", err)
	}
	var policy containerregistry.RetentionPolicy
	if err := json.Unmarshal(data, &policy); err != nil {
		return nil, fmt.Errorf("failed to parse retention policy JSON: %w", err)
	}
	return &policy, nil
}

func buildRetentionPolicy(enabled, deleteUntagged bool, days, count int, retentionPolicyFile string) (*containerregistry.RetentionPolicy, error) {
	if retentionPolicyFile != "" {
		return loadRetentionPolicyFromFile(retentionPolicyFile)
	}
	if !enabled && !deleteUntagged && days == 0 && count == 0 {
		return nil, nil
	}

	policy := &containerregistry.RetentionPolicy{
		Enabled:              enabled,
		DeleteUntaggedImages: deleteUntagged,
	}
	if days > 0 || count > 0 {
		rule := containerregistry.RetentionPolicyRule{
			Scope: containerregistry.RetentionPolicyScopeTags,
		}
		if days > 0 {
			rule.Days = &days
		}
		if count > 0 {
			rule.Count = &count
		}
		policy.Rules = []containerregistry.RetentionPolicyRule{rule}
	}
	return policy, nil
}
