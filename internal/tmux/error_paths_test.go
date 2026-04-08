package tmux

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
)

// ---------------------------------------------------------------------------
// Kill nonexistent session
// ---------------------------------------------------------------------------

func TestKillSession_Nonexistent(t *testing.T) {
	requireTmux(t)

	td := findTool(t, "tmux_kill_session")
	req := makeReq(map[string]any{"name": "nonexistent-session-xyz-99999"})
	result, err := td.Handler(context.Background(), req)
	if err != nil {
		// Error is expected
		return
	}
	if result == nil || !result.IsError {
		t.Fatal("expected error for killing nonexistent session")
	}
}

// ---------------------------------------------------------------------------
// New session with duplicate name
// ---------------------------------------------------------------------------

func TestNewSession_Duplicate(t *testing.T) {
	requireTmux(t)

	sessName := fmt.Sprintf("test-dup-%d", rand.Intn(100000))

	// Create first session
	td := findTool(t, "tmux_new_session")
	_, err := td.Handler(context.Background(), makeReq(map[string]any{"name": sessName}))
	if err != nil {
		t.Skipf("could not create session: %v", err)
	}
	defer func() {
		killTd := findTool(t, "tmux_kill_session")
		_, _ = killTd.Handler(context.Background(), makeReq(map[string]any{"name": sessName}))
	}()

	// Try to create a second session with the same name
	result, err := td.Handler(context.Background(), makeReq(map[string]any{"name": sessName}))
	if err != nil {
		// Expected: duplicate session error
		return
	}
	if result == nil || !result.IsError {
		t.Fatal("expected error for duplicate session name")
	}
}

// ---------------------------------------------------------------------------
// Capture pane on nonexistent session
// ---------------------------------------------------------------------------

func TestCapture_NonexistentSession(t *testing.T) {
	requireTmux(t)

	td := findTool(t, "tmux_capture_pane")
	req := makeReq(map[string]any{"session": "nonexistent-xyz-12345"})
	result, err := td.Handler(context.Background(), req)
	if err != nil {
		return // error expected
	}
	if result == nil || !result.IsError {
		t.Fatal("expected error for capture on nonexistent session")
	}
}

// ---------------------------------------------------------------------------
// Send keys to nonexistent session
// ---------------------------------------------------------------------------

func TestSendKeys_NonexistentSession(t *testing.T) {
	requireTmux(t)

	td := findTool(t, "tmux_send_keys")
	req := makeReq(map[string]any{
		"session": "nonexistent-xyz-12345",
		"keys":    "echo test",
	})
	result, err := td.Handler(context.Background(), req)
	if err != nil {
		return // error expected
	}
	if result == nil || !result.IsError {
		t.Fatal("expected error for send_keys on nonexistent session")
	}
}

// ---------------------------------------------------------------------------
// List panes on nonexistent session
// ---------------------------------------------------------------------------

func TestListPanes_NonexistentSession(t *testing.T) {
	requireTmux(t)

	td := findTool(t, "tmux_list_panes")
	req := makeReq(map[string]any{"session": "nonexistent-xyz-12345"})
	result, err := td.Handler(context.Background(), req)
	if err != nil {
		return // error expected
	}
	// Could be error or empty list depending on tmux version
	if result != nil && !result.IsError {
		var out ListPanesOutput
		unmarshalResult(t, result, &out)
		// Should be empty
	}
}

// ---------------------------------------------------------------------------
// New window on nonexistent session
// ---------------------------------------------------------------------------

func TestNewWindow_NonexistentSession(t *testing.T) {
	requireTmux(t)

	td := findTool(t, "tmux_new_window")
	req := makeReq(map[string]any{"session": "nonexistent-xyz-12345"})
	result, err := td.Handler(context.Background(), req)
	if err != nil {
		return // error expected
	}
	if result == nil || !result.IsError {
		t.Fatal("expected error for new_window on nonexistent session")
	}
}

// ---------------------------------------------------------------------------
// Wait for text with context cancellation
// ---------------------------------------------------------------------------

func TestWaitForText_ShortTimeout(t *testing.T) {
	requireTmux(t)

	sessName := fmt.Sprintf("test-wft-short-%d", rand.Intn(100000))
	newTd := findTool(t, "tmux_new_session")
	_, err := newTd.Handler(context.Background(), makeReq(map[string]any{"name": sessName}))
	if err != nil {
		t.Skipf("could not create session: %v", err)
	}
	defer func() {
		killTd := findTool(t, "tmux_kill_session")
		_, _ = killTd.Handler(context.Background(), makeReq(map[string]any{"name": sessName}))
	}()

	td := findTool(t, "tmux_wait_for_text")
	req := makeReq(map[string]any{
		"session":      sessName,
		"text":         "NEVER_APPEARING_TEXT_XYZ",
		"timeout_secs": 1,
	})
	result, err := td.Handler(context.Background(), req)
	// Should timeout
	if err == nil && (result == nil || !result.IsError) {
		t.Fatal("expected timeout error for text that never appears")
	}
}

// ---------------------------------------------------------------------------
// Search panes on nonexistent session
// ---------------------------------------------------------------------------

func TestSearchPanes_NonexistentSession(t *testing.T) {
	requireTmux(t)

	td := findTool(t, "tmux_search_panes")
	req := makeReq(map[string]any{
		"session": "nonexistent-xyz-12345",
		"pattern": "test",
	})
	result, err := td.Handler(context.Background(), req)
	if err != nil {
		return // error expected
	}
	// Graceful empty result or error result
	if result != nil && !result.IsError {
		var out SearchPanesOutput
		unmarshalResult(t, result, &out)
		if len(out.Matches) != 0 {
			t.Error("expected no matches for nonexistent session")
		}
	}
}
