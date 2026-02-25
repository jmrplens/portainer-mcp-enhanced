package mcp

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (s *PortainerMCPServer) AddRoleFeatures() {
	s.addToolIfExists(ToolListRoles, s.HandleListRoles())
}

func (s *PortainerMCPServer) HandleListRoles() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		roles, err := s.cli.GetRoles()
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to list roles", err), nil
		}

		data, err := json.Marshal(roles)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to marshal roles", err), nil
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}
