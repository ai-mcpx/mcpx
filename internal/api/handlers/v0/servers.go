package v0

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/modelcontextprotocol/registry/internal/auth"
	"github.com/modelcontextprotocol/registry/internal/config"
	"github.com/modelcontextprotocol/registry/internal/database"
	"github.com/modelcontextprotocol/registry/internal/service"
	apiv0 "github.com/modelcontextprotocol/registry/pkg/api/v0"
	"github.com/modelcontextprotocol/registry/pkg/model"
)

// ListServersInput represents the input for listing servers
type ListServersInput struct {
	Cursor       string `query:"cursor" doc:"Pagination cursor (UUID)" format:"uuid" required:"false" example:"550e8400-e29b-41d4-a716-446655440000"`
	Limit        int    `query:"limit" doc:"Number of items per page" default:"30" minimum:"1" maximum:"100" example:"50"`
	UpdatedSince string `query:"updated_since" doc:"Filter servers updated since timestamp (RFC3339 datetime)" required:"false" example:"2025-08-07T13:15:04.280Z"`
	Search       string `query:"search" doc:"Search servers by name (substring match)" required:"false" example:"filesystem"`
	Version      string `query:"version" doc:"Filter by version ('latest' for latest version, or an exact version like '1.2.3')" required:"false" example:"latest"`
}

// ServerDetailInput represents the input for getting server details
type ServerDetailInput struct {
	ID string `path:"id" doc:"Server ID (UUID)" format:"uuid"`
}

// UpdateServerInput represents the input for updating a server
type UpdateServerInput struct {
	ID            string           `path:"id" doc:"Server ID (UUID)" format:"uuid"`
	Authorization string           `header:"Authorization" doc:"Registry JWT token" required:"false"`
	Body          apiv0.ServerJSON `body:""`
}

// DeleteServerInput represents the input for deleting a server
type DeleteServerInput struct {
	ID            string `path:"id" doc:"Server ID (UUID)" format:"uuid"`
	Authorization string `header:"Authorization" doc:"Registry JWT token" required:"false"`
}

// RegisterServersEndpoints registers all server-related endpoints
func RegisterServersEndpoints(api huma.API, registry service.RegistryService, cfg *config.Config) {
	registerListServersEndpoint(api, registry)
	registerGetServerEndpoint(api, registry)
	registerUpdateServerEndpoint(api, registry, cfg)
	registerDeleteServerEndpoint(api, registry, cfg)
}

// registerListServersEndpoint registers the list servers endpoint
func registerListServersEndpoint(api huma.API, registry service.RegistryService) {
	huma.Register(api, huma.Operation{
		OperationID: "list-servers",
		Method:      http.MethodGet,
		Path:        "/v0/servers",
		Summary:     "List MCP servers",
		Description: "Get a paginated list of MCP servers from the registry",
		Tags:        []string{"servers"},
	}, func(_ context.Context, input *ListServersInput) (*Response[apiv0.ServerListResponse], error) {
		// Validate cursor if provided
		if input.Cursor != "" {
			_, err := uuid.Parse(input.Cursor)
			if err != nil {
				return nil, huma.Error400BadRequest("Invalid cursor parameter")
			}
		}

		// Build filter from input parameters
		filter := &database.ServerFilter{}

		// Parse updated_since parameter
		if input.UpdatedSince != "" {
			// Parse RFC3339 format
			if updatedTime, err := time.Parse(time.RFC3339, input.UpdatedSince); err == nil {
				filter.UpdatedSince = &updatedTime
			} else {
				return nil, huma.Error400BadRequest("Invalid updated_since format: expected RFC3339 timestamp (e.g., 2025-08-07T13:15:04.280Z)")
			}
		}

		// Handle search parameter
		if input.Search != "" {
			filter.SubstringName = &input.Search
		}

		// Handle version parameter
		if input.Version != "" {
			if input.Version == "latest" {
				// Special case: filter for latest versions
				isLatest := true
				filter.IsLatest = &isLatest
			} else {
				// Future: exact version matching
				filter.Version = &input.Version
			}
		}

		// Get paginated results with filtering
		servers, nextCursor, err := registry.List(filter, input.Cursor, input.Limit)
		if err != nil {
			return nil, huma.Error500InternalServerError("Failed to get registry list", err)
		}

		return &Response[apiv0.ServerListResponse]{
			Body: apiv0.ServerListResponse{
				Servers: servers,
				Metadata: apiv0.Metadata{
					NextCursor: nextCursor,
					Count:      len(servers),
				},
			},
		}, nil
	})
}

// registerGetServerEndpoint registers the get server details endpoint
func registerGetServerEndpoint(api huma.API, registry service.RegistryService) {
	huma.Register(api, huma.Operation{
		OperationID: "get-server",
		Method:      http.MethodGet,
		Path:        "/v0/servers/{id}",
		Summary:     "Get MCP server details",
		Description: "Get detailed information about a specific MCP server",
		Tags:        []string{"servers"},
	}, func(_ context.Context, input *ServerDetailInput) (*Response[apiv0.ServerJSON], error) {
		// Get the server details from the registry service
		serverDetail, err := registry.GetByID(input.ID)
		if err != nil {
			if err.Error() == database.RecordNotFoundMsg {
				return nil, huma.Error404NotFound("Server not found")
			}
			return nil, huma.Error500InternalServerError("Failed to get server details", err)
		}

		return &Response[apiv0.ServerJSON]{
			Body: *serverDetail,
		}, nil
	})
}

// registerUpdateServerEndpoint registers the update server endpoint
func registerUpdateServerEndpoint(api huma.API, registry service.RegistryService, cfg *config.Config) {
	huma.Register(api, huma.Operation{
		OperationID: "update-server",
		Method:      http.MethodPut,
		Path:        "/v0/servers/{id}",
		Summary:     "Update MCP server",
		Description: "Update an existing MCP server in the registry",
		Tags:        []string{"servers"},
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}, func(ctx context.Context, input *UpdateServerInput) (*Response[apiv0.ServerJSON], error) {
		// Get the existing server to check permissions
		existingServer, err := registry.GetByID(input.ID)
		if err != nil {
			if err.Error() == database.RecordNotFoundMsg {
				return nil, huma.Error404NotFound("Server not found")
			}
			return nil, huma.Error500InternalServerError("Failed to get server details", err)
		}

		// Create JWT manager for token validation
		jwtManager := auth.NewJWTManager(cfg)

		// Determine auth method based on server namespace
		var authMethod auth.Method
		if strings.HasPrefix(existingServer.Name, "io.github.") {
			authMethod = auth.MethodGitHubAT
		} else {
			authMethod = auth.MethodNone
		}

		// Extract token if provided
		var token string
		if input.Authorization != "" {
			const bearerPrefix = "Bearer "
			authHeader := input.Authorization
			if len(authHeader) < len(bearerPrefix) || !strings.EqualFold(authHeader[:len(bearerPrefix)], bearerPrefix) {
				return nil, huma.Error401Unauthorized("Invalid Authorization header format. Expected 'Bearer <token>'")
			}
			token = authHeader[len(bearerPrefix):]
		}

		// Validate authentication only if required by auth method or if token is provided
		if authMethod != auth.MethodNone {
			// GitHub auth method requires authentication
			if token == "" {
				return nil, huma.Error401Unauthorized("Authentication is required for this server namespace")
			}

			claims, err := jwtManager.ValidateToken(ctx, token)
			if err != nil {
				return nil, huma.Error401Unauthorized("Invalid or expired Registry JWT token", err)
			}

			// Verify that the token's repository matches the server being updated
			if !jwtManager.HasPermission(existingServer.Name, auth.PermissionActionEdit, claims.Permissions) {
				return nil, huma.Error403Forbidden("You do not have permission to update this server")
			}
		} else if token != "" {
			// If a token is provided but auth method is None, validate it anyway
			claims, err := jwtManager.ValidateToken(ctx, token)
			if err != nil {
				return nil, huma.Error401Unauthorized("Invalid or expired Registry JWT token", err)
			}

			// Verify that the token's repository matches the server being updated
			if !jwtManager.HasPermission(existingServer.Name, auth.PermissionActionEdit, claims.Permissions) {
				return nil, huma.Error403Forbidden("You do not have permission to update this server")
			}
		}

		// Prevent renaming servers
		if existingServer.Name != input.Body.Name {
			return nil, huma.Error400BadRequest("Cannot rename server")
		}

		// Prevent undeleting servers - once deleted, they stay deleted
		if existingServer.Status == model.StatusDeleted && input.Body.Status != model.StatusDeleted {
			return nil, huma.Error400BadRequest("Cannot change status of deleted server. Deleted servers cannot be undeleted.")
		}

		// Update the server
		updatedServer, err := registry.Update(input.ID, input.Body)
		if err != nil {
			return nil, huma.Error400BadRequest("Failed to update server", err)
		}

		return &Response[apiv0.ServerJSON]{
			Body: *updatedServer,
		}, nil
	})
}

// registerDeleteServerEndpoint registers the delete server endpoint
func registerDeleteServerEndpoint(api huma.API, registry service.RegistryService, cfg *config.Config) {
	huma.Register(api, huma.Operation{
		OperationID: "delete-server",
		Method:      http.MethodDelete,
		Path:        "/v0/servers/{id}",
		Summary:     "Delete MCP server",
		Description: "Delete an MCP server from the registry",
		Tags:        []string{"servers"},
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}, func(ctx context.Context, input *DeleteServerInput) (*Response[struct{}], error) {
		// Get the existing server to check permissions
		existingServer, err := registry.GetByID(input.ID)
		if err != nil {
			if err.Error() == database.RecordNotFoundMsg {
				return nil, huma.Error404NotFound("Server not found")
			}
			return nil, huma.Error500InternalServerError("Failed to get server details", err)
		}

		// Create JWT manager for token validation
		jwtManager := auth.NewJWTManager(cfg)

		// Determine auth method based on server namespace
		var authMethod auth.Method
		if strings.HasPrefix(existingServer.Name, "io.github.") {
			authMethod = auth.MethodGitHubAT
		} else {
			authMethod = auth.MethodNone
		}

		// Extract token if provided
		var token string
		if input.Authorization != "" {
			const bearerPrefix = "Bearer "
			authHeader := input.Authorization
			if len(authHeader) < len(bearerPrefix) || !strings.EqualFold(authHeader[:len(bearerPrefix)], bearerPrefix) {
				return nil, huma.Error401Unauthorized("Invalid Authorization header format. Expected 'Bearer <token>'")
			}
			token = authHeader[len(bearerPrefix):]
		}

		// Validate authentication only if required by auth method or if token is provided
		if authMethod != auth.MethodNone {
			// GitHub auth method requires authentication
			if token == "" {
				return nil, huma.Error401Unauthorized("Authentication is required for this server namespace")
			}

			claims, err := jwtManager.ValidateToken(ctx, token)
			if err != nil {
				return nil, huma.Error401Unauthorized("Invalid or expired Registry JWT token", err)
			}

			// Verify that the token's repository matches the server being deleted
			if !jwtManager.HasPermission(existingServer.Name, auth.PermissionActionEdit, claims.Permissions) {
				return nil, huma.Error403Forbidden("You do not have permission to delete this server")
			}
		} else if token != "" {
			// If a token is provided but auth method is None, validate it anyway
			claims, err := jwtManager.ValidateToken(ctx, token)
			if err != nil {
				return nil, huma.Error401Unauthorized("Invalid or expired Registry JWT token", err)
			}

			// Verify that the token's repository matches the server being deleted
			if !jwtManager.HasPermission(existingServer.Name, auth.PermissionActionEdit, claims.Permissions) {
				return nil, huma.Error403Forbidden("You do not have permission to delete this server")
			}
		}

		// Delete the server
		err = registry.Delete(input.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("Failed to delete server", err)
		}

		return &Response[struct{}]{
			Body: struct{}{},
		}, nil
	})
}
