package registries

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/registry/pkg/model"
)

// Docker validation simplified for development/testing

// ValidateDocker validates Docker image format (simplified for development/testing)
func ValidateDocker(ctx context.Context, pkg model.Package, serverName string) error {
	// Set default registry base URL if empty
	if pkg.RegistryBaseURL == "" {
		pkg.RegistryBaseURL = model.RegistryURLDocker
	}

	// Validate that the registry base URL matches Docker exactly
	if pkg.RegistryBaseURL != model.RegistryURLDocker {
		return fmt.Errorf("registry type and base URL do not match: '%s' is not valid for registry type '%s'. Expected: %s",
			pkg.RegistryBaseURL, model.RegistryTypeDocker, model.RegistryURLDocker)
	}

	// Parse image reference (namespace/repo or repo) to validate format
	_, _, err := parseDockerImageReference(pkg.Identifier)
	if err != nil {
		return fmt.Errorf("invalid Docker image reference: %w", err)
	}

	// Skip HTTP validation for development/testing
	return nil
}

func parseDockerImageReference(identifier string) (string, string, error) {
	parts := strings.Split(identifier, "/")
	switch len(parts) {
	case 2:
		return parts[0], parts[1], nil
	case 1:
		return "library", parts[0], nil
	default:
		return "", "", fmt.Errorf("invalid image reference: %s", identifier)
	}
}
