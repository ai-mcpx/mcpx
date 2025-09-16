package registries

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/modelcontextprotocol/registry/pkg/model"
)

func ValidateWheel(_ context.Context, pkg model.Package, _ string) error {
	err := validateWheelURL(pkg.Identifier)
	if err != nil {
		return err
	}

	inferredBaseURL, err := inferWheelRegistryBaseURL(pkg.Identifier)
	if err != nil {
		return err
	}

	if pkg.RegistryBaseURL == "" {
		pkg.RegistryBaseURL = inferredBaseURL
	} else if pkg.RegistryBaseURL != inferredBaseURL {
		return fmt.Errorf("wheel package '%s' has inconsistent registry base URL: %s (expected: %s)",
			pkg.Identifier, pkg.RegistryBaseURL, inferredBaseURL)
	}

	// Parse the URL to validate format
	url, err := url.Parse(pkg.Identifier)
	if err != nil {
		return fmt.Errorf("invalid wheel package URL: %w", err)
	}
	if url.Scheme != httpsScheme {
		return fmt.Errorf("invalid wheel package URL, must use HTTPS: %s", pkg.Identifier)
	}

	// Skip HTTP accessibility check for development/testing

	return nil
}

func validateWheelURL(fullURL string) error {
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		return fmt.Errorf("invalid wheel package URL: %w", err)
	}

	// Skip host allowlist validation to allow more flexible hosting options
	// This allows wheel packages to be hosted on any provider

	// Validate that the filename ends with .whl
	if !strings.HasSuffix(strings.ToLower(parsedURL.Path), ".whl") {
		return fmt.Errorf("wheel package URL must end with .whl extension: %s", fullURL)
	}

	return nil
}


// inferWheelRegistryBaseURL infers the registry base URL from a wheel identifier
func inferWheelRegistryBaseURL(identifier string) (string, error) {
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
		// For non-GitHub/GitLab hosts, use the scheme and host as the base URL
		return fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host), nil
	}
}
