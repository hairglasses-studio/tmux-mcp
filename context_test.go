package main

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

// ---------------------------------------------------------------------------
// Resource module
// ---------------------------------------------------------------------------

func TestTmuxResourceModule_Metadata(t *testing.T) {
	m := &tmuxResourceModule{}
	if m.Name() != "tmux_context" {
		t.Errorf("Name() = %q, want %q", m.Name(), "tmux_context")
	}
	if m.Description() == "" {
		t.Error("Description() is empty")
	}
}

func TestTmuxResourceModule_Resources(t *testing.T) {
	m := &tmuxResourceModule{}
	resources := m.Resources()
	if len(resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(resources))
	}

	rd := resources[0]
	if rd.Category != "workflow" {
		t.Errorf("Category = %q, want %q", rd.Category, "workflow")
	}
	if len(rd.Tags) == 0 {
		t.Error("expected tags")
	}

	// Call the handler
	contents, err := rd.Handler(context.Background(), mcp.ReadResourceRequest{})
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if len(contents) != 1 {
		t.Fatalf("expected 1 content, got %d", len(contents))
	}
	tc, ok := contents[0].(mcp.TextResourceContents)
	if !ok {
		t.Fatalf("expected TextResourceContents, got %T", contents[0])
	}
	if tc.Text == "" {
		t.Error("resource text is empty")
	}
	if tc.URI != "tmux://workflows/session-debug" {
		t.Errorf("URI = %q, want %q", tc.URI, "tmux://workflows/session-debug")
	}
}

func TestTmuxResourceModule_NilTemplates(t *testing.T) {
	m := &tmuxResourceModule{}
	if m.Templates() != nil {
		t.Error("expected nil templates")
	}
}

// ---------------------------------------------------------------------------
// Prompt module
// ---------------------------------------------------------------------------

func TestTmuxPromptModule_Metadata(t *testing.T) {
	m := &tmuxPromptModule{}
	if m.Name() != "tmux_prompts" {
		t.Errorf("Name() = %q, want %q", m.Name(), "tmux_prompts")
	}
	if m.Description() == "" {
		t.Error("Description() is empty")
	}
}

func TestTmuxPromptModule_Prompts(t *testing.T) {
	m := &tmuxPromptModule{}
	prompts := m.Prompts()
	if len(prompts) != 1 {
		t.Fatalf("expected 1 prompt, got %d", len(prompts))
	}

	pd := prompts[0]
	if pd.Category != "workflow" {
		t.Errorf("Category = %q, want %q", pd.Category, "workflow")
	}
}

func TestTmuxPrompt_Handler(t *testing.T) {
	m := &tmuxPromptModule{}
	pd := m.Prompts()[0]

	req := mcp.GetPromptRequest{}
	req.Params.Arguments = map[string]string{
		"session": "test-session",
		"goal":    "check logs",
	}

	result, err := pd.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if result.Description == "" {
		t.Error("Description is empty")
	}
	if len(result.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(result.Messages))
	}
}

func TestTmuxPrompt_DefaultGoal(t *testing.T) {
	m := &tmuxPromptModule{}
	pd := m.Prompts()[0]

	req := mcp.GetPromptRequest{}
	req.Params.Arguments = map[string]string{
		"session": "test-session",
	}

	result, err := pd.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if result == nil {
		t.Fatal("result is nil")
	}
	// Default goal should produce valid output
	if len(result.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(result.Messages))
	}
}
