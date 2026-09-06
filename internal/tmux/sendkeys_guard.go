package tmux

import (
	"fmt"
	"os"
	"strings"
)

// sendKeysAgentRuntimes lists pane foreground commands that are agent
// runtimes and therefore legitimate keystroke targets. A shell prompt is the
// UNSAFE state: on 2026-09-06 four fleet loops typed directives into twelve
// bash panes every 25-30 s (the cline-fleet echo storm).
var sendKeysAgentRuntimes = map[string]bool{
	"node": true, ".cline": true, "cline": true, "claude": true, "claude-real": true,
	"codex": true, "gemini": true, "agy": true,
}

const (
	sendKeysShellOKEnv     = "FLEET_TMUX_SHELL_OK"
	sendKeysAcceptPaneOpt  = "@accept_shell_injection"
	sendKeysGuardFormatStr = "#{session_name}:#{window_index}.#{pane_index}|#{pane_current_command}|#{@accept_shell_injection}"
)

// sendKeysAllowed decides whether tmux_send_keys may type into target.
// It fails closed: any resolution error refuses. paneInfo is the raw
// display-message output so callers can report what was resolved.
func sendKeysAllowed(target string, query func(args ...string) (string, error)) (allowed bool, reason string) {
	if os.Getenv(sendKeysShellOKEnv) == "1" {
		return true, "FLEET_TMUX_SHELL_OK=1 override"
	}
	out, err := query("display-message", "-p", "-t", target, sendKeysGuardFormatStr)
	if err != nil {
		return false, fmt.Sprintf("refused: cannot resolve pane %q: %v", target, err)
	}
	parts := strings.Split(strings.TrimSpace(out), "|")
	if len(parts) < 2 {
		return false, fmt.Sprintf("refused: unparseable pane info for %q", target)
	}
	resolved, cmd := parts[0], parts[1]
	accept := len(parts) > 2 && parts[2] == "1"
	if sendKeysAgentRuntimes[cmd] {
		return true, fmt.Sprintf("agent runtime %s in %s", cmd, resolved)
	}
	if accept {
		return true, fmt.Sprintf("pane %s opted in via %s", resolved, sendKeysAcceptPaneOpt)
	}
	return false, fmt.Sprintf("refused: pane %s runs %q, not an agent runtime (2026-09-06 echo-storm guard); use tmux display-message for notices or set %s=1 / pane option %s=1 for a pane you own", resolved, cmd, sendKeysShellOKEnv, sendKeysAcceptPaneOpt)
}
