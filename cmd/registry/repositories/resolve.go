package repositories

import (
	"fmt"
	"sort"
	"strings"

	"github.com/thalassa-cloud/client-go/containerregistry"
)

// ResolveNamespaceByName returns the namespace with the given name.
// An identity must be passed with --id.
func ResolveNamespaceByName(namespaces []containerregistry.ContainerRegistryNamespace, name string) (*containerregistry.ContainerRegistryNamespace, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("namespace name is required")
	}

	var exact []containerregistry.ContainerRegistryNamespace
	var folded []containerregistry.ContainerRegistryNamespace
	var identityMatch *containerregistry.ContainerRegistryNamespace
	for _, ns := range namespaces {
		switch {
		case ns.Namespace == name:
			exact = append(exact, ns)
		case strings.EqualFold(ns.Namespace, name):
			folded = append(folded, ns)
		}
		if ns.Identity == name {
			match := ns
			identityMatch = &match
		}
	}

	matches := exact
	if len(matches) == 0 {
		matches = folded
	}
	switch len(matches) {
	case 1:
		match := matches[0]
		if match.Identity == "" {
			return nil, fmt.Errorf("namespace %q has no identity", name)
		}
		return &match, nil
	case 0:
		if identityMatch != nil {
			return nil, fmt.Errorf("namespace %q is an identity; use --id to select a namespace by identity", name)
		}
		return nil, fmt.Errorf("namespace not found: %s", name)
	default:
		ids := make([]string, len(matches))
		for i, ns := range matches {
			ids[i] = ns.Identity
		}
		sort.Strings(ids)
		return nil, fmt.Errorf("namespace %q matches more than one registry namespace (%s); use --id", name, strings.Join(ids, ", "))
	}
}
