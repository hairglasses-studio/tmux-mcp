#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$REPO_ROOT"
export GOWORK=off

# Default is official since mcpkit v0.8.2 (P78.38 schema_reflect fix; wire-verified
# 28/28 described props, equal to legacy). TMUX_MCP_SDK=legacy still selects the
# legacy build explicitly.
case "${TMUX_MCP_SDK:-official}" in
  official)
    exec go run -tags=official_sdk ./cmd/tmux-mcp
    ;;
  legacy)
    exec go run ./cmd/tmux-mcp
    ;;
  *)
    printf 'unsupported TMUX_MCP_SDK=%q (expected official or legacy)\n' "$TMUX_MCP_SDK" >&2
    exit 2
    ;;
esac
