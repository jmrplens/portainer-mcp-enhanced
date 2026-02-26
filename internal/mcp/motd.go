package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func (s *PortainerMCPServer) AddMotdFeatures() {
	s.addToolIfExists(ToolGetMOTD, s.HandleGetMOTD())
}

func (s *PortainerMCPServer) HandleGetMOTD() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		motd, err := s.cli.GetMOTD()
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to get MOTD", err), nil
		}

		return jsonResult(motd, "failed to marshal MOTD")
	}
}
