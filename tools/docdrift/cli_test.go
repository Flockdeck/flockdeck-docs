package main

import "testing"

func realCommands(t *testing.T) map[string]bool {
	t.Helper()
	cmds, err := parseDesktopHelp(fixture(t, "desktop-v0.3.40-help.txt"))
	if err != nil {
		t.Fatal(err)
	}
	return cmds
}

func TestParseDesktopHelp(t *testing.T) {
	cmds := realCommands(t)
	for _, want := range []string{
		"spawn", "peer-name", "close", "agents", "chat", "keys", "remote", "update", "help",
		"remote/enable", "remote/pair", "remote/status", "remote/devices", "remote/revoke",
		"remote/disable", "remote/rename", "remote/move",
	} {
		if !cmds[want] {
			t.Errorf("%s not parsed; got %v", want, sortedKeys(cmds))
		}
	}
	// "close [-force] <pane-id> | close -finished" is one command, not "close/close".
	for _, not := range []string{"close/close", "keys/list", "update/-check", "spawn/-worktree"} {
		if cmds[not] {
			t.Errorf("%s parsed", not)
		}
	}
	// Environment: entries are not subcommands.
	for k := range cmds {
		if k == "FLOCKDECK_UPDATE" || k == "flockdeck" {
			t.Errorf("parsed %q", k)
		}
	}
}

func TestParseDesktopHelpFailsLoudly(t *testing.T) {
	if _, err := parseDesktopHelp("Usage:\n  flockdeck [flags]\n"); err == nil {
		t.Error("help without a Subcommands: heading parsed")
	}
	if _, err := parseDesktopHelp("Subcommands:\n\nEnvironment:\n"); err == nil {
		t.Error("help with no subcommands parsed")
	}
}

// page wraps the way the generated HTML does: inline <code>, and <pre><code>
// blocks with span-highlighted text and escaped brackets.
func page(body string) string {
	return "<html><body>" + body + "</body></html>"
}

const fullPage = `<h2>spawn</h2><pre><code class="language-sh">flockdeck spawn [--worktree &lt;branch&gt;] &lt;task&gt;
</code></pre>
<p>peer-name: <code>flockdeck peer-name &lt;name&gt;</code>, and <code>flockdeck close --finished</code>.</p>
<p>Use <code>flockdeck agents</code> or <code>flockdeck update -check</code>; <code>flockdeck help</code>.</p>
<pre><code>flockdeck
flockdeck-chat chat
</code></pre>
<p><code>flockdeck keys list</code> <code>flockdeck chat</code></p>
<p><code>flockdeck remote enable</code> <code>flockdeck remote pair</code> <code>flockdeck remote status</code>
<code>flockdeck remote devices</code> <code>flockdeck remote revoke &lt;id&gt;</code> <code>flockdeck remote disable</code>
<code>flockdeck remote rename &lt;name&gt;</code> <code>flockdeck remote move &lt;relay&gt;</code></p>`

func TestCompareCLIRealCommandsAgainstCompletePage(t *testing.T) {
	if got := compareCLI(realCommands(t), page(fullPage), nil); len(got) != 0 {
		t.Errorf("complete page reported: %q", got)
	}
}

func TestCompareCLICatchesMissingCommand(t *testing.T) {
	// The peer-name and close gap this check exists for.
	broken := page(`<p><code>flockdeck spawn x</code> <code>flockdeck agents</code> <code>flockdeck keys list</code>
<code>flockdeck chat</code> <code>flockdeck update</code> <code>flockdeck help</code>
<code>flockdeck remote enable</code> <code>flockdeck remote pair</code> <code>flockdeck remote status</code>
<code>flockdeck remote devices</code> <code>flockdeck remote revoke</code> <code>flockdeck remote disable</code>
<code>flockdeck remote rename</code> <code>flockdeck remote move</code></p>`)
	got := compareCLI(realCommands(t), broken, nil)
	wantFinding(t, got, "`peer-name` is not documented")
	wantFinding(t, got, "`close` is not documented")
	if len(got) != 2 {
		t.Errorf("got %q", got)
	}
}

func TestCompareCLICatchesMissingSubcommand(t *testing.T) {
	// remote is documented, but not remote move.
	broken := fullPage
	broken = replace(broken, "<code>flockdeck remote move &lt;relay&gt;</code>", "")
	got := compareCLI(realCommands(t), page(broken), nil)
	if len(got) != 1 {
		t.Fatalf("got %q", got)
	}
	wantFinding(t, got, "`remote move` is not documented")
}

func TestCompareCLICatchesRemovedCommand(t *testing.T) {
	got := compareCLI(realCommands(t), page(fullPage+`<p><code>flockdeck frobnicate now</code> <code>flockdeck remote nuke</code></p>`), nil)
	wantFinding(t, got, "`flockdeck frobnicate`, which the desktop no longer has")
	wantFinding(t, got, "`flockdeck remote nuke`, which the desktop no longer has")
}

func TestCompareCLIIgnoresProseAndFlags(t *testing.T) {
	// "flockdeck is built without a console" is not a subcommand; flags and
	// the -chat binary's name are not either.
	extra := `<p>flockdeck is built without a console; <code>flockdeck -version</code>, <code>flockdeck -C ~/code</code>, <code>flockdeck.exe</code>.</p>`
	if got := compareCLI(realCommands(t), page(fullPage+extra), nil); len(got) != 0 {
		t.Errorf("got %q", got)
	}
}

func TestCompareCLIAllowlist(t *testing.T) {
	noHelp := replace(fullPage, "<code>flockdeck help</code>", "")
	if got := compareCLI(realCommands(t), page(noHelp), nil); len(got) != 1 {
		t.Fatalf("got %q", got)
	}
	allow := map[string]string{"help": "same as -h"}
	if got := compareCLI(realCommands(t), page(noHelp), allow); len(got) != 0 {
		t.Errorf("allowlisted command reported: %q", got)
	}
	wantFinding(t, compareCLI(realCommands(t), page(fullPage), allow), "remove it from the allowlist")
	wantFinding(t, compareCLI(realCommands(t), page(fullPage), map[string]string{"gone": "x"}), "allowlist entry gone")
}

func replace(s, old, repl string) string {
	i := indexOf(s, old)
	if i < 0 {
		panic("test fixture does not contain " + old)
	}
	return s[:i] + repl + s[i+len(old):]
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
