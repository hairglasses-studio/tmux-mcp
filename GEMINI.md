# tmux-mcp — Gemini CLI Instructions

This repo uses [AGENTS.md](AGENTS.md) as the canonical instruction file. Treat this file as compatibility guidance for Gemini-specific workflows.

MCP server for tmux session, window, and pane management. Built with [mcpkit](https://github.com/hairglasses-studio/mcpkit).

## Build & Test

```bash
go build -o tmux-mcp ./...
go vet ./...
go test ./... -count=1
make check              # All three above
```

## Architecture

Go program (`main.go` + `context.go`). One `TmuxModule` registers all 11 tools. Shells out to `tmux` CLI for all operations.

## Key Conventions

- Graceful handling of "no server running" — list operations return empty arrays.
- Target strings built as `session:window.pane` with optional window/pane components.
- All sessions created detached (`-d` flag).
- Error codes: `handler.CodedErrorResult(handler.ErrInvalidParam, err)` — never `(nil, error)`.
