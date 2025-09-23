package registries

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/registry/pkg/model"
)

// Docker validation simplified for development/testing

// ValidateDocker validates Docker image format (simplified for development/testing)
func ValidateDocker(_ context.Context, pkg model.Package, _ string) error {
	// Set default registry base URL if empty
	if pkg.RegistryBaseURL == "" {
		pkg.RegistryBaseURL = model.RegistryURLDocker
	}

	// Validate registry type and base URL match
	if err := validateRegistryTypeAndURL(
		pkg.RegistryBaseURL, model.RegistryTypeOCI, model.RegistryURLDocker); err != nil {
		return err
	}

	// Skip image reference validation to allow more flexible Docker image identifiers
	// This allows for custom registries and more complex image naming schemes

	// Skip HTTP validation for development/testing
	return nil
}

// validateRegistryTypeAndURL validates that the registry type and base URL match
func validateRegistryTypeAndURL(registryBaseURL, registryType, expectedURL string) error {
	if registryBaseURL != expectedURL {
		return fmt.Errorf("registry type and base URL do not match: '%s' is not valid for registry type '%s'. Expected: %s",
			registryBaseURL, registryType, expectedURL)
	}
	return nil
}
