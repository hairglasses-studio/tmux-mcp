# tmux-mcp

> **Archived** -- This repo has been consolidated into [hairglasses-studio/dotfiles](https://github.com/hairglasses-studio/dotfiles) at `mcp/tmux-mcp/`. For continued updates, use the consolidated version.

[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Glama](https://glama.ai/mcp/servers/hairglasses-studio/tmux-mcp/badges/score.svg)](https://glama.ai/mcp/servers/hairglasses-studio/tmux-mcp)

MCP server for tmux session, window, and pane management. Gives AI assistants the ability to create terminal workspaces, send commands, and capture output — including a declarative workspace tool that sets up multi-pane layouts in a single call.

Canonical development lives in [`hairglasses-studio/dotfiles`](https://github.com/hairglasses-studio/dotfiles/tree/main/mcp/tmux-mcp) under `dotfiles/mcp/tmux-mcp`. The standalone [`tmux-mcp`](https://github.com/hairglasses-studio/tmux-mcp) repo is a publish mirror kept in parity for installation and discovery.

Built with [mcpkit](https://github.com/hairglasses-studio/mcpkit) using stdio transport.

## Install

```bash
go install github.com/hairglasses-studio/tmux-mcp@latest
```

Or build from source:

```bash
git clone https://github.com/hairglasses-studio/tmux-mcp
cd tmux-mcp
go build -o tmux-mcp .
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

### Composed
| Tool | Description |
|------|-------------|
| `tmux_workspace` | Create a multi-window, multi-pane workspace from a declarative spec |

## Declarative Workspace

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
