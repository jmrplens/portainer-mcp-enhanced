---
title: Home
layout: home
nav_order: 1
---

# Portainer MCP Server
{: .fs-9 }

Manage your entire Portainer infrastructure through AI assistants using the Model Context Protocol.
{: .fs-6 .fw-300 }

[Get Started]({% link getting-started.md %}){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 }
[View on GitHub](https://github.com/portainer/portainer-mcp){: .btn .fs-5 .mb-4 .mb-md-0 }

---

## What is Portainer MCP?

Portainer MCP is a [Model Context Protocol](https://modelcontextprotocol.io/) server that connects AI assistants — like **Claude Desktop**, **VS Code Copilot**, and **Cursor** — to your [Portainer](https://www.portainer.io/) instance. It exposes **98 tools** covering the complete Portainer API, enabling natural language management of your container infrastructure.

### Key Features

- **15 grouped meta-tools** (default) for optimal LLM tool selection, or **98 granular tools** for full control
- **Complete Portainer API coverage** — environments, stacks, Docker, Kubernetes, Helm, users, teams, registries, edge computing, backups, and more
- **Read-only mode** — restrict to observation-only operations for safe monitoring
- **MCP annotations** — every tool includes `readOnlyHint`, `destructiveHint`, `idempotentHint`, and `openWorldHint` for safe AI operation
- **Version-pinned compatibility** — validated against specific Portainer versions to prevent API mismatches
- **Zero dependencies** — single static binary, no runtime requirements

### How It Works

```
┌─────────────────────┐      MCP Protocol       ┌─────────────────────┐      HTTPS       ┌───────────────┐
│   AI Assistant       │ ◄──── (stdio/JSON-RPC) ──►│  Portainer MCP      │ ◄──────────────► │  Portainer    │
│  Claude / Copilot    │                          │  Server             │                  │  API          │
│  Cursor / etc.       │                          │  (15 meta-tools)    │                  │  v2.31.2      │
└─────────────────────┘                          └─────────────────────┘                  └───────────────┘
```

The server acts as a bridge, translating MCP tool calls into Portainer API requests and returning structured results that the AI assistant can interpret and present.

### Quick Example

Ask your AI assistant:

> "List all environments and show me the Docker dashboard for the local one"

The assistant will call `manage_environments` with `action: list_environments`, then `manage_docker` with `action: get_docker_dashboard` — all handled automatically.

---

## Compatibility

| MCP Server Version | Supported Portainer |
|:-------------------|:-------------------|
| v0.6.x | 2.31.2 |
| v0.5.x | 2.30.0 |
| v0.4.x | 2.27.4 |

## License

Copyright © 2025 [Portainer.io](https://www.portainer.io/). See [LICENSE](https://github.com/portainer/portainer-mcp/blob/main/LICENSE) for details.
