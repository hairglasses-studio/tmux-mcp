# tmux-mcp — Gemini CLI Instructions

MCP server for tmux session, window, and pane management. Built with [mcpkit](https://github.com/hairglasses-studio/mcpkit).

## Build & Test

```bash
go build ./...
go vet ./...
go test ./... -count=1
go install .
```

## Key Conventions

- Graceful handling of "no server running" -- list operations return empty arrays
- Target strings built as `session:window.pane` with optional window/pane components
- All sessions created detached (`-d` flag)

