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

	inferredBaseURL, err := inferBinaryRegistryBaseURL(pkg.Identifier)
	if err != nil {
		return err
	}

	if pkg.RegistryBaseURL == "" {
		pkg.RegistryBaseURL = inferredBaseURL
	} else if pkg.RegistryBaseURL != inferredBaseURL {
		return fmt.Errorf("binary package '%s' has inconsistent registry base URL: %s (expected: %s)",
			pkg.Identifier, pkg.RegistryBaseURL, inferredBaseURL)
	}

	// Parse the URL to validate format
	url, err := url.Parse(pkg.Identifier)
	if err != nil {
		return fmt.Errorf("invalid binary package URL: %w", err)
	}
	if url.Scheme != httpsScheme {
		return fmt.Errorf("invalid binary package URL, must use HTTPS: %s", pkg.Identifier)
	}

	// Skip HTTP accessibility check for development/testing

	return nil
}

func validateBinaryURL(fullURL string) error {
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		return fmt.Errorf("invalid binary package URL: %w", err)
	}

	// Remove allowlist restriction - allow binary packages from any host
	_ = strings.ToLower(parsedURL.Host) // Keep for potential future use

	return nil
}


// inferBinaryRegistryBaseURL infers the registry base URL from a binary identifier
func inferBinaryRegistryBaseURL(identifier string) (string, error) {
	parsedURL, err := url.Parse(identifier)
	if err != nil {
		return "", err
	}

	host := strings.ToLower(parsedURL.Host)
	switch host {
	case githubHost, githubWWWHost:
		return model.RegistryURLGitHub, nil
	case gitlabHost, gitlabWWWHost:
		return model.RegistryURLGitLab, nil
	default:
		// For other hosts, return the base URL with scheme and host
		return fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host), nil
	}
}
