package registries

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/modelcontextprotocol/registry/pkg/model"
)

func ValidateBinary(ctx context.Context, pkg model.Package, _ string) error {

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

	host := strings.ToLower(parsedURL.Host)
	allowedHosts := []string{
		githubHost,
		githubWWWHost,
		gitlabHost,
		gitlabWWWHost,
	}

	isAllowed := false
	for _, allowed := range allowedHosts {
		if host == allowed {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		return fmt.Errorf("binary packages must be hosted on allowlisted providers (GitHub or GitLab). Host '%s' is not allowed", host)
	}


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
		return "", fmt.Errorf("invalid host for binary package: %s, expected github or gitlab", host)
	}
}
