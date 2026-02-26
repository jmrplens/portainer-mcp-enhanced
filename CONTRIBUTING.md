# Contributing to Portainer MCP Server

Thank you for your interest in contributing! This guide will help you get started.

## Development Setup

### Prerequisites

- **Go 1.24+** — [Install Go](https://go.dev/doc/install)
- **Make** — Build automation
- **Docker** — Required for integration tests
- A running **Portainer 2.31.2** instance (for live tests)

### Clone & Build

```bash
git clone https://github.com/portainer/portainer-mcp.git
cd portainer-mcp
make build
```

### Run Tests

```bash
make test                  # Unit tests only
make test-integration      # Integration tests (needs Docker)
make test-all              # Everything
make test-coverage         # Coverage report → coverage.out
```

### MCP Inspector

Test your changes interactively with the MCP Inspector:

```bash
make inspector
```

This launches a web UI where you can invoke tools and inspect responses.

## Code Style

### General Conventions

- **Go naming**: `PascalCase` for exported identifiers, `camelCase` for private
- **Error handling**: Always wrap errors with context: `fmt.Errorf("failed to get stack: %w", err)`
- **Imports**: Group in order — standard library, external packages, internal packages
- **Comments**: Document all exported functions with `Parameters` and `Returns` sections

### Package Structure

```
cmd/portainer-mcp/       Entry point and CLI
internal/mcp/            MCP server, handlers, tool registration
internal/tooldef/        Tool definitions (embedded YAML)
pkg/toolgen/             YAML parser and parameter utilities
pkg/portainer/client/    Portainer API wrapper client
pkg/portainer/models/    Local models with conversions
tests/                   Integration and live tests
```

### Models & Client Architecture

The project uses a two-layer model architecture:

- **Raw Models** (`portainer/client-api-go/v2/pkg/models`) — Direct SDK types, prefixed with `raw`
- **Local Models** (`pkg/portainer/models`) — Simplified types for MCP use

```go
import (
    "github.com/portainer/portainer-mcp/pkg/portainer/models"              // Local
    apimodels "github.com/portainer/client-api-go/v2/pkg/models"           // Raw SDK
)
```

Always convert from raw → local models before returning to handlers.

## Adding a New Tool

### 1. Define the Tool in `tools.yaml`

```yaml
- name: myNewTool
  description: "Does something useful"
  parameters:
    type: object
    properties:
      id:
        type: integer
        description: "Resource ID"
    required:
      - id
  annotations:
    title: "My New Tool"
    readOnlyHint: true        # Set to true if read-only
    destructiveHint: false
    idempotentHint: true
    openWorldHint: true
```

> **Important**: Use `parameters:` not `inputSchema:` — the parser only recognizes `parameters`.

### 2. Add the Handler

Create or extend a handler file in `internal/mcp/`:

```go
func (s *PortainerMCPServer) handleMyNewTool(arguments map[string]interface{}) (*mcp.CallToolResult, error) {
    params := toolgen.NewParameterParser(arguments)

    id, err := params.GetRequiredInt("id")
    if err != nil {
        return nil, fmt.Errorf("invalid parameters: %w", err)
    }

    result, err := s.client.MyNewMethod(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get resource: %w", err)
    }

    return toolgen.TextResult(result)
}
```

### 3. Register the Handler

In the server's feature registration (`internal/mcp/` or `cmd/portainer-mcp/mcp.go`):

```go
func (s *PortainerMCPServer) AddMyFeatures() {
    s.AddToolHandler("myNewTool", s.handleMyNewTool, false)
    // Use `true` for the third parameter if it's a read-only tool
}
```

### 4. Implement the Client Method

1. Add to `PortainerClient` interface in `internal/mcp/server.go`
2. Implement in `pkg/portainer/client/`
3. Add local model in `pkg/portainer/models/` if needed

### 5. Write Tests

- **Unit test**: Mock the client interface, test handler logic
- **Integration test**: Use `testcontainers-go` for end-to-end validation

## Testing Patterns

### Unit Tests

Use table-driven tests:

```go
func TestHandleMyNewTool(t *testing.T) {
    tests := []struct {
        name    string
        args    map[string]interface{}
        want    string
        wantErr bool
    }{
        {
            name:    "valid ID",
            args:    map[string]interface{}{"id": float64(1)},
            want:    `"id":1`,
            wantErr: false,
        },
        {
            name:    "missing ID",
            args:    map[string]interface{}{},
            wantErr: true,
        },
    }
    // ...
}
```

### Integration Tests

Integration tests use Docker to spin up real Portainer instances:

```go
func TestMyFeature(t *testing.T) {
    env := helpers.SetupTestEnv(t)
    defer env.Cleanup()
    // Test against real Portainer
}
```

## Commit Messages

Use conventional commit format:

```
feat: add webhook management tools
fix: correct tools.yaml schema keys for Helm tools
docs: update README with architecture diagram
test: add integration tests for backup operations
refactor: simplify parameter parsing in handlers
```

## Pull Request Process

1. **Fork** the repository and create a branch from `main`
2. **Implement** your changes with tests
3. **Run** `make test-all` to verify nothing is broken
4. **Format** code: `gofmt -s -w .`
5. **Submit** a PR with a clear description of what and why

### PR Checklist

- [ ] Tests pass (`make test-all`)
- [ ] Code is formatted (`gofmt`)
- [ ] New tools are added to `tools.yaml` with correct `parameters:` key
- [ ] Handler includes proper error wrapping
- [ ] Read-only annotations are set correctly
- [ ] Documentation updated if adding user-facing changes

## Design Decisions

Significant architectural decisions are documented in `docs/design/`. Before making major changes, review existing decisions and create a new record if needed:

- Naming: `YYMMDD-N-short-description.md`
- Template: See `docs/design_summary.md`

## Questions?

Open an issue or discussion on GitHub. We're happy to help!
