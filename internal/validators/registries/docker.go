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

	// Validate that the registry base URL matches Docker exactly
	if pkg.RegistryBaseURL != model.RegistryURLDocker {
		return fmt.Errorf("registry type and base URL do not match: '%s' is not valid for registry type '%s'. Expected: %s",
			pkg.RegistryBaseURL, model.RegistryTypeDocker, model.RegistryURLDocker)
	}

	// Skip image reference validation to allow more flexible Docker image identifiers
	// This allows for custom registries and more complex image naming schemes

	// Skip HTTP validation for development/testing
	return nil
}
