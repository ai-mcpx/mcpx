package registries

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/modelcontextprotocol/registry/pkg/model"
)

// Wheel validation for development/testing

// ValidateWheel validates wheel package format (simplified for development/testing)
func ValidateWheel(_ context.Context, pkg model.Package, _ string) error {
	// Set default registry base URL if empty
	if pkg.RegistryBaseURL == "" {
		pkg.RegistryBaseURL = model.RegistryURLPyPI
	}

	// Skip strict registry type and base URL validation to allow custom wheel registries
	// This allows for private registries and more flexible wheel package identifiers
	// Only validate that the base URL is a valid URL format
	if pkg.RegistryBaseURL != "" {
		if err := validateWheelRegistryURL(pkg.RegistryBaseURL); err != nil {
			return err
		}
	}

	// Validate wheel package URL
	if err := validateWheelURL(pkg.Identifier); err != nil {
		return err
	}

	// Skip HTTP validation for development/testing
	return nil
}

func validateWheelURL(fullURL string) error {
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		return fmt.Errorf("invalid wheel package URL: %w", err)
	}

	// Check if it's a valid HTTP/HTTPS URL
	if parsedURL.Scheme != HTTPScheme && parsedURL.Scheme != HTTPSScheme {
		return fmt.Errorf("package URL must use http or https scheme")
	}

	// Check if the URL has a valid host
	if parsedURL.Host == "" {
		return fmt.Errorf("package URL must have a valid host")
	}

	// Check if the URL ends with .whl extension
	if !strings.HasSuffix(strings.ToLower(parsedURL.Path), WheelExtension) {
		return fmt.Errorf("wheel package URL must end with %s extension", WheelExtension)
	}

	// Remove allowlist restriction - allow wheel packages from any host
	_ = strings.ToLower(parsedURL.Host) // Keep for potential future use

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
