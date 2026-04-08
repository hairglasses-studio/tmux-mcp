package tmux

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
)

// ---------------------------------------------------------------------------
// tmux_capture_pane edge cases
// ---------------------------------------------------------------------------

func TestCapture_LineCounts(t *testing.T) {
	requireTmux(t)

	sessName := fmt.Sprintf("test-edge-cap-%d", rand.Intn(100000))
	newTd := findTool(t, "tmux_new_session")
	_, err := newTd.Handler(context.Background(), makeReq(map[string]any{"name": sessName}))
	if err != nil {
		t.Skipf("could not create session: %v", err)
	}
	defer func() {
		killTd := findTool(t, "tmux_kill_session")
		_, _ = killTd.Handler(context.Background(), makeReq(map[string]any{"name": sessName}))
	}()

	td := findTool(t, "tmux_capture_pane")

	tests := []struct {
		name  string
		lines int
	}{
		{"zero_defaults_to_50", 0},
		{"one_line", 1},
		{"large_count", 5000},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := makeReq(map[string]any{"session": sessName, "lines": tc.lines})
			result, err := td.Handler(context.Background(), req)
			if err != nil {
				t.Fatalf("capture error: %v", err)
			}
			var out CaptureOutput
			unmarshalResult(t, result, &out)
			// Should not panic and should return valid output
			if out.Target == "" {
				t.Error("target is empty")
			}
		})
	}
}

// ---------------------------------------------------------------------------
// tmux_send_keys with special characters
// ---------------------------------------------------------------------------

func TestSendKeys_SpecialChars(t *testing.T) {
	requireTmux(t)

	sessName := fmt.Sprintf("test-edge-keys-%d", rand.Intn(100000))
	newTd := findTool(t, "tmux_new_session")
	_, err := newTd.Handler(context.Background(), makeReq(map[string]any{"name": sessName}))
	if err != nil {
		t.Skipf("could not create session: %v", err)
	}
	defer func() {
		killTd := findTool(t, "tmux_kill_session")
		_, _ = killTd.Handler(context.Background(), makeReq(map[string]any{"name": sessName}))
	}()

	td := findTool(t, "tmux_send_keys")

	tests := []struct {
		name string
		keys string
	}{
		{"enter", "Enter"},
		{"escape", "Escape"},
		{"ctrl_c", "C-c"},
		{"tab", "Tab"},
		{"literal_text", "echo test"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := makeReq(map[string]any{"session": sessName, "keys": tc.keys})
			result, err := td.Handler(context.Background(), req)
			if err != nil {
				t.Fatalf("send_keys error: %v", err)
			}
			var out SendKeysOutput
			unmarshalResult(t, result, &out)
			if out.Sent != tc.keys {
				t.Errorf("sent = %q, want %q", out.Sent, tc.keys)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// tmux_workspace with various layouts
// ---------------------------------------------------------------------------

func TestWorkspace_ValidLayouts(t *testing.T) {
	requireTmux(t)

	layouts := []string{
		"tiled",
		"even-horizontal",
		"even-vertical",
		"main-horizontal",
		"main-vertical",
	}

	for _, layout := range layouts {
		t.Run(layout, func(t *testing.T) {
			sessName := fmt.Sprintf("test-layout-%d", rand.Intn(100000))
			td := findTool(t, "tmux_workspace")

			req := makeReq(map[string]any{
				"session": sessName,
				"windows": []any{
					map[string]any{
						"name":   "win",
						"layout": layout,
						"panes": []any{
							map[string]any{"command": "true"},
							map[string]any{"command": "true"},
						},
					},
				},
			})

			result, err := td.Handler(context.Background(), req)
			if err != nil {
				t.Fatalf("workspace error: %v", err)
			}

			defer func() {
				killTd := findTool(t, "tmux_kill_session")
				_, _ = killTd.Handler(context.Background(), makeReq(map[string]any{"name": sessName}))
			}()

			var out WorkspaceOutput
			unmarshalResult(t, result, &out)
			if out.PaneCount != 2 {
				t.Errorf("expected 2 panes, got %d", out.PaneCount)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// tmux_new_session with directory and command
// ---------------------------------------------------------------------------

func TestNewSession_WithDirAndCommand(t *testing.T) {
	requireTmux(t)

	sessName := fmt.Sprintf("test-dir-cmd-%d", rand.Intn(100000))
	td := findTool(t, "tmux_new_session")
	req := makeReq(map[string]any{
		"name":      sessName,
		"directory": "/tmp",
		"command":   "echo hello",
	})
	result, err := td.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("new_session error: %v", err)
	}
	defer func() {
		killTd := findTool(t, "tmux_kill_session")
		_, _ = killTd.Handler(context.Background(), makeReq(map[string]any{"name": sessName}))
	}()

	var out NewSessionOutput
	unmarshalResult(t, result, &out)
	if !out.Created {
		t.Error("expected created=true")
	}
	if out.Session != sessName {
		t.Errorf("session = %q, want %q", out.Session, sessName)
	}
}

// ---------------------------------------------------------------------------
// tmux_new_window with name and command
// ---------------------------------------------------------------------------

func TestNewWindow_WithNameAndCommand(t *testing.T) {
	requireTmux(t)

	sessName := fmt.Sprintf("test-nw-%d", rand.Intn(100000))
	newTd := findTool(t, "tmux_new_session")
	_, err := newTd.Handler(context.Background(), makeReq(map[string]any{"name": sessName}))
	if err != nil {
		t.Skipf("could not create session: %v", err)
	}
	defer func() {
		killTd := findTool(t, "tmux_kill_session")
		_, _ = killTd.Handler(context.Background(), makeReq(map[string]any{"name": sessName}))
	}()

	td := findTool(t, "tmux_new_window")
	req := makeReq(map[string]any{
		"session": sessName,
		"name":    "my-win",
		"command": "echo window-test",
	})
	result, err := td.Handler(context.Background(), req)
	if err != nil {
		t.Fatalf("new_window error: %v", err)
	}
	var out NewWindowOutput
	unmarshalResult(t, result, &out)
	if out.Session != sessName {
		t.Errorf("session = %q, want %q", out.Session, sessName)
	}
}

