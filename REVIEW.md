# Review Guidelines — tmux-mcp

Inherits from org-wide [REVIEW.md](https://github.com/hairglasses-studio/.github/blob/main/REVIEW.md).

## Additional Focus
- **Session state races**: Multiple tools may target the same session — use proper locking
- **Pane capture bounds**: Validate line ranges before capture, handle empty panes
- **Workspace declarations**: Validate YAML/JSON workspace specs before applying
- **Shell injection**: Sanitize all inputs passed to `tmux send-keys`
