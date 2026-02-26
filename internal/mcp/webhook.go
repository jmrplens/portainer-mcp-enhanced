package mcp

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/portainer/portainer-mcp/pkg/toolgen"
)

func (s *PortainerMCPServer) AddWebhookFeatures() {
	s.addToolIfExists(ToolListWebhooks, s.HandleListWebhooks())

	if !s.readOnly {
		s.addToolIfExists(ToolCreateWebhook, s.HandleCreateWebhook())
		s.addToolIfExists(ToolDeleteWebhook, s.HandleDeleteWebhook())
	}
}

func (s *PortainerMCPServer) HandleListWebhooks() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		webhooks, err := s.cli.GetWebhooks()
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to get webhooks", err), nil
		}

		return jsonResult(webhooks, "failed to marshal webhooks")
	}
}

func (s *PortainerMCPServer) HandleCreateWebhook() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		parser := toolgen.NewParameterParser(request)

		resourceId, err := parser.GetString("resourceId", true)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("invalid resourceId parameter", err), nil
		}

		endpointId, err := parser.GetInt("endpointId", true)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("invalid endpointId parameter", err), nil
		}

		webhookType, err := parser.GetInt("webhookType", true)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("invalid webhookType parameter", err), nil
		}

		id, err := s.cli.CreateWebhook(resourceId, endpointId, webhookType)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to create webhook", err), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Webhook created successfully with ID: %d", id)), nil
	}
}

func (s *PortainerMCPServer) HandleDeleteWebhook() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		parser := toolgen.NewParameterParser(request)

		id, err := parser.GetInt("id", true)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("invalid id parameter", err), nil
		}

		err = s.cli.DeleteWebhook(id)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to delete webhook", err), nil
		}

		return mcp.NewToolResultText("Webhook deleted successfully"), nil
	}
}
