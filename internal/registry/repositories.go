package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/thalassa-cloud/client-go/containerregistry"
	"github.com/thalassa-cloud/client-go/filters"
)

// Repository is a container registry repository.
// Namespace is decoded separately because the API returns either a string or an object.
type Repository struct {
	Identity       string       `json:"identity"`
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
	ObjectVersion  int64        `json:"objectVersion"`
	Namespace      NamespaceRef `json:"namespace,omitempty"`
	Image          string       `json:"image"`
	FullName       string       `json:"full_name"`
	Description    string       `json:"description"`
	LastPulledAt   *time.Time   `json:"last_pulled_at"`
	LastPushedAt   *time.Time   `json:"last_pushed_at"`
	TotalSizeBytes int64        `json:"total_size_bytes"`
	TagCount       int64        `json:"tag_count"`
	ArtifactCount  int64        `json:"artifact_count"`
	Tags           []ImageTag   `json:"tags,omitempty"`
}

// ImageTag is a tag on a repository. Nested repository objects are ignored.
type ImageTag struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Tag       string    `json:"tag"`
	Sha256    string    `json:"sha256"`
	SizeMb    float64   `json:"size_mb"`
}

// NamespaceRef accepts a namespace identity string or a namespace object.
type NamespaceRef struct {
	Identity  string
	Namespace string
}

// UnmarshalJSON accepts the namespace forms returned by the registry API.
func (n *NamespaceRef) UnmarshalJSON(data []byte) error {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" || trimmed == "null" {
		return nil
	}
	if trimmed[0] == '"' {
		var value string
		if err := json.Unmarshal(data, &value); err != nil {
			return fmt.Errorf("invalid namespace string: %w", err)
		}
		n.Identity = value
		n.Namespace = value
		return nil
	}

	var object struct {
		Identity  string `json:"identity"`
		Namespace string `json:"namespace"`
	}
	if err := json.Unmarshal(data, &object); err != nil {
		return fmt.Errorf("invalid namespace object: %w", err)
	}
	n.Identity = object.Identity
	n.Namespace = object.Namespace
	return nil
}

// ListRepositories lists repositories for a namespace identity.
func ListRepositories(ctx context.Context, client *containerregistry.Client, namespaceIdentity string, requestFilters filters.Filters) ([]Repository, error) {
	if strings.TrimSpace(namespaceIdentity) == "" {
		return nil, fmt.Errorf("namespace identity is required")
	}

	path, err := repositoriesPath(namespaceIdentity, "", requestFilters)
	if err != nil {
		return nil, err
	}
	body, err := getRegistry(ctx, client, path)
	if err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(string(body))) == 0 {
		return nil, nil
	}

	var repositories []Repository
	if err := json.Unmarshal(body, &repositories); err != nil {
		return nil, fmt.Errorf("failed to read repositories: %w", err)
	}
	return repositories, nil
}

// GetRepository returns one repository by identity.
func GetRepository(ctx context.Context, client *containerregistry.Client, namespaceIdentity, repositoryIdentity string) (*Repository, error) {
	if strings.TrimSpace(namespaceIdentity) == "" {
		return nil, fmt.Errorf("namespace identity is required")
	}
	if strings.TrimSpace(repositoryIdentity) == "" {
		return nil, fmt.Errorf("repository identity is required")
	}

	path, err := repositoriesPath(namespaceIdentity, repositoryIdentity, nil)
	if err != nil {
		return nil, err
	}
	body, err := getRegistry(ctx, client, path)
	if err != nil {
		return nil, err
	}

	var repository Repository
	if err := json.Unmarshal(body, &repository); err != nil {
		return nil, fmt.Errorf("failed to read repository: %w", err)
	}
	return &repository, nil
}

func getRegistry(ctx context.Context, client *containerregistry.Client, path string) ([]byte, error) {
	if client == nil {
		return nil, fmt.Errorf("registry client is required")
	}
	resp, err := client.RawRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	if err := client.Check(resp); err != nil {
		return nil, err
	}
	return resp.Body(), nil
}

func repositoriesPath(namespaceIdentity, repositoryIdentity string, requestFilters filters.Filters) (string, error) {
	path := containerregistry.ContainerRegistryEndpoint + "/" + url.PathEscape(namespaceIdentity) + "/repositories"
	if repositoryIdentity != "" {
		path += "/" + url.PathEscape(repositoryIdentity)
	}

	query := url.Values{}
	for _, filter := range requestFilters {
		if filter == nil {
			continue
		}
		for key, value := range filter.ToParams() {
			query.Set(key, value)
		}
	}
	if len(query) == 0 {
		return path, nil
	}
	return path + "?" + query.Encode(), nil
}
