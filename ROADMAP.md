# Roadmap

## Current State

tmux-mcp provides 9 tools for tmux session, window, and pane management via MCP. Includes a declarative `tmux_workspace` composed tool that creates multi-window, multi-pane layouts in a single call. Graceful handling of no-server-running state. Built on mcpkit with stdio transport.

All tools functional and tested. MIT licensed, README and CLAUDE.md in place.

## Planned

### Phase 1 — Robustness & Testing
- Add integration tests using `mcptest.NewServer()` with mock tmux state
- Improve `tmux_capture_pane` with scroll-back support (configurable history depth)
- Add `tmux_resize_pane` tool for programmatic pane resizing
- Better error messages when target session/window/pane doesn't exist

### Phase 2 — Workspace Templates
- Named workspace templates (save and recall `tmux_workspace` specs)
- `tmux_snapshot` — capture current layout as a reusable workspace spec
- `tmux_workspace` support for split direction (horizontal/vertical) and size ratios

### Phase 3 — Session Intelligence
- `tmux_search` — search pane contents across all sessions
- `tmux_wait_for_output` — block until a pane produces matching output (useful for build/test workflows)

## Future Considerations
- SSE transport for real-time pane output streaming
- Integration with tmux hooks for event-driven notifications
- Workspace persistence across tmux server restarts (complement tmux-resurrect)
