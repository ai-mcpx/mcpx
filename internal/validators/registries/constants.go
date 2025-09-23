package registries

// Registry validation constants
const (
	// Supported registry types
	RegistryTypeNPM    = "npm"
	RegistryTypePyPI   = "pypi"
	RegistryTypeNuGet  = "nuget"
	RegistryTypeDocker = "docker"
	RegistryTypeBinary = "binary"
	RegistryTypeWheel  = "wheel"

	// Supported repository sources
	RepositorySourceGitHub = "github"
	RepositorySourceGitLab = "gitlab"
	RepositorySourceGerrit = "gerrit"

	// Common URL schemes and hosts
	HTTPSScheme    = "https"
	HTTPScheme     = "http"
	GitHubHost     = "github.com"
	GitHubWWWHost  = "www.github.com"
	GitLabHost     = "gitlab.com"
	GitLabWWWHost  = "www.gitlab.com"
	WheelExtension = ".whl"
)