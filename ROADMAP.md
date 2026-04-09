# Roadmap

## Current State

tmux-mcp provides 11 tools for tmux session, window, and pane management via MCP. Includes a declarative `tmux_workspace` composed tool that creates multi-window, multi-pane layouts in a single call, plus `tmux_wait_for_text` for polling pane output and `tmux_search_panes` for regex search across panes. Graceful handling of no-server-running state. Built on mcpkit with stdio transport.

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

## Future Considerations
- SSE transport for real-time pane output streaming
- Integration with tmux hooks for event-driven notifications
- Workspace persistence across tmux server restarts (complement tmux-resurrect)

<!-- whiteclaw-rollout:start -->
## Whiteclaw-Derived Overhaul (2026-04-08)

This tranche applies the highest-value whiteclaw findings that fit this repo's real surface: engineer briefs, bounded skills/runbooks, searchable provenance, scoped MCP packaging, and explicit verification ladders.

### Strategic Focus
- Treat this repo as a public mirror with a user-facing terminal contract, not just a generic mirror shell.
- Use whiteclaw patterns to make the exported tmux workspace/wait surfaces easy to understand and safe to verify.
- Keep source-of-truth and release parity explicit so mirror maintenance stays mechanical.

### Recommended Work
- [ ] [Mirror contract] Keep the canonical-source mapping to `dotfiles/mcp/tmux-mcp` explicit and verifiable.
- [ ] [Schema/examples] Snapshot and document the exported `tmux_workspace`, `tmux_wait_for_text`, and related surface contracts.
- [ ] [Compatibility tests] Add smoke tests for 'no tmux server running' and other common compatibility paths.
- [ ] [Publish verification] Add mirror smoke tests that prove the released artifact matches the canonical source surface.

### Rationale Snapshot
- Tier / lifecycle: `standalone` / `publish-mirror`
- Language profile: `Go`
- Visibility / sensitivity: `PUBLIC` / `public`
- Surface baseline: AGENTS=yes, skills=yes, codex=yes, mcp_manifest=configured, ralph=yes, roadmap=yes
- Whiteclaw transfers in scope: mirror contract, schema/examples, compatibility smoke tests, publish verification
- Live repo notes: AGENTS, skills, Codex config, configured .mcp.json, .ralph, 1 workflow(s)

<!-- whiteclaw-rollout:end -->
