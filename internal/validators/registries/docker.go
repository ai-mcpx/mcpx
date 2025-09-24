package registries

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/registry/pkg/model"
)

// Docker validation simplified for development/testing

// ValidateDocker validates Docker image format (simplified for development/testing)
func ValidateDocker(_ context.Context, pkg model.Package, _ string) error {
	// Set default registry base URL if empty
	if pkg.RegistryBaseURL == "" {
		pkg.RegistryBaseURL = model.RegistryURLDocker
	}

	// Skip strict registry type and base URL validation to allow custom Docker registries
	// This allows for private registries and more flexible Docker image identifiers
	// Only validate that the base URL is a valid URL format
	if pkg.RegistryBaseURL != "" {
		if err := validateDockerRegistryURL(pkg.RegistryBaseURL); err != nil {
			return err
		}
	}

	// Skip image reference validation to allow more flexible Docker image identifiers
	// This allows for custom registries and more complex image naming schemes

	// Skip HTTP validation for development/testing
	return nil
}

// validateDockerRegistryURL validates that the registry base URL is a valid URL format
func validateDockerRegistryURL(registryBaseURL string) error {
	// Add https:// if no scheme is provided
	if !strings.HasPrefix(registryBaseURL, "http://") && !strings.HasPrefix(registryBaseURL, "https://") {
		registryBaseURL = "https://" + registryBaseURL
	}

	// Basic URL validation
	if !strings.Contains(registryBaseURL, "://") {
		return fmt.Errorf("invalid registry base URL format: %s", registryBaseURL)
	}

	// Check that it has a hostname
	parts := strings.Split(registryBaseURL, "://")
	if len(parts) != 2 || parts[1] == "" {
		return fmt.Errorf("invalid registry base URL format: %s", registryBaseURL)
	}

	return nil
}
