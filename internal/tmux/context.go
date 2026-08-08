package tmux

import (
	"context"
	"fmt"

	"github.com/hairglasses-studio/mcpkit/prompts"
	"github.com/hairglasses-studio/mcpkit/registry"
	"github.com/hairglasses-studio/mcpkit/resources"
)

// newDocResource builds a registry-aliased Resource via mcpkit's compat
// constructor plus field assignment, instead of mcp-go's functional options
// (mcp.WithResourceDescription etc.), which the compat layer deliberately
// does not re-export. This file must not import mark3labs/mcp-go (or
// modelcontextprotocol/go-sdk) directly so it compiles under both build
// tags unchanged.
func newDocResource(uri, name, description, mimeType string) registry.Resource {
	r := registry.NewResource(uri, name)
	r.Description = description
	r.MIMEType = mimeType
	return r
}

type tmuxResourceModule struct{}

func (m *tmuxResourceModule) Name() string        { return "tmux_context" }
func (m *tmuxResourceModule) Description() string { return "Reusable tmux debugging context" }

func (m *tmuxResourceModule) Resources() []resources.ResourceDefinition {
	return []resources.ResourceDefinition{
		{
			Resource: newDocResource(
				"tmux://workflows/session-debug",
				"Tmux Session Debug Workflow",
				"Compact sequence for inspecting a tmux session without spamming pane output",
				"text/markdown",
			),
			Handler: resources.TextResourceHandler(func(_ context.Context, _ string) (string, string, error) {
				return "text/markdown", "1. Use `tmux_list_sessions`, `tmux_list_windows`, and `tmux_list_panes` to target the right pane.\n2. Use `tmux_capture_pane` with a bounded line count before asking for more output.\n3. Use `tmux_search_panes` when the issue is pattern-driven.\n4. Use `tmux_send_keys` only after identifying the target pane and the exact command to send.", nil
			}),
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
			Prompt: registry.MakePrompt(
				"tmux_debug_session",
				"Investigate a tmux session with bounded pane inspection before sending input",
				registry.PromptArgument{Name: "session", Description: "Tmux session name", Required: true},
				registry.PromptArgument{Name: "goal", Description: "What you are trying to inspect or fix"},
			),
			Handler: prompts.TextPromptHandler(func(_ context.Context, args map[string]string) (string, string, error) {
				session := args["session"]
				goal := args["goal"]
				if goal == "" {
					goal = "inspect the current state"
				}
				return "Debug tmux session " + session, fmt.Sprintf(
					"Investigate tmux session %q to %s. Start with `tmux_list_windows` and `tmux_list_panes`, then use `tmux_capture_pane` or `tmux_search_panes` with bounded output. Only use `tmux_send_keys` after identifying the exact target pane and command.",
					session, goal,
				), nil
			}),
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
