# tmux-mcp

> **Mirror** -- Canonical development lives in [hairglasses-studio/dotfiles](https://github.com/hairglasses-studio/dotfiles) at `mcp/tmux-mcp/`. This repo is a publish mirror kept in parity for `go install` and MCP registry discovery.

[![Go Reference](https://pkg.go.dev/badge/github.com/hairglasses-studio/tmux-mcp.svg)](https://pkg.go.dev/github.com/hairglasses-studio/tmux-mcp)
[![Go Report Card](https://goreportcard.com/badge/github.com/hairglasses-studio/tmux-mcp)](https://goreportcard.com/report/github.com/hairglasses-studio/tmux-mcp)
[![CI](https://github.com/hairglasses-studio/tmux-mcp/actions/workflows/ci.yml/badge.svg)](https://github.com/hairglasses-studio/tmux-mcp/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Glama](https://glama.ai/mcp/servers/hairglasses-studio/tmux-mcp/badges/score.svg)](https://glama.ai/mcp/servers/hairglasses-studio/tmux-mcp)

MCP server for tmux session, window, and pane management. Gives AI assistants the ability to create terminal workspaces, send commands, and capture output — including a declarative workspace tool that sets up multi-pane layouts in a single call.

Built with [mcpkit](https://github.com/hairglasses-studio/mcpkit) using stdio transport.

## Install

```bash
go install github.com/hairglasses-studio/tmux-mcp@latest
```

Or build from source:

```bash
git clone https://github.com/hairglasses-studio/tmux-mcp
cd tmux-mcp
make build
make check
```

## Configure

Add to your MCP client config (for example Codex or Claude Code):

```json
{
  "mcpServers": {
    "tmux": {
      "command": "tmux-mcp"
    }
  }
}
```

For a local checkout, the repo now also ships `.mcp.json` plus a repo-local
launcher script so MCP clients can attach directly without reconstructing the
command manually.

## Tools

### Session Management
| Tool | Description |
|------|-------------|
| `tmux_list_sessions` | List all sessions with window count, creation time, attached status |
| `tmux_new_session` | Create a detached session with optional start directory and command |
| `tmux_kill_session` | Kill a session by name |

### Window Management
| Tool | Description |
|------|-------------|
| `tmux_list_windows` | List windows with layout, pane count, active status |
| `tmux_new_window` | Create a window with optional name and command |

### Pane Management
| Tool | Description |
|------|-------------|
| `tmux_list_panes` | List panes with dimensions, running command, active status |
| `tmux_capture_pane` | Capture visible text from a pane (default 50 lines) |
| `tmux_send_keys` | Send keystrokes (literal text or tmux key names like `Enter`, `C-c`) |
| `tmux_wait_for_text` | Poll a pane until specific text appears or timeout |
| `tmux_search_panes` | Search across all panes in a session for a regex pattern |

### Composed
| Tool | Description |
|------|-------------|
| `tmux_workspace` | Create a multi-window, multi-pane workspace from a declarative spec |

## Usage Examples

The `tmux_workspace` tool replaces sequences of `new_session` + `new_window` + `send_keys` calls with a single declarative specification:

```
"Set up a dev workspace with an editor, server, and log watcher"
→ tmux_workspace(
    session: "dev",
    windows: [
      { name: "editor", command: "nvim ." },
      { name: "server", command: "go run ./cmd/server" },
      { name: "logs", command: "tail -f /var/log/app.log" }
    ]
  )
```

## Key Patterns

- **Graceful degradation**: List operations return empty arrays when no tmux server is running (no errors)
- **Target strings**: Built as `session:window.pane` with optional components
- **Detached by default**: All sessions created with `-d` flag

## License

MIT
