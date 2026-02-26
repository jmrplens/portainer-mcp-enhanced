---
title: Contributing
nav_order: 8
---

# Contributing
{: .no_toc }

How to contribute to the Portainer MCP project.
{: .fs-6 .fw-300 }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Development Setup

### Prerequisites

- **Go 1.24+**
- **Make**
- **Docker** (for integration tests)
- A running **Portainer 2.31.2** instance (for integration tests)

### Clone and Build

```bash
git clone https://github.com/portainer/portainer-mcp.git
cd portainer-mcp
make build
```

### Run Tests

```bash
# Unit tests
go test -v ./...

# Single test
go test -v ./internal/mcp/ -run TestMetaToolCount

# Coverage report
make test-coverage

# Integration tests (requires Docker + Portainer)
make test-integration

# All tests
make test-all
```

### Code Quality

```bash
# Format code
make fmt

# Static analysis
make vet

# Lint (requires golint)
make lint
```

## Code Style

- **Naming**: PascalCase for exported, camelCase for private
- **Errors**: wrap with context via `fmt.Errorf("failed to do X: %w", err)`
- **Imports**: group stdlib, external, internal — separated by blank lines
- **Comments**: document all exported functions with Parameters/Returns sections
- **Tests**: use table-driven pattern with descriptive case names

## Adding a New Tool

### 1. Define the Tool in `tools.yaml`

Add a new entry to `tools.yaml` with the tool name, description, parameters, and annotations:

```yaml
- name: my_new_tool
  description: "Does something useful"
  annotations:
    readOnlyHint: true
    destructiveHint: false
    idempotentHint: true
    openWorldHint: false
  parameters:
    - name: id
      type: integer
      description: "Resource ID"
      required: true
```

### 2. Add a Tool Constant

In `internal/mcp/schema.go`, add a constant:

```go
const ToolMyNewTool = "my_new_tool"
```

### 3. Create a Handler

In the appropriate handler file (or a new one), create the handler function:

```go
func (s *PortainerMCPServer) HandleMyNewTool(
    ctx context.Context,
    request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
    // Extract parameters
    id, err := toolgen.GetInt(request, "id")
    if err != nil {
        return nil, fmt.Errorf("failed to get id: %w", err)
    }

    // Call Portainer API
    result, err := s.client.GetSomething(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get something: %w", err)
    }

    // Serialize and return
    data, err := json.Marshal(result)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal result: %w", err)
    }

    return mcp.NewToolResultText(string(data)), nil
}
```

### 4. Register the Handler

In the `AddXxxFeatures` function for the relevant domain:

```go
func (s *PortainerMCPServer) AddMyFeatures(mcpServer *server.MCPServer, tools []toolgen.ToolDef, readOnly bool) {
    for _, t := range tools {
        switch t.Name {
        case ToolMyNewTool:
            mcpServer.AddTool(t.Tool, s.HandleMyNewTool)
        }
    }
}
```

### 5. Add to Meta-Tool Group

In `internal/mcp/metatool_registry.go`, add the action to the appropriate meta-tool:

```go
{name: "my_new_action", handler: (*PortainerMCPServer).HandleMyNewTool, readOnly: true},
```

### 6. Test

Write table-driven tests for your handler and verify it works with both meta-tools and granular mode.

## Architecture Guidelines

- **Handlers** should validate inputs, call the client, and return JSON results
- **Client methods** transform between raw API models and local MCP models
- **Local models** should only include fields relevant to the MCP use case
- **Design decisions** should be documented in `docs/design/`

See the [Architecture]({% link architecture.md %}) page for the full system overview.

## Pull Request Process

1. Fork the repository
2. Create a feature branch from `main`
3. Make your changes with tests
4. Run `make fmt && make vet && go test ./...`
5. Submit a pull request with a clear description

## Documentation

- Update existing docs if your changes affect them
- Add new documentation for new features
- Design decisions go in `docs/design/YYMMDD-N-short-description.md`
