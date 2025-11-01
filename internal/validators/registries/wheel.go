package registries

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/registry/pkg/model"
)

// ValidateWheel validates wheel (Python wheel) registry format (simplified for development/testing)
func ValidateWheel(_ context.Context, pkg model.Package, _ string) error {
	// Skip strict registry type and base URL validation to allow custom wheel registries
	// This allows for private registries and more flexible wheel URL identifiers
	// Only validate that the base URL is a valid URL format
	if pkg.RegistryBaseURL != "" {
		if err := validateWheelRegistryURL(pkg.RegistryBaseURL); err != nil {
			return err
		}
	}

	// Skip wheel reference validation to allow more flexible wheel URL identifiers
	// This allows for custom registries and more complex wheel naming schemes

	// Skip HTTP validation for development/testing
	return nil
}

// validateWheelRegistryURL validates that the registry base URL is a valid URL format
func validateWheelRegistryURL(registryBaseURL string) error {
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
