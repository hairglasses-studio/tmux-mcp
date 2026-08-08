//go:build !official_sdk

// prompt_invoke_mcpgo_test.go provides the mcp-go build's callPrompt test
// helper — prompts.PromptHandlerFunc's request shape differs per SDK
// (mcp-go: func(ctx, mcp.GetPromptRequest) (*mcp.GetPromptResult, error), a
// value-typed Params; official: func(ctx, *mcp.GetPromptRequest) (...), a
// pointer-typed Params — mcpkit/prompts has no dual-tag request constructor
// for this, mirroring its own handler_adapters_test.go split), so this
// helper is tag-split within the consumer repo instead of shared. See
// prompt_invoke_official_test.go for the official_sdk counterpart.
package tmux

import (
	"context"

	"github.com/hairglasses-studio/mcpkit/prompts"
	"github.com/hairglasses-studio/mcpkit/registry"
	"github.com/mark3labs/mcp-go/mcp"
)

func callPrompt(handler prompts.PromptHandlerFunc, args map[string]string) (*registry.GetPromptResult, error) {
	req := mcp.GetPromptRequest{}
	req.Params.Arguments = args
	return handler(context.Background(), req)
}
