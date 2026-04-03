// Command tmux-mcp is an MCP server for tmux session, window, and pane
// management via the Model Context Protocol (stdio transport).
//
// Usage:
//
//	tmux-mcp
package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/hairglasses-studio/mcpkit/handler"
	"github.com/hairglasses-studio/mcpkit/registry"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func runCmd(name string, args ...string) (string, string, error) {
	cmd := exec.Command(name, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), err
}

func runTmux(args ...string) (string, error) {
	stdout, stderr, err := runCmd("tmux", args...)
	if err != nil {
		return "", fmt.Errorf("tmux %s: %s: %w", strings.Join(args, " "), stderr, err)
	}
	return stdout, nil
}

// isNoServer returns true if the error indicates tmux has no running server.
func isNoServer(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "no server running") ||
		strings.Contains(msg, "no current client") ||
		strings.Contains(msg, "error connecting")
}

// tmuxTarget builds a tmux target string like "session:window.pane".
func tmuxTarget(session, window, pane string) string {
	t := session
	if window != "" {
		t += ":" + window
	}
	if pane != "" {
		t += "." + pane
	}
	return t
}

// ---------------------------------------------------------------------------
// I/O types
// ---------------------------------------------------------------------------

// --- list_sessions ---

type ListSessionsInput struct{}

type sessionInfo struct {
	Name     string `json:"name"`
	Windows  int    `json:"windows"`
	Created  string `json:"created"`
	Attached bool   `json:"attached"`
}

type ListSessionsOutput struct {
	Sessions []sessionInfo `json:"sessions"`
}

// --- list_windows ---

type ListWindowsInput struct {
	Session string `json:"session" jsonschema:"required,description=Tmux session name"`
}

type windowInfo struct {
	Index  int    `json:"index"`
	Name   string `json:"name"`
	Panes  int    `json:"panes"`
	Active bool   `json:"active"`
	Layout string `json:"layout"`
}

type ListWindowsOutput struct {
	Session string       `json:"session"`
	Windows []windowInfo `json:"windows"`
}

// --- list_panes ---

type ListPanesInput struct {
	Session string `json:"session" jsonschema:"required,description=Tmux session name"`
	Window  string `json:"window,omitempty" jsonschema:"description=Window index or name"`
}

type paneInfo struct {
	Index   int    `json:"index"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Command string `json:"command"`
	Active  bool   `json:"active"`
}

type ListPanesOutput struct {
	Target string     `json:"target"`
	Panes  []paneInfo `json:"panes"`
}

// --- capture_pane ---

type CaptureInput struct {
	Session string `json:"session" jsonschema:"required,description=Tmux session name"`
	Window  string `json:"window,omitempty" jsonschema:"description=Window index or name"`
	Pane    string `json:"pane,omitempty" jsonschema:"description=Pane index"`
	Lines   int    `json:"lines,omitempty" jsonschema:"description=Number of lines to capture. Default 50."`
}

type CaptureOutput struct {
	Target  string `json:"target"`
	Content string `json:"content"`
}

// --- send_keys ---

type SendKeysInput struct {
	Session string `json:"session" jsonschema:"required,description=Tmux session name"`
	Window  string `json:"window,omitempty" jsonschema:"description=Window index or name"`
	Pane    string `json:"pane,omitempty" jsonschema:"description=Pane index"`
	Keys    string `json:"keys" jsonschema:"required,description=Keys to send. Literal text or tmux key names like Enter or C-c."`
}

type SendKeysOutput struct {
	Target string `json:"target"`
	Sent   string `json:"sent"`
}

// --- new_session ---

type NewSessionInput struct {
	Name      string `json:"name" jsonschema:"required,description=Session name"`
	Command   string `json:"command,omitempty" jsonschema:"description=Initial command to run"`
	Directory string `json:"directory,omitempty" jsonschema:"description=Working directory for the session"`
}

type NewSessionOutput struct {
	Session string `json:"session"`
	Created bool   `json:"created"`
}

// --- new_window ---

type NewWindowInput struct {
	Session string `json:"session" jsonschema:"required,description=Tmux session name"`
	Name    string `json:"name,omitempty" jsonschema:"description=Window name"`
	Command string `json:"command,omitempty" jsonschema:"description=Command to run in the new window"`
}

type NewWindowOutput struct {
	Session string `json:"session"`
	Window  string `json:"window"`
}

// --- kill_session ---

type KillSessionInput struct {
	Name string `json:"name" jsonschema:"required,description=Session name to kill"`
}

type KillSessionOutput struct {
	Session string `json:"session"`
	Killed  bool   `json:"killed"`
}

// ---------------------------------------------------------------------------
// TmuxModule
// ---------------------------------------------------------------------------

type TmuxModule struct{}

func (m *TmuxModule) Name() string        { return "tmux" }
func (m *TmuxModule) Description() string { return "Tmux session, window, and pane management" }

func (m *TmuxModule) Tools() []registry.ToolDefinition {
	return []registry.ToolDefinition{
		// ---------------------------------------------------------------
		// tmux_list_sessions
		// ---------------------------------------------------------------
		handler.TypedHandler[ListSessionsInput, ListSessionsOutput](
			"tmux_list_sessions",
			"List all tmux sessions with window count, creation time, and attached status.",
			func(_ context.Context, _ ListSessionsInput) (ListSessionsOutput, error) {
				out, err := runTmux("list-sessions", "-F",
					"#{session_name}\t#{session_windows}\t#{session_created}\t#{session_attached}")
				if err != nil {
					if isNoServer(err) {
						return ListSessionsOutput{Sessions: []sessionInfo{}}, nil
					}
					return ListSessionsOutput{}, fmt.Errorf("[%s] list sessions: %w", handler.ErrAPIError, err)
				}

				var sessions []sessionInfo
				for _, line := range strings.Split(out, "\n") {
					if line == "" {
						continue
					}
					parts := strings.SplitN(line, "\t", 4)
					if len(parts) < 4 {
						continue
					}
					windows, _ := strconv.Atoi(parts[1])
					attached := parts[3] == "1"
					sessions = append(sessions, sessionInfo{
						Name:     parts[0],
						Windows:  windows,
						Created:  parts[2],
						Attached: attached,
					})
				}
				if sessions == nil {
					sessions = []sessionInfo{}
				}
				return ListSessionsOutput{Sessions: sessions}, nil
			},
		),

		// ---------------------------------------------------------------
		// tmux_list_windows
		// ---------------------------------------------------------------
		handler.TypedHandler[ListWindowsInput, ListWindowsOutput](
			"tmux_list_windows",
			"List windows in a tmux session with index, name, pane count, active status, and layout.",
			func(_ context.Context, input ListWindowsInput) (ListWindowsOutput, error) {
				out, err := runTmux("list-windows", "-t", input.Session, "-F",
					"#{window_index}\t#{window_name}\t#{window_panes}\t#{window_active}\t#{window_layout}")
				if err != nil {
					if isNoServer(err) {
						return ListWindowsOutput{Session: input.Session, Windows: []windowInfo{}}, nil
					}
					return ListWindowsOutput{}, fmt.Errorf("[%s] list windows: %w", handler.ErrNotFound, err)
				}

				var windows []windowInfo
				for _, line := range strings.Split(out, "\n") {
					if line == "" {
						continue
					}
					parts := strings.SplitN(line, "\t", 5)
					if len(parts) < 5 {
						continue
					}
					idx, _ := strconv.Atoi(parts[0])
					panes, _ := strconv.Atoi(parts[2])
					active := parts[3] == "1"
					windows = append(windows, windowInfo{
						Index:  idx,
						Name:   parts[1],
						Panes:  panes,
						Active: active,
						Layout: parts[4],
					})
				}
				if windows == nil {
					windows = []windowInfo{}
				}
				return ListWindowsOutput{
					Session: input.Session,
					Windows: windows,
				}, nil
			},
		),

		// ---------------------------------------------------------------
		// tmux_list_panes
		// ---------------------------------------------------------------
		handler.TypedHandler[ListPanesInput, ListPanesOutput](
			"tmux_list_panes",
			"List panes in a session or window with index, dimensions, current command, and active status.",
			func(_ context.Context, input ListPanesInput) (ListPanesOutput, error) {
				target := tmuxTarget(input.Session, input.Window, "")
				out, err := runTmux("list-panes", "-t", target, "-F",
					"#{pane_index}\t#{pane_width}\t#{pane_height}\t#{pane_current_command}\t#{pane_active}")
				if err != nil {
					if isNoServer(err) {
						return ListPanesOutput{Target: target, Panes: []paneInfo{}}, nil
					}
					return ListPanesOutput{}, fmt.Errorf("[%s] list panes: %w", handler.ErrNotFound, err)
				}

				var panes []paneInfo
				for _, line := range strings.Split(out, "\n") {
					if line == "" {
						continue
					}
					parts := strings.SplitN(line, "\t", 5)
					if len(parts) < 5 {
						continue
					}
					idx, _ := strconv.Atoi(parts[0])
					width, _ := strconv.Atoi(parts[1])
					height, _ := strconv.Atoi(parts[2])
					active := parts[4] == "1"
					panes = append(panes, paneInfo{
						Index:   idx,
						Width:   width,
						Height:  height,
						Command: parts[3],
						Active:  active,
					})
				}
				if panes == nil {
					panes = []paneInfo{}
				}
				return ListPanesOutput{
					Target: target,
					Panes:  panes,
				}, nil
			},
		),

		// ---------------------------------------------------------------
		// tmux_capture_pane
		// ---------------------------------------------------------------
		handler.TypedHandler[CaptureInput, CaptureOutput](
			"tmux_capture_pane",
			"Capture visible text content from a tmux pane. Default 50 lines of scrollback.",
			func(_ context.Context, input CaptureInput) (CaptureOutput, error) {
				target := tmuxTarget(input.Session, input.Window, input.Pane)
				lines := input.Lines
				if lines <= 0 {
					lines = 50
				}
				out, err := runTmux("capture-pane", "-p", "-t", target,
					"-S", fmt.Sprintf("-%d", lines))
				if err != nil {
					return CaptureOutput{}, fmt.Errorf("[%s] capture pane: %w", handler.ErrNotFound, err)
				}
				return CaptureOutput{
					Target:  target,
					Content: out,
				}, nil
			},
		),

		// ---------------------------------------------------------------
		// tmux_send_keys
		// ---------------------------------------------------------------
		handler.TypedHandler[SendKeysInput, SendKeysOutput](
			"tmux_send_keys",
			"Send keystrokes to a tmux pane. Keys can be literal text or tmux key names like Enter, C-c, etc.",
			func(_ context.Context, input SendKeysInput) (SendKeysOutput, error) {
				target := tmuxTarget(input.Session, input.Window, input.Pane)
				_, err := runTmux("send-keys", "-t", target, input.Keys)
				if err != nil {
					return SendKeysOutput{}, fmt.Errorf("[%s] send keys: %w", handler.ErrNotFound, err)
				}
				return SendKeysOutput{
					Target: target,
					Sent:   input.Keys,
				}, nil
			},
		),

		// ---------------------------------------------------------------
		// tmux_new_session
		// ---------------------------------------------------------------
		handler.TypedHandler[NewSessionInput, NewSessionOutput](
			"tmux_new_session",
			"Create a new detached tmux session with optional working directory and initial command.",
			func(_ context.Context, input NewSessionInput) (NewSessionOutput, error) {
				args := []string{"new-session", "-d", "-s", input.Name}
				if input.Directory != "" {
					args = append(args, "-c", input.Directory)
				}
				if input.Command != "" {
					args = append(args, input.Command)
				}
				_, err := runTmux(args...)
				if err != nil {
					return NewSessionOutput{}, fmt.Errorf("[%s] new session: %w", handler.ErrAPIError, err)
				}
				return NewSessionOutput{
					Session: input.Name,
					Created: true,
				}, nil
			},
		),

		// ---------------------------------------------------------------
		// tmux_new_window
		// ---------------------------------------------------------------
		handler.TypedHandler[NewWindowInput, NewWindowOutput](
			"tmux_new_window",
			"Create a new window in an existing tmux session with optional name and command.",
			func(_ context.Context, input NewWindowInput) (NewWindowOutput, error) {
				args := []string{"new-window", "-t", input.Session, "-P", "-F", "#{window_index}"}
				if input.Name != "" {
					args = append(args, "-n", input.Name)
				}
				if input.Command != "" {
					args = append(args, input.Command)
				}
				out, err := runTmux(args...)
				if err != nil {
					return NewWindowOutput{}, fmt.Errorf("[%s] new window: %w", handler.ErrNotFound, err)
				}
				return NewWindowOutput{
					Session: input.Session,
					Window:  out,
				}, nil
			},
		),

		// ---------------------------------------------------------------
		// tmux_kill_session
		// ---------------------------------------------------------------
		handler.TypedHandler[KillSessionInput, KillSessionOutput](
			"tmux_kill_session",
			"Kill a tmux session by name.",
			func(_ context.Context, input KillSessionInput) (KillSessionOutput, error) {
				_, err := runTmux("kill-session", "-t", input.Name)
				if err != nil {
					return KillSessionOutput{}, fmt.Errorf("[%s] kill session: %w", handler.ErrNotFound, err)
				}
				return KillSessionOutput{
					Session: input.Name,
					Killed:  true,
				}, nil
			},
		),
	}
}

// ---------------------------------------------------------------------------
// main
// ---------------------------------------------------------------------------

func main() {
	reg := registry.NewToolRegistry()
	reg.RegisterModule(&TmuxModule{})

	s := registry.NewMCPServer("tmux-mcp", "1.0.0")
	reg.RegisterWithServer(s)

	if err := registry.ServeStdio(s); err != nil {
		log.Fatal(err)
	}
}
