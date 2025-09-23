package registries

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/modelcontextprotocol/registry/pkg/model"
)


func ValidateBinary(_ context.Context, pkg model.Package, _ string) error {
	err := validateBinaryURL(pkg.Identifier)
	if err != nil {
		return err
	}

	// Infer registry base URL if not provided
	if pkg.RegistryBaseURL == "" {
		baseURL, err := inferBinaryRegistryBaseURL(pkg.Identifier)
		if err != nil {
			return fmt.Errorf("failed to infer registry base URL: %w", err)
		}
		pkg.RegistryBaseURL = baseURL
	}

	return nil
}

func validateBinaryURL(fullURL string) error {
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		return fmt.Errorf("invalid binary package URL: %w", err)
	}

	// Check if it's a valid HTTP/HTTPS URL
	if parsedURL.Scheme != HTTPScheme && parsedURL.Scheme != HTTPSScheme {
		return fmt.Errorf("package URL must use http or https scheme")
	}

	// Check if the URL has a valid host
	if parsedURL.Host == "" {
		return fmt.Errorf("package URL must have a valid host")
	}

	// Remove allowlist restriction - allow binary packages from any host
	_ = strings.ToLower(parsedURL.Host) // Keep for potential future use

	return nil
}

func inferBinaryRegistryBaseURL(identifier string) (string, error) {
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
		// For other hosts, return the base URL with scheme and host
		return fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host), nil
	}
}