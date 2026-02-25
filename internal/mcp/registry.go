package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/portainer/portainer-mcp/pkg/toolgen"
)

func (s *PortainerMCPServer) AddRegistryFeatures() {
	s.addToolIfExists(ToolListRegistries, s.HandleListRegistries())
	s.addToolIfExists(ToolGetRegistry, s.HandleGetRegistry())

	if !s.readOnly {
		s.addToolIfExists(ToolCreateRegistry, s.HandleCreateRegistry())
		s.addToolIfExists(ToolUpdateRegistry, s.HandleUpdateRegistry())
		s.addToolIfExists(ToolDeleteRegistry, s.HandleDeleteRegistry())
	}
}

func (s *PortainerMCPServer) HandleListRegistries() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		registries, err := s.cli.GetRegistries()
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to list registries", err), nil
		}

		data, err := json.Marshal(registries)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to marshal registries", err), nil
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}

func (s *PortainerMCPServer) HandleGetRegistry() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		parser := toolgen.NewParameterParser(request)

		id, err := parser.GetInt("id", true)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("invalid id parameter", err), nil
		}

		registry, err := s.cli.GetRegistry(id)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to get registry", err), nil
		}

		data, err := json.Marshal(registry)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to marshal registry", err), nil
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}

func (s *PortainerMCPServer) HandleCreateRegistry() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		parser := toolgen.NewParameterParser(request)

		name, err := parser.GetString("name", true)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("invalid name parameter", err), nil
		}

		registryType, err := parser.GetInt("type", true)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("invalid type parameter", err), nil
		}

		url, err := parser.GetString("url", true)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("invalid url parameter", err), nil
		}

		authentication, err := parser.GetBoolean("authentication", true)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("invalid authentication parameter", err), nil
		}

		username, _ := parser.GetString("username", false)
		password, _ := parser.GetString("password", false)
		baseURL, _ := parser.GetString("baseURL", false)

		id, err := s.cli.CreateRegistry(name, registryType, url, authentication, username, password, baseURL)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to create registry", err), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Registry created successfully with ID: %d", id)), nil
	}
}

func (s *PortainerMCPServer) HandleUpdateRegistry() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		parser := toolgen.NewParameterParser(request)

		id, err := parser.GetInt("id", true)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("invalid id parameter", err), nil
		}

		args := request.GetArguments()

		var name *string
		if _, ok := args["name"]; ok {
			v, err := parser.GetString("name", false)
			if err != nil {
				return mcp.NewToolResultErrorFromErr("invalid name parameter", err), nil
			}
			name = &v
		}

		var url *string
		if _, ok := args["url"]; ok {
			v, err := parser.GetString("url", false)
			if err != nil {
				return mcp.NewToolResultErrorFromErr("invalid url parameter", err), nil
			}
			url = &v
		}

		var authentication *bool
		if _, ok := args["authentication"]; ok {
			v, err := parser.GetBoolean("authentication", false)
			if err != nil {
				return mcp.NewToolResultErrorFromErr("invalid authentication parameter", err), nil
			}
			authentication = &v
		}

		var username *string
		if _, ok := args["username"]; ok {
			v, err := parser.GetString("username", false)
			if err != nil {
				return mcp.NewToolResultErrorFromErr("invalid username parameter", err), nil
			}
			username = &v
		}

		var password *string
		if _, ok := args["password"]; ok {
			v, err := parser.GetString("password", false)
			if err != nil {
				return mcp.NewToolResultErrorFromErr("invalid password parameter", err), nil
			}
			password = &v
		}

		var baseURL *string
		if _, ok := args["baseURL"]; ok {
			v, err := parser.GetString("baseURL", false)
			if err != nil {
				return mcp.NewToolResultErrorFromErr("invalid baseURL parameter", err), nil
			}
			baseURL = &v
		}

		err = s.cli.UpdateRegistry(id, name, url, authentication, username, password, baseURL)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to update registry", err), nil
		}

		return mcp.NewToolResultText("Registry updated successfully"), nil
	}
}

func (s *PortainerMCPServer) HandleDeleteRegistry() server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		parser := toolgen.NewParameterParser(request)

		id, err := parser.GetInt("id", true)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("invalid id parameter", err), nil
		}

		err = s.cli.DeleteRegistry(id)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("failed to delete registry", err), nil
		}

		return mcp.NewToolResultText("Registry deleted successfully"), nil
	}
}
