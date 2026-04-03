# tmux-mcp

MCP server for tmux session, window, and pane management. Built with [mcpkit](https://github.com/hairglasses-studio/mcpkit).

## Build & Test
```bash
go build ./...
go vet ./...
go test ./... -count=1
go install .
```

## Tools (9)

### Session Management (3)
- `tmux_list_sessions` -- List all tmux sessions with window count, creation time, attached status
- `tmux_new_session` -- Create a new detached session with optional start directory and command
- `tmux_kill_session` -- Kill a tmux session by name

### Window Management (2)
- `tmux_list_windows` -- List windows in a session with layout, pane count, active status
- `tmux_new_window` -- Create a new window in a session with optional name and command

### Pane Management (3)
- `tmux_list_panes` -- List panes in a session/window with dimensions, command, active status
- `tmux_capture_pane` -- Capture visible text content from a pane (default 50 lines)
- `tmux_send_keys` -- Send keystrokes to a pane (literal text or tmux key names like Enter, C-c)

### Composed (1)
- `tmux_workspace` -- **Composed**: create multi-window, multi-pane workspace from declarative spec. Replaces new_session + new_window + send_keys sequences.

## Key Patterns
- Graceful handling of "no server running" -- list operations return empty arrays
- Target strings built as `session:window.pane` with optional window/pane components
- All sessions created detached (`-d` flag)
