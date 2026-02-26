---
title: Architecture
nav_order: 6
---

# Architecture
{: .no_toc }

Technical architecture of the Portainer MCP server.
{: .fs-6 .fw-300 }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## High-Level Overview

```
┌────────────────────┐         ┌─────────────────────────────────┐         ┌──────────────┐
│  AI Assistant       │  stdio  │  Portainer MCP Server           │  HTTPS  │  Portainer   │
│  (Claude, Copilot,  │◄───────►│                                 │◄───────►│  API         │
│   Cursor, etc.)     │ JSON-RPC│  ┌─────────────────────────┐   │         │  v2.31.2     │
└────────────────────┘         │  │  Meta-Tool Layer (15)    │   │         └──────────────┘
                               │  │  action → handler route  │   │
                               │  └──────────┬──────────────┘   │
                               │             │                   │
                               │  ┌──────────▼──────────────┐   │
                               │  │  Handler Layer (~98)     │   │
                               │  │  HandleGetEnvironments   │   │
                               │  │  HandleCreateStack       │   │
                               │  │  HandleDockerProxy       │   │
                               │  │  ...                     │   │
                               │  └──────────┬──────────────┘   │
                               │             │                   │
                               │  ┌──────────▼──────────────┐   │
                               │  │  Wrapper Client          │   │
                               │  │  pkg/portainer/client    │   │
                               │  └──────────┬──────────────┘   │
                               │             │                   │
                               │  ┌──────────▼──────────────┐   │
                               │  │  Raw Client (SDK)        │   │
                               │  │  client-api-go/v2        │   │
                               │  └─────────────────────────┘   │
                               └─────────────────────────────────┘
```

## Project Structure

```
portainer-mcp/
├── cmd/portainer-mcp/
│   └── mcp.go                    # Entry point, CLI flags, server bootstrap
├── internal/
│   ├── mcp/
│   │   ├── server.go             # Server core, PortainerClient interface
│   │   ├── metatool_registry.go  # 15 meta-tool definitions
│   │   ├── metatool_handler.go   # Meta-tool routing logic
│   │   ├── metatool_test.go      # Meta-tool tests
│   │   ├── schema.go             # Tool constants, validation helpers
│   │   ├── environment.go        # Environment handlers
│   │   ├── stack.go              # Stack handlers
│   │   ├── docker.go             # Docker proxy handler
│   │   ├── kubernetes.go         # Kubernetes proxy handler
│   │   └── ...                   # 22 handler files total
│   ├── tooldef/
│   │   └── tooldef.go            # Embeds tools.yaml at build time
│   └── k8sutil/
│       └── stripper.go           # K8s response metadata stripping
├── pkg/
│   ├── toolgen/
│   │   ├── yaml.go               # YAML parser for tool definitions
│   │   └── param.go              # Parameter extraction helpers
│   └── portainer/
│       ├── client/
│       │   └── adapter.go        # Wrapper client over SDK
│       └── models/
│           └── *.go              # Local MCP models with conversions
├── tests/integration/            # Integration tests (Docker-based)
├── docs/                         # Documentation (this site)
├── tools.yaml                    # 98 granular tool definitions
├── Makefile                      # Build, test, release targets
└── go.mod                        # Go module definition
```

## Component Details

### Entry Point (`cmd/portainer-mcp/mcp.go`)

The entry point handles:
1. **CLI flag parsing** — all 7 flags including `-server`, `-token`, `-read-only`, `-granular-tools`
2. **Client construction** — creates the Portainer SDK client with the configured server URL and token
3. **Server construction** — builds the MCP server with functional options (`WithClient()`, `WithReadOnly()`, `WithGranularTools()`)
4. **Tool registration** — conditionally registers either 15 meta-tools or 98 granular tools
5. **Stdio transport** — starts the MCP server on stdin/stdout using JSON-RPC 2.0

### Server Core (`internal/mcp/server.go`)

The `PortainerMCPServer` struct holds:
- **`PortainerClient` interface** — ~170 methods covering the entire Portainer API
- **`readOnly` flag** — controls which tools are registered
- **`granularTools` flag** — controls meta-tools vs granular registration

The server uses the **functional options pattern** for configuration:

```go
server := mcp.NewPortainerMCPServer(
    mcp.WithClient(client),
    mcp.WithReadOnly(true),
    mcp.WithGranularTools(false),
)
```

### Meta-Tool Layer (`internal/mcp/metatool_registry.go` + `metatool_handler.go`)

The meta-tool system consists of:

1. **Registry** — `metaToolDefinitions()` returns 15 `metaToolDef` structs, each containing:
   - Tool name and description
   - List of `metaAction` entries mapping action names to handler functions
   - MCP annotation (readOnly, destructive, idempotent, openWorld hints)

2. **Registration** — `RegisterMetaTools()` iterates definitions, filters actions by read-only mode, and programmatically builds MCP tools using the SDK

3. **Routing** — `makeMetaHandler()` creates a closure that extracts the `action` parameter from the request, looks up the handler, and delegates

### Two-Layer Client Architecture

See the [Clients & Models]({% link clients_and_models.md %}) page for complete details.

**Raw Client** (`portainer/client-api-go/v2`):
- Auto-generated Go SDK from Portainer's Swagger/OpenAPI spec
- Direct HTTP communication with Portainer API
- Works with raw API models

**Wrapper Client** (`pkg/portainer/client/adapter.go`):
- Simplifies the SDK interface for the MCP server
- Handles data transformation between raw API models and local MCP models
- Provides consistent error handling

**Local Models** (`pkg/portainer/models/`):
- Simplified structs tailored for MCP tool responses
- Conversion functions from raw API models (e.g., `NewEnvironmentFromAPI()`)
- Only include fields relevant to the MCP use case

### Tool Definitions (`tools.yaml`)

The `tools.yaml` file defines all 98 granular tools with:
- Tool name and description
- Parameters with types, descriptions, and required flags
- MCP annotations (readOnly, destructive, idempotent, openWorld)

This file is embedded in the binary at build time via `internal/tooldef/tooldef.go` and the Go `embed` directive. An external file can override the embedded version.

{: .note }
> Meta-tools are defined programmatically in Go, not in YAML. This allows dynamic action enum filtering based on the read-only mode at runtime.

## MCP Protocol Details

### Transport

The server uses **stdio transport** — it reads JSON-RPC 2.0 requests from stdin and writes responses to stdout. This is the standard MCP transport for local tool integrations.

### Tool Annotations

Every tool (both meta and granular) includes MCP annotations:

| Annotation | Purpose |
|:-----------|:--------|
| `readOnlyHint` | `true` if the tool only reads data |
| `destructiveHint` | `true` if the tool can delete or irreversibly modify data |
| `idempotentHint` | `true` if calling multiple times produces the same result |
| `openWorldHint` | `true` if the tool interacts with external systems (Docker/K8s APIs) |

These annotations help AI assistants make safer decisions about when to confirm actions with the user.

### Response Format

All tool handlers return JSON-serialized results as MCP text content. Error responses include descriptive messages with context about what failed.

## Design Decisions

See the [Design Summary]({% link design_summary.md %}) for the complete record of architectural decisions, including:

- External tools file for maintainability
- Tools-based resource access over MCP resources
- Two-layer client architecture
- Read-only mode implementation
- Meta-tools grouping strategy
