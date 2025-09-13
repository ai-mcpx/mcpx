# MCP Registry

The MCP registry provides MCP clients with a list of MCP servers, like an app store for MCP servers.

[**📤 Publish my MCP server**](docs/guides/publishing/publish-server.md) | [**⚡️ Live API docs**](https://registry.modelcontextprotocol.io/docs) | [**👀 Ecosystem vision**](docs/explanations/ecosystem-vision.md) | 📖 **[Full documentation](./docs)**

## Development Status

**2025-09-08 update**: The registry has launched in preview 🎉 ([announcement blog post](https://blog.modelcontextprotocol.io/posts/2025-09-08-mcp-registry-preview/)). While the system is now more stable, this is still a preview release and breaking changes or data resets may occur. A general availability (GA) release will follow later. We'd love your feedback in [GitHub discussions](https://github.com/modelcontextprotocol/registry/discussions/new?category=ideas) or in the [#registry-dev Discord](https://discord.com/channels/1358869848138059966/1369487942862504016) ([joining details here](https://modelcontextprotocol.io/community/communication)).

Current key maintainers:
- **Adam Jones** (Anthropic) [@domdomegg](https://github.com/domdomegg)
- **Tadas Antanavicius** (PulseMCP) [@tadasant](https://github.com/tadasant)
- **Toby Padilla** (GitHub) [@toby](https://github.com/toby)

## Contributing

We use multiple channels for collaboration - see [modelcontextprotocol.io/community/communication](https://modelcontextprotocol.io/community/communication).

Often (but not always) ideas flow through this pipeline:

- **[Discord](https://modelcontextprotocol.io/community/communication)** - Real-time community discussions
- **[Discussions](https://github.com/modelcontextprotocol/registry/discussions)** - Propose and discuss product/technical requirements
- **[Issues](https://github.com/modelcontextprotocol/registry/issues)** - Track well-scoped technical work
- **[Pull Requests](https://github.com/modelcontextprotocol/registry/pulls)** - Contribute work towards issues

### Quick start:

#### Pre-requisites

- **Docker**
- **Go 1.24.x**
- **golangci-lint v2.4.0**

#### Running the server

```bash
# Start full development environment
make dev-compose
```

This starts the registry at [`localhost:8080`](http://localhost:8080) with PostgreSQL and seed data. It can be configured with environment variables in [docker-compose.yml](./docker-compose.yml) - see [.env.example](./.env.example) for a reference.

<details>
<summary>Alternative: Local setup without Docker</summary>

**Prerequisites:**
- PostgreSQL running locally
- Go 1.24.x installed

```bash
# Build and run locally
make build
make dev-local
```

The service runs on [`localhost:8080`](http://localhost:8080) by default. This can be configured with environment variables in `.env` - see [.env.example](./.env.example) for a reference.

</details>

<details>
<summary>Alternative: Running a pre-built Docker image</summary>

Pre-built Docker images are automatically published to GitHub Container Registry:

```bash
# Run latest stable release
docker run -p 8080:8080 ghcr.io/modelcontextprotocol/registry:latest

# Run latest from main branch (continuous deployment)
docker run -p 8080:8080 ghcr.io/modelcontextprotocol/registry:main

# Run specific release version
docker run -p 8080:8080 ghcr.io/modelcontextprotocol/registry:v1.0.0

# Run development build from main branch
docker run -p 8080:8080 ghcr.io/modelcontextprotocol/registry:main-20250906-abc123d
```

**Available tags:**
- **Releases**: `latest`, `v1.0.0`, `v1.1.0`, etc.
- **Continuous**: `main` (latest main branch build)
- **Development**: `main-<date>-<sha>` (specific commit builds)

</details>

#### Publishing a server

To publish a server, we've built a simple CLI. You can use it with:

```bash
# Build the latest CLI
make publisher

# Use it!
./bin/mcp-publisher --help
```

**Example server configurations** are available in the `mcpx-cli/` directory:
- `example-server-binary.json` - Binary package distribution
- `example-server-wheel.json` - Python wheel package distribution
- `example-server-docker.json` - Docker image distribution

See [the publisher guide](./docs/guides/publishing/publish-server.md) for more details.

#### Other commands

```bash
# Run lint, unit tests and integration tests
make check
```

There are also a few more helpful commands for development. Run `make help` to learn more, or look in [Makefile](./Makefile).

<!--
For Claude and other AI tools: Always prefer make targets over custom commands where possible.
-->

## Architecture

### Project Structure

```
├── cmd/                     # Application entry points
│   └── publisher/           # Server publishing tool
├── data/                    # Seed data
├── deploy/                  # Deployment configuration (Pulumi)
├── docs/                    # Documentation
├── internal/                # Private application code
│   ├── api/                 # HTTP handlers and routing
│   ├── auth/                # Authentication (GitHub OAuth, JWT, namespace blocking)
│   ├── config/              # Configuration management
│   ├── database/            # Data persistence (PostgreSQL, in-memory)
│   ├── service/             # Business logic
│   ├── telemetry/           # Metrics and monitoring
│   └── validators/          # Input validation
├── pkg/                     # Public packages
│   ├── api/                 # API types and structures
│   │   └── v0/              # Version 0 API types
│   └── model/               # Data models for server.json
├── scripts/                 # Development and testing scripts
├── tests/                   # Integration tests
└── tools/                   # CLI tools and utilities
    └── validate-*.sh        # Schema validation tools
```

### Authentication

Publishing supports multiple authentication methods:
- **GitHub OAuth** - For publishing by logging into GitHub
- **GitHub OIDC** - For publishing from GitHub Actions
- **DNS verification** - For proving ownership of a domain and its subdomains
- **HTTP verification** - For proving ownership of a domain
- **Anonymous** - For publishing to anonymous namespaces (no authentication required)

The registry validates namespace ownership when publishing. E.g. to publish...:
- `io.github.domdomegg/my-cool-mcp` you must login to GitHub as `domdomegg`, or be in a GitHub Action on domdomegg's repos
- `me.adamjones/my-cool-mcp` you must prove ownership of `adamjones.me` via DNS or HTTP challenge
- `io.modelcontextprotocol.anonymous/my-cool-mcp` you can publish without authentication (anonymous namespace)

**Note**: Authentication is automatically determined based on the server namespace. GitHub namespaces (`io.github.*`) require GitHub authentication, while other namespaces can use alternative authentication methods or anonymous publishing.

### Repository Sources

The registry supports multiple repository sources for hosting your MCP server code:

- **GitHub** - Standard GitHub repositories (`https://github.com/user/repo`)
- **GitLab** - Standard GitLab repositories (`https://gitlab.com/user/repo`)
- **Gerrit** - Gerrit code review repositories (`http://host:port/project/path`)

When specifying a repository in your server configuration, use the appropriate `source` field value (`github`, `gitlab`, or `gerrit`) along with the repository URL.

### Package Registry Types

The registry supports multiple package registry types for distributing your MCP server:

- **npm** - Node.js packages via npm registry
- **pypi** - Python packages via PyPI registry
- **oci** - Container images via OCI-compatible registries (Docker Hub, etc.)
- **docker** - Docker images via Docker Hub (simplified validation)
- **nuget** - .NET packages via NuGet registry
- **mcpb** - MCP binary packages via direct download
- **binary** - Generic binary packages via direct download (GitHub/GitLab releases)
- **wheel** - Python wheel packages via PyPI or direct download (GitHub/GitLab releases)

**Validation Notes:**
- **npm, pypi, nuget, mcpb**: Full validation including package existence and ownership verification
- **oci, docker, binary, wheel**: Simplified validation for development/testing - only format validation, no HTTP checks

**Development Features:**
- **Localhost URLs**: Remote URLs can use localhost for development and testing
- **No SHA Requirements**: Binary and wheel packages don't require file_sha256 for integrity verification
- **No HTTP Validation**: OCI, docker, binary, and wheel packages skip HTTP accessibility checks

Each package type has specific requirements for versioning, file hashing, and distribution methods. See the [API documentation](#api-documentation) for detailed schema requirements.

## API Documentation

The API is documented using Swagger/OpenAPI. This page provides a complete reference of all endpoints with request/response schemas and examples, and allows you to test the API directly from your browser.

### Key Endpoints

- `GET /v0/servers` - List all registered servers with pagination
- `GET /v0/servers/{id}` - Get details of a specific server by ID
- `PUT /v0/servers/{id}` - Update a specific server by ID
- `DELETE /v0/servers/{id}` - Delete a specific server by ID
- `POST /v0/publish` - Publish a new server to the registry

The `PUT /v0/servers/{id}` endpoint allows updating server details including version information. When updating a version, it must be greater than the existing version to maintain version ordering.

**Note**: The `DELETE /v0/servers/{id}` endpoint permanently removes a server from the registry. This action cannot be undone.

## Configuration

The service can be configured using environment variables. See [.env.example](./.env.example) for details.

## More documentation

See the [documentation](./docs) for more details if your question has not been answered here!
