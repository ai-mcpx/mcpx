package registries

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/modelcontextprotocol/registry/pkg/model"
)


func ValidateWheel(_ context.Context, pkg model.Package, _ string) error {
	err := validateWheelURL(pkg.Identifier)
	if err != nil {
		return err
	}

	// Infer registry base URL if not provided
	if pkg.RegistryBaseURL == "" {
		baseURL, err := inferWheelRegistryBaseURL(pkg.Identifier)
		if err != nil {
			return fmt.Errorf("failed to infer registry base URL: %w", err)
		}
		pkg.RegistryBaseURL = baseURL
	}

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

	// Skip host allowlist validation to allow more flexible hosting options
	// This allows wheel packages to be hosted on any provider

	// Validate that the filename ends with .whl
	if !strings.HasSuffix(strings.ToLower(parsedURL.Path), WheelExtension) {
		return fmt.Errorf("wheel package URL must end with .whl extension: %s", fullURL)
	}

	// Basic wheel filename validation (name-version-py3-none-any.whl pattern)
	wheelFilename := strings.ToLower(parsedURL.Path)
	if strings.Contains(wheelFilename, "/") {
		parts := strings.Split(wheelFilename, "/")
		wheelFilename = parts[len(parts)-1]
	}

	// Remove .whl extension for validation
	wheelFilename = strings.TrimSuffix(wheelFilename, WheelExtension)

	// Basic pattern validation for wheel filename
	wheelPattern := regexp.MustCompile(`^[a-zA-Z0-9_-]+-[0-9]+(\.[0-9]+)*[a-zA-Z0-9_-]*$`)
	if !wheelPattern.MatchString(wheelFilename) {
		return fmt.Errorf("package URL must point to a valid wheel filename")
	}

	return nil
}

func inferWheelRegistryBaseURL(identifier string) (string, error) {
	parsedURL, err := url.Parse(identifier)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	host := strings.ToLower(parsedURL.Host)

	switch host {
	case GitHubHost, GitHubWWWHost:
		return model.RegistryURLGitHub, nil
	case GitLabHost, GitLabWWWHost:
		return model.RegistryURLGitLab, nil
	default:
		// For non-GitHub/GitLab hosts, use the scheme and host as the base URL
		return fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host), nil
	}
}