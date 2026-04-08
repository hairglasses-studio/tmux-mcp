---
name: tmux_mcp
description: 'Operate the tmux-mcp server for tmux session, window, and pane automation. Use this when changing tmux MCP tools or session orchestration in this repo, not for higher-level agent orchestration that belongs in ralphglasses or hgmux.'
---

# tmux-mcp

Use this repo for the dedicated MCP surface around tmux sessions, windows, panes, and related terminal orchestration primitives.

Focus paths:
- `cmd/`
- `internal/`
- `AGENTS.md`
- `README.md`

Keep the tool surface aligned with tmux semantics and avoid baking in higher-level workflow assumptions that belong in consumer repos.
