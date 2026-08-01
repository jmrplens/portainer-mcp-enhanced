<div align="center">

# Portainer MCP Server (Enhanced) — archived

**This repository has moved to [jmrplens/portainer-mcp](https://github.com/jmrplens/portainer-mcp).**

</div>

---

## What happened

This project was a fork of the [official Portainer MCP Server](https://github.com/portainer/portainer-mcp), extended to 98 tools. In August 2026 it was rewritten from scratch as a standalone project rather than a fork, and development continues at:

### → **https://github.com/jmrplens/portainer-mcp**

The rewrite is not an incremental change. It replaces the entire Go codebase:

| | This repository (archived) | [jmrplens/portainer-mcp](https://github.com/jmrplens/portainer-mcp) |
|---|---|---|
| API coverage | 98 tools | Full 1:1 with the Portainer REST API — 265 operations on CE, 442 on EE — verified in CI |
| Editions | No edition awareness | CE and EE detected at runtime; EE-only operations gated |
| Portainer versions | Pinned to 2.39.1, refused to start otherwise | LTS and STS channels both supported, with per-operation version gating |
| Tool surfaces | 15 meta-tools | `dynamic` (2 tools), `meta` and `individual`, all projected from one action catalog |
| MCP SDK | `mark3labs/mcp-go` | Official `modelcontextprotocol/go-sdk`, with stateless streamable HTTP |
| Client | Hand-written adapters over a partial SDK | Generated from Portainer's official OpenAPI specification |

## Why a new repository rather than a rename

This repository is a fork, and GitHub excludes forks from search results. Detaching it was not possible: GitHub refuses to remove a fork from its network while it has child forks, and this one has three. Creating a standalone repository was the only route to an independent project — and it leaves this one intact as a reference instead of destroying its history, issues and stars, which detaching would have done.

## This repository is read-only

It is archived. Issues and pull requests belong on the new repository. The full git history and every release tag remain here for anyone who needs them.

Installations pinned to `github.com/jmrplens/portainer-mcp-enhanced` continue to resolve — nothing published has been removed — but they receive no further updates, **including security fixes**. Report vulnerabilities against the new repository.

## License

MIT, unchanged. See [LICENSE](LICENSE).
