# Contributing to tmux-mcp

Thank you for your interest in contributing to tmux-mcp. This document covers
everything you need to get started.

## Reporting Bugs

Open a [GitHub Issue](https://github.com/hairglasses-studio/tmux-mcp/issues)
with:

- Go version, OS, and tmux version
- Minimal reproduction steps
- Expected vs. actual behavior
- Relevant error output or logs

## Suggesting Features

Open a [GitHub Issue](https://github.com/hairglasses-studio/tmux-mcp/issues)
with the `enhancement` label. Describe the use case and which tmux operations
are involved. If suggesting a new MCP tool, include the proposed tool name and
parameter schema.

## Submitting Pull Requests

1. Fork the repository and clone your fork.
2. Create a branch from `main`: `git checkout -b feat/my-change`
3. Make your changes, following the code style below.
4. Run the check suite: `make build && make vet && make test`
5. Commit with a clear, descriptive message.
6. Push your branch and open a PR against `main`.

Keep PRs focused. One logical change per PR is easier to review than a combined
refactor-plus-feature.

## Development Setup

**Requirements:** Go 1.26.1+, tmux installed

```bash
git clone https://github.com/hairglasses-studio/tmux-mcp
cd tmux-mcp
make build    # go build -o tmux-mcp ./...
make test     # go test ./... -count=1
make vet      # go vet ./...
```

## Code Style

- Format with `gofmt` (or `goimports`).
- Pass `go vet ./...` with no warnings.
- Follow existing patterns in the codebase.
- New tools must implement the `ToolModule` interface (`Name()`, `Description()`, `Tools()`).
- Use `handler.TypedHandler` for new tool handlers.
- Return errors via `handler.CodedErrorResult` -- never `(nil, error)`.
- Protect shared state with `sync.RWMutex` (`RLock` for reads, `Lock` for writes).

## Testing Requirements

- All existing tests must pass before submitting a PR.
- Add tests for new features and bug fixes.
- Run tests with race detection: `go test ./... -count=1 -race`
- Integration tests use `mcptest.NewServer()`; unit tests use stdlib `testing`.

## Commit Messages

Use conventional-style prefixes:

```
feat: add declarative workspace layout tool
fix: handle detached session edge case
docs: document pane capture options
test: add tests for window creation
```

## License

By contributing, you agree that your contributions will be licensed under the
[MIT License](LICENSE).

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md).
