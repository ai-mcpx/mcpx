package v0

import (
	"context"
	"net/http"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/modelcontextprotocol/registry/internal/auth"
	"github.com/modelcontextprotocol/registry/internal/config"
	"github.com/modelcontextprotocol/registry/internal/service"
	"github.com/modelcontextprotocol/registry/internal/validators"
	apiv0 "github.com/modelcontextprotocol/registry/pkg/api/v0"
)

// PublishServerInput represents the input for publishing a server
type PublishServerInput struct {
	Authorization string           `header:"Authorization" doc:"Registry JWT token (obtained from /v0/auth/token/github)" required:"false"`
	Body          apiv0.ServerJSON `body:""`
}

// RegisterPublishEndpoint registers the publish endpoint
func RegisterPublishEndpoint(api huma.API, registry service.RegistryService, cfg *config.Config) {
	// Create JWT manager for token validation
	jwtManager := auth.NewJWTManager(cfg)

	huma.Register(api, huma.Operation{
		OperationID: "publish-server",
		Method:      http.MethodPost,
		Path:        "/v0/publish",
		Summary:     "Publish MCP server",
		Description: "Publish a new MCP server to the registry or update an existing one",
		Tags:        []string{"publish"},
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}, func(ctx context.Context, input *PublishServerInput) (*Response[apiv0.ServerJSON], error) {
		// Validate the publish request
		if err := validators.ValidatePublishRequest(ctx, input.Body, cfg); err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		// Get server details from request body
		serverDetail := input.Body

		// Determine auth method based on server namespace
		var authMethod auth.Method
		if strings.HasPrefix(serverDetail.Name, "io.github.") {
			authMethod = auth.MethodGitHubAT // or MethodGitHubOIDC - both require GitHub auth
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

			// Verify that the token's repository matches the server being published
			if !jwtManager.HasPermission(serverDetail.Name, auth.PermissionActionPublish, claims.Permissions) {
				return nil, huma.Error403Forbidden("You do not have permission to publish this server")
			}
		} else if token != "" {
			// If a token is provided but auth method is None, validate it anyway
			claims, err := jwtManager.ValidateToken(ctx, token)
			if err != nil {
				return nil, huma.Error401Unauthorized("Invalid or expired Registry JWT token", err)
			}

			// Verify that the token's repository matches the server being published
			if !jwtManager.HasPermission(serverDetail.Name, auth.PermissionActionPublish, claims.Permissions) {
				return nil, huma.Error403Forbidden("You do not have permission to publish this server")
			}
		}

		// Publish the server with extensions
		publishedServer, err := registry.Publish(input.Body)
		if err != nil {
			return nil, huma.Error400BadRequest("Failed to publish server", err)
		}

		// Return the published server in flattened format
		return &Response[apiv0.ServerJSON]{
			Body: *publishedServer,
		}, nil
	})
}

