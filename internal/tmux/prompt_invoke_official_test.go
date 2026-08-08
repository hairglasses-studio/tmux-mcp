//go:build official_sdk

// prompt_invoke_official_test.go is the official_sdk counterpart to
// prompt_invoke_mcpgo_test.go — same callPrompt signature, adapted to the
// official SDK's pointer-typed *mcp.GetPromptRequest/GetPromptParams.
package tmux

import (
	"context"

	"github.com/hairglasses-studio/mcpkit/prompts"
	"github.com/hairglasses-studio/mcpkit/registry"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func callPrompt(handler prompts.PromptHandlerFunc, args map[string]string) (*registry.GetPromptResult, error) {
	req := &mcp.GetPromptRequest{Params: &mcp.GetPromptParams{Arguments: args}}
	return handler(context.Background(), req)
}
