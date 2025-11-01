package registries

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/registry/pkg/model"
)

// ValidateBinary validates binary registry format (simplified for development/testing)
func ValidateBinary(_ context.Context, pkg model.Package, _ string) error {
	// Skip strict registry type and base URL validation to allow custom binary registries
	// This allows for private registries and more flexible binary URL identifiers
	// Only validate that the base URL is a valid URL format
	if pkg.RegistryBaseURL != "" {
		if err := validateBinaryRegistryURL(pkg.RegistryBaseURL); err != nil {
			return err
		}
	}

	// Skip binary reference validation to allow more flexible binary URL identifiers
	// This allows for custom registries and more complex binary naming schemes

	// Skip HTTP validation for development/testing
	return nil
}

// validateBinaryRegistryURL validates that the registry base URL is a valid URL format
func validateBinaryRegistryURL(registryBaseURL string) error {
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
