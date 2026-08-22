#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$REPO_ROOT"
export GOWORK=off

# DEFAULT REVERTED to legacy (2026-08-22): the pinned mcpkit (v0.8.x) predates
# the P78.38 schema fix, so the official_sdk build serves ZERO property
# descriptions and DROPS every omitempty input param (live-verified: the
# legacy build serves 22/22 described props incl. confirm; official served
# none). Re-default to official once mcpkit v0.8.2 (schema_reflect) is
# tagged and the pin bumped. TMUX_MCP_SDK=official still selects it explicitly.
case "${TMUX_MCP_SDK:-legacy}" in
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
