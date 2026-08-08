package tmux

import (
	"context"
	"testing"

	"github.com/hairglasses-studio/mcpkit/resources"
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
	defs := m.Resources()
	if len(defs) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(defs))
	}

	rd := defs[0]
	if rd.Category != "workflow" {
		t.Errorf("Category = %q, want %q", rd.Category, "workflow")
	}
	if len(rd.Tags) == 0 {
		t.Error("expected tags")
	}
	if rd.Resource.URI != "tmux://workflows/session-debug" {
		t.Errorf("URI = %q, want %q", rd.Resource.URI, "tmux://workflows/session-debug")
	}

	mimeType, text, err := resources.CallHandlerText(context.Background(), rd.Handler, "tmux://workflows/session-debug")
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if mimeType != "text/markdown" {
		t.Errorf("mimeType = %q, want %q", mimeType, "text/markdown")
	}
	if text == "" {
		t.Error("resource text is empty")
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
	defs := m.Prompts()
	if len(defs) != 1 {
		t.Fatalf("expected 1 prompt, got %d", len(defs))
	}

	pd := defs[0]
	if pd.Category != "workflow" {
		t.Errorf("Category = %q, want %q", pd.Category, "workflow")
	}
}

func TestTmuxPrompt_Handler(t *testing.T) {
	m := &tmuxPromptModule{}
	pd := m.Prompts()[0]

	result, err := callPrompt(pd.Handler, map[string]string{
		"session": "test-session",
		"goal":    "check logs",
	})
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

	result, err := callPrompt(pd.Handler, map[string]string{
		"session": "test-session",
	})
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if result == nil {
		t.Fatal("result is nil")
	}
	if len(result.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(result.Messages))
	}
}
