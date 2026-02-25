package mcp

import (
	"context"
	"encoding/json"

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

		data, err := json.Marshal(motd)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to marshal MOTD", err), nil
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}
