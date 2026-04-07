package main

import (
	"context"
	"fmt"

	"github.com/hairglasses-studio/mcpkit/prompts"
	"github.com/hairglasses-studio/mcpkit/resources"
	"github.com/mark3labs/mcp-go/mcp"
)

type tmuxResourceModule struct{}

func (m *tmuxResourceModule) Name() string        { return "tmux_context" }
func (m *tmuxResourceModule) Description() string { return "Reusable tmux debugging context" }

func (m *tmuxResourceModule) Resources() []resources.ResourceDefinition {
	return []resources.ResourceDefinition{
		{
			Resource: mcp.NewResource(
				"tmux://workflows/session-debug",
				"Tmux Session Debug Workflow",
				mcp.WithResourceDescription("Compact sequence for inspecting a tmux session without spamming pane output"),
				mcp.WithMIMEType("text/markdown"),
			),
			Handler: func(_ context.Context, _ mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
				return []mcp.ResourceContents{
					mcp.TextResourceContents{
						URI:      "tmux://workflows/session-debug",
						MIMEType: "text/markdown",
						Text:     "1. Use `tmux_list_sessions`, `tmux_list_windows`, and `tmux_list_panes` to target the right pane.\n2. Use `tmux_capture_pane` with a bounded line count before asking for more output.\n3. Use `tmux_search_panes` when the issue is pattern-driven.\n4. Use `tmux_send_keys` only after identifying the target pane and the exact command to send.",
					},
				}, nil
			},
			Category: "workflow",
			Tags:     []string{"tmux", "debugging", "workflow"},
		},
	}
}

func (m *tmuxResourceModule) Templates() []resources.TemplateDefinition { return nil }

type tmuxPromptModule struct{}

func (m *tmuxPromptModule) Name() string        { return "tmux_prompts" }
func (m *tmuxPromptModule) Description() string { return "Prompt workflows for tmux investigations" }

func (m *tmuxPromptModule) Prompts() []prompts.PromptDefinition {
	return []prompts.PromptDefinition{
		{
			Prompt: mcp.NewPrompt(
				"tmux_debug_session",
				mcp.WithPromptDescription("Investigate a tmux session with bounded pane inspection before sending input"),
				mcp.WithArgument("session", mcp.RequiredArgument(), mcp.ArgumentDescription("Tmux session name")),
				mcp.WithArgument("goal", mcp.ArgumentDescription("What you are trying to inspect or fix")),
			),
			Handler: func(_ context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
				session := req.Params.Arguments["session"]
				goal := req.Params.Arguments["goal"]
				if goal == "" {
					goal = "inspect the current state"
				}
				return &mcp.GetPromptResult{
					Description: "Debug tmux session " + session,
					Messages: []mcp.PromptMessage{
						mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(fmt.Sprintf(
							"Investigate tmux session %q to %s. Start with `tmux_list_windows` and `tmux_list_panes`, then use `tmux_capture_pane` or `tmux_search_panes` with bounded output. Only use `tmux_send_keys` after identifying the exact target pane and command.",
							session, goal,
						))),
					},
				}, nil
			},
			Category: "workflow",
			Tags:     []string{"tmux", "workflow", "debugging"},
		},
	}
}

func buildTmuxResourceRegistry() *resources.ResourceRegistry {
	reg := resources.NewResourceRegistry()
	reg.RegisterModule(&tmuxResourceModule{})
	return reg
}

func buildTmuxPromptRegistry() *prompts.PromptRegistry {
	reg := prompts.NewPromptRegistry()
	reg.RegisterModule(&tmuxPromptModule{})
	return reg
}
