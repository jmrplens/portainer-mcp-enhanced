package mcp

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/portainer/portainer-mcp/pkg/toolgen"
)

func (s *PortainerMCPServer) AddAuthFeatures() {
	s.addToolIfExists(ToolAuthenticate, s.HandleAuthenticateUser())

	if !s.readOnly {
		s.addToolIfExists(ToolLogout, s.HandleLogout())
	}
}

func (s *PortainerMCPServer) HandleAuthenticateUser() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		parser := toolgen.NewParameterParser(request)

		username, err := parser.GetString("username", true)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("invalid username parameter", err), nil
		}

		password, err := parser.GetString("password", true)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("invalid password parameter", err), nil
		}

		authResponse, err := s.cli.AuthenticateUser(username, password)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to authenticate user", err), nil
		}

		data, err := json.Marshal(authResponse)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to marshal authentication response", err), nil
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}

func (s *PortainerMCPServer) HandleLogout() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		err := s.cli.Logout()
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to logout", err), nil
		}

		return mcp.NewToolResultText("Logged out successfully"), nil
	}
}
