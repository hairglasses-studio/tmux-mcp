package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os/exec"
	"testing"

	"github.com/hairglasses-studio/mcpkit/mcptest"
	"github.com/hairglasses-studio/mcpkit/registry"
)

// ---------------------------------------------------------------------------
// Registration
// ---------------------------------------------------------------------------

func TestModuleRegistration(t *testing.T) {
	m := &TmuxModule{}
	tools := m.Tools()
	if len(tools) != 8 {
		t.Fatalf("expected 8 tools, got %d", len(tools))
	}

	reg := registry.NewToolRegistry()
	reg.RegisterModule(m)
	srv := mcptest.NewServer(t, reg)

	names := srv.ToolNames()
	if len(names) != 8 {
		t.Fatalf("expected 8 registered tools, got %d", len(names))
	}

	for _, want := range []string{
		"tmux_list_sessions", "tmux_new_session", "tmux_kill_session",
		"tmux_list_windows", "tmux_new_window",
		"tmux_list_panes", "tmux_capture_pane", "tmux_send_keys",
	} {
		if !srv.HasTool(want) {
			t.Errorf("missing tool: %s", want)
		}
	}
}

// ---------------------------------------------------------------------------
// Helper unit tests
// ---------------------------------------------------------------------------

func TestIsNoServer(t *testing.T) {
	tests := []struct {
		err  error
		want bool
	}{
		{nil, false},
		{fmt.Errorf("no server running on /tmp/tmux"), true},
		{fmt.Errorf("no current client"), true},
		{fmt.Errorf("error connecting to /tmp/tmux"), true},
		{fmt.Errorf("session not found"), false},
		{fmt.Errorf("some other error"), false},
	}

	for _, tc := range tests {
		got := isNoServer(tc.err)
		if got != tc.want {
			t.Errorf("isNoServer(%v) = %v, want %v", tc.err, got, tc.want)
		}
	}
}

func TestTmuxTarget(t *testing.T) {
	tests := []struct {
		session, window, pane string
		want                  string
	}{
		{"sess", "", "", "sess"},
		{"sess", "1", "", "sess:1"},
		{"sess", "1", "0", "sess:1.0"},
		{"sess", "", "0", "sess.0"},
		{"my-session", "main", "2", "my-session:main.2"},
	}

	for _, tc := range tests {
		got := tmuxTarget(tc.session, tc.window, tc.pane)
		if got != tc.want {
			t.Errorf("tmuxTarget(%q, %q, %q) = %q, want %q",
				tc.session, tc.window, tc.pane, got, tc.want)
		}
	}
}

// ---------------------------------------------------------------------------
// List sessions (graceful when no server)
// ---------------------------------------------------------------------------

func TestListSessions_GracefulNoServer(t *testing.T) {
	requireTmux(t)

	// This test works whether or not tmux is running — it should not error
	td := findTool(t, "tmux_list_sessions")
	req := makeReq(nil)
	result, err := td.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("expected graceful handling, got error: %v", err)
	}

	var out ListSessionsOutput
	unmarshalResult(t, result, &out)
	if out.Sessions == nil {
		t.Error("sessions should be non-nil (empty slice)")
	}
}

// ---------------------------------------------------------------------------
// Session lifecycle (create → list → windows → panes → capture → send → kill)
// ---------------------------------------------------------------------------

func TestSessionLifecycle(t *testing.T) {
	requireTmux(t)

	sessName := fmt.Sprintf("test-mcp-%d", rand.Intn(100000))

	// Create session
	newTd := findTool(t, "tmux_new_session")
	req := makeReq(map[string]any{"name": sessName})
	result, err := newTd.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("new_session error: %v", err)
	}
	var newOut NewSessionOutput
	unmarshalResult(t, result, &newOut)
	if !newOut.Created {
		t.Error("expected created=true")
	}

	// Cleanup: kill session at end
	defer func() {
		killTd := findTool(t, "tmux_kill_session")
		killReq := makeReq(map[string]any{"name": sessName})
		killTd.Handler(context.Background(), killReq)
	}()

	// List sessions — should include our session
	listTd := findTool(t, "tmux_list_sessions")
	result, err = listTd.Handler(context.Background(), makeReq(nil))
	if err != nil {
		t.Fatalf("list_sessions error: %v", err)
	}
	var listOut ListSessionsOutput
	unmarshalResult(t, result, &listOut)
	found := false
	for _, s := range listOut.Sessions {
		if s.Name == sessName {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("session %q not found in list", sessName)
	}

	// List windows
	winTd := findTool(t, "tmux_list_windows")
	result, err = winTd.Handler(context.Background(), makeReq(map[string]any{"session": sessName}))
	if err != nil {
		t.Fatalf("list_windows error: %v", err)
	}
	var winOut ListWindowsOutput
	unmarshalResult(t, result, &winOut)
	if len(winOut.Windows) == 0 {
		t.Error("expected at least 1 window")
	}

	// List panes
	paneTd := findTool(t, "tmux_list_panes")
	result, err = paneTd.Handler(context.Background(), makeReq(map[string]any{"session": sessName}))
	if err != nil {
		t.Fatalf("list_panes error: %v", err)
	}
	var paneOut ListPanesOutput
	unmarshalResult(t, result, &paneOut)
	if len(paneOut.Panes) == 0 {
		t.Error("expected at least 1 pane")
	}

	// Send keys
	sendTd := findTool(t, "tmux_send_keys")
	result, err = sendTd.Handler(context.Background(), makeReq(map[string]any{
		"session": sessName,
		"keys":    "echo hello-mcp-test",
	}))
	if err != nil {
		t.Fatalf("send_keys error: %v", err)
	}
	var sendOut SendKeysOutput
	unmarshalResult(t, result, &sendOut)
	if sendOut.Sent != "echo hello-mcp-test" {
		t.Errorf("expected sent=%q, got %q", "echo hello-mcp-test", sendOut.Sent)
	}

	// Capture pane
	capTd := findTool(t, "tmux_capture_pane")
	result, err = capTd.Handler(context.Background(), makeReq(map[string]any{
		"session": sessName,
		"lines":   10,
	}))
	if err != nil {
		t.Fatalf("capture_pane error: %v", err)
	}
	var capOut CaptureOutput
	unmarshalResult(t, result, &capOut)
	// Content may or may not contain our text depending on timing

	// New window
	newWinTd := findTool(t, "tmux_new_window")
	result, err = newWinTd.Handler(context.Background(), makeReq(map[string]any{
		"session": sessName,
		"name":    "test-win",
	}))
	if err != nil {
		t.Fatalf("new_window error: %v", err)
	}
	var newWinOut NewWindowOutput
	unmarshalResult(t, result, &newWinOut)
	if newWinOut.Session != sessName {
		t.Errorf("expected session=%q, got %q", sessName, newWinOut.Session)
	}

	// Kill session (explicit)
	killTd := findTool(t, "tmux_kill_session")
	result, err = killTd.Handler(context.Background(), makeReq(map[string]any{"name": sessName}))
	if err != nil {
		t.Fatalf("kill_session error: %v", err)
	}
	var killOut KillSessionOutput
	unmarshalResult(t, result, &killOut)
	if !killOut.Killed {
		t.Error("expected killed=true")
	}
}

// ---------------------------------------------------------------------------
// Error cases
// ---------------------------------------------------------------------------

func TestListWindows_BadSession(t *testing.T) {
	requireTmux(t)

	td := findTool(t, "tmux_list_windows")
	req := makeReq(map[string]any{"session": "nonexistent-session-xyz-99999"})
	result, err := td.Handler(context.Background(), req)
	if err != nil {
		// Error is acceptable — session doesn't exist
		return
	}
	// Typed handler wraps errors as IsError results
	if result != nil && result.IsError {
		return
	}
	// If no error, should return empty windows (graceful)
	var out ListWindowsOutput
	unmarshalResult(t, result, &out)
	if len(out.Windows) != 0 {
		t.Errorf("expected 0 windows for nonexistent session, got %d", len(out.Windows))
	}
}

func TestCapturePane_DefaultLines(t *testing.T) {
	// Verify the handler defaults lines=50 when 0 is provided
	// We can't easily test this without a running session, so just verify
	// the code path doesn't panic with lines=0
	requireTmux(t)

	sessName := fmt.Sprintf("test-cap-%d", rand.Intn(100000))
	newTd := findTool(t, "tmux_new_session")
	_, err := newTd.Handler(context.Background(), makeReq(map[string]any{"name": sessName}))
	if err != nil {
		t.Skipf("could not create test session: %v", err)
	}
	defer func() {
		killTd := findTool(t, "tmux_kill_session")
		killTd.Handler(context.Background(), makeReq(map[string]any{"name": sessName}))
	}()

	td := findTool(t, "tmux_capture_pane")
	req := makeReq(map[string]any{"session": sessName, "lines": 0})
	result, err := td.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("capture with lines=0 error: %v", err)
	}
	var out CaptureOutput
	unmarshalResult(t, result, &out)
	// Should have captured something (default 50 lines)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func requireTmux(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("tmux"); err != nil {
		t.Skip("tmux not available")
	}
}

func findTool(t *testing.T, name string) registry.ToolDefinition {
	t.Helper()
	m := &TmuxModule{}
	for _, td := range m.Tools() {
		if td.Tool.Name == name {
			return td
		}
	}
	t.Fatalf("tool %q not found", name)
	return registry.ToolDefinition{}
}

func makeReq(args map[string]any) registry.CallToolRequest {
	req := registry.CallToolRequest{}
	if args == nil {
		args = map[string]any{}
	}
	req.Params.Arguments = args
	return req
}

func extractText(t *testing.T, result *registry.CallToolResult) string {
	t.Helper()
	if len(result.Content) == 0 {
		t.Fatal("result has no content")
	}
	tc, ok := result.Content[0].(registry.TextContent)
	if !ok {
		t.Fatalf("content is not TextContent, got %T", result.Content[0])
	}
	return tc.Text
}

func unmarshalResult(t *testing.T, result *registry.CallToolResult, out any) {
	t.Helper()
	text := extractText(t, result)
	if err := json.Unmarshal([]byte(text), out); err != nil {
		t.Fatalf("unmarshal error: %v; text=%s", err, text)
	}
}
