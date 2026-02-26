---
title: Getting Started
nav_order: 2
---

# Getting Started
{: .no_toc }

Get Portainer MCP running with your AI assistant in under 5 minutes.
{: .fs-6 .fw-300 }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

---

## Prerequisites

- A running **Portainer 2.31.2** instance (Business Edition or Community Edition)
- An **API token** from your Portainer instance
- An MCP-compatible AI assistant (Claude Desktop, VS Code with Copilot, Cursor, etc.)

## Installation

### Download a Pre-built Binary

Grab the latest binary from [GitHub Releases](https://github.com/portainer/portainer-mcp/releases/latest) for your platform:

| Platform | Architecture | Binary |
|:---------|:------------|:-------|
| Linux | amd64 | `portainer-mcp-linux-amd64` |
| Linux | arm64 | `portainer-mcp-linux-arm64` |
| macOS | arm64 (Apple Silicon) | `portainer-mcp-darwin-arm64` |

```bash
# Example: Linux amd64
curl -L -o portainer-mcp \
  https://github.com/portainer/portainer-mcp/releases/latest/download/portainer-mcp-linux-amd64
chmod +x portainer-mcp
```

### Build from Source

Requires **Go 1.24+** and **Make**:

```bash
git clone https://github.com/portainer/portainer-mcp.git
cd portainer-mcp
make build
# Binary: dist/portainer-mcp
```

Cross-compile for a specific platform:

```bash
make PLATFORM=linux ARCH=arm64 build
```

## Get a Portainer API Token

1. Log in to your Portainer instance
2. Go to **My Account** → **API Keys**
3. Click **Add API Key**, give it a name
4. Copy the generated token (starts with `ptr_`)

{: .warning }
> Keep your API token secure. It provides full access to your Portainer instance unless you use `-read-only` mode.

## Configure Your AI Assistant

### Claude Desktop

Edit your Claude Desktop configuration file:

- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "portainer": {
      "command": "/path/to/portainer-mcp",
      "args": [
        "-server", "https://your-portainer:9443",
        "-token", "ptr_your_api_token"
      ]
    }
  }
}
```

### VS Code (GitHub Copilot)

Create `.vscode/mcp.json` in your workspace:

```json
{
  "servers": {
    "portainer": {
      "type": "stdio",
      "command": "/path/to/portainer-mcp",
      "args": [
        "-server", "https://your-portainer:9443",
        "-token", "ptr_your_api_token"
      ]
    }
  }
}
```

### Cursor

Go to **Cursor Settings → MCP** and add:

```json
{
  "mcpServers": {
    "portainer": {
      "command": "/path/to/portainer-mcp",
      "args": [
        "-server", "https://your-portainer:9443",
        "-token", "ptr_your_api_token"
      ]
    }
  }
}
```

## Verify It Works

Once configured, ask your AI assistant:

> "What environments are available in Portainer?"

You should see the assistant call `manage_environments` with `action: list_environments` and return your environment list.

### More Examples

```
"Show me the Docker dashboard for environment 1"
"List all stacks and their status"
"Create a new team called 'developers'"
"Search for nginx Helm charts"
"What's the system status of Portainer?"
```

## Next Steps

- [Configuration]({% link configuration.md %}) — all CLI flags and options
- [Meta-Tools Guide]({% link meta-tools.md %}) — understand the 15 grouped tools
- [Tools Reference]({% link api-reference.md %}) — complete parameter details for all 98 tools
- [Security]({% link security.md %}) — security considerations and read-only mode
