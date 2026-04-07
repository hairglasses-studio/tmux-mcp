# Changelog

Format based on [Keep a Changelog](https://keepachangelog.com/).

## [Unreleased]

### Added
- `tmux_wait_for_text` tool — poll a pane until specific text appears
- `tmux_search_panes` tool — regex search across all panes in a session

### Changed
- Expanded .gitignore, SECURITY.md, CONTRIBUTING.md, Makefile
- Added Go Report Card, pkg.go.dev, CI badges to README
- Consolidated goreleaser configs, fixed golangci-lint config
- Fixed errcheck lint issues across codebase
- Server card capabilities now reflect resources and prompts support

## [1.0.0] - 2026-04-04

### Added
- 9 MCP tools for tmux session, window, and pane management
- Declarative `tmux_workspace` composed tool
- Graceful handling of no-server-running state
- Structured logging via slog
- MIT license and documentation
