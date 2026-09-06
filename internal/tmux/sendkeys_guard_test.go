package tmux

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func fakeQuery(out string, err error) func(args ...string) (string, error) {
	return func(args ...string) (string, error) { return out, err }
}

func TestSendKeysAllowed_RefusesShellPane(t *testing.T) {
	t.Setenv("FLEET_TMUX_SHELL_OK", "")
	ok, why := sendKeysAllowed("cline-fleet:5.1", fakeQuery("cline-fleet:5.1|bash|", nil))
	if ok {
		t.Fatalf("shell pane must be refused, got allowed: %s", why)
	}
	if !strings.Contains(why, "not an agent runtime") {
		t.Fatalf("unexpected reason: %s", why)
	}
}

func TestSendKeysAllowed_AgentPane(t *testing.T) {
	t.Setenv("FLEET_TMUX_SHELL_OK", "")
	for _, cmd := range []string{"node", "claude", "codex", ".cline", "gemini"} {
		ok, why := sendKeysAllowed("s:1.1", fakeQuery("s:1.1|"+cmd+"|", nil))
		if !ok {
			t.Fatalf("%s pane must be allowed: %s", cmd, why)
		}
	}
}

func TestSendKeysAllowed_FailsClosedOnError(t *testing.T) {
	t.Setenv("FLEET_TMUX_SHELL_OK", "")
	if ok, _ := sendKeysAllowed("s:1.1", fakeQuery("", errors.New("no server"))); ok {
		t.Fatal("resolution error must refuse")
	}
	if ok, _ := sendKeysAllowed("s:1.1", fakeQuery("garbage", nil)); ok {
		t.Fatal("unparseable info must refuse")
	}
}

func TestSendKeysAllowed_Overrides(t *testing.T) {
	t.Setenv("FLEET_TMUX_SHELL_OK", "")
	if ok, _ := sendKeysAllowed("s:1.1", fakeQuery("s:1.1|bash|1", nil)); !ok {
		t.Fatal("pane option @accept_shell_injection=1 must allow")
	}
	t.Setenv("FLEET_TMUX_SHELL_OK", "1")
	if ok, _ := sendKeysAllowed("s:1.1", fakeQuery("s:1.1|bash|", nil)); !ok {
		t.Fatal("FLEET_TMUX_SHELL_OK=1 must allow")
	}
}

// Real-tmux mutation check: a bash pane must not receive the keystrokes.
func TestSendKeys_RefusedForRealShellPane(t *testing.T) {
	requireTmux(t)
	t.Setenv("FLEET_TMUX_SHELL_OK", "")
	if _, err := exec.LookPath("bash"); err != nil {
		t.Skip("bash not available")
	}
	sess := "tmuxmcp-guard-" + strings.ReplaceAll(t.Name(), "/", "-")
	_ = exec.Command("tmux", "kill-session", "-t", sess).Run()
	if err := exec.Command("tmux", "new-session", "-d", "-s", sess, "bash").Run(); err != nil {
		t.Skipf("cannot create tmux session: %v", err)
	}
	t.Cleanup(func() { _ = exec.Command("tmux", "kill-session", "-t", sess).Run() })

	td := findTool(t, "tmux_send_keys")
	req := makeReq(map[string]any{"session": sess, "keys": "echo GUARD_INJECTED"})
	result, err := td.Handler(context.Background(), req)
	if err == nil && (result == nil || !result.IsError) {
		t.Fatal("send_keys into a bash pane must be refused")
	}
	out, _ := exec.Command("tmux", "capture-pane", "-p", "-t", sess).Output()
	if strings.Contains(string(out), "GUARD_INJECTED") {
		t.Fatal("keystrokes reached the shell pane despite refusal")
	}
	_ = os.Unsetenv("FLEET_TMUX_SHELL_OK")
}
