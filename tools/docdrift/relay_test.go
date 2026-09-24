package main

import (
	"os"
	"strings"
	"testing"
)

func fixture(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func realRelayFlags(t *testing.T) map[string]relayFlag {
	t.Helper()
	flags, err := parseRelayHelp(fixture(t, "relay-v0.2.50-serve-help.txt"))
	if err != nil {
		t.Fatal(err)
	}
	return flags
}

func TestParseRelayHelp(t *testing.T) {
	flags := realRelayFlags(t)
	if len(flags) != 48 {
		t.Errorf("got %d flags, want 48", len(flags))
	}
	for _, c := range []struct{ name, def string }{
		{"-addr", ":8080"},                     // "(:443 by default ...) (default ":8080")": the last group wins
		{"-session-ttl", "720h0m0s"},           // a duration
		{"-max-hosts", "10"},                   // an int
		{"-trust-proxy", ""},                   // a bool prints no default
		{"-memory-limit", ""},                  // "(default: GOMEMLIMIT, or none)" is prose, not a default
		{"-audit-retention", ""},               // "(0: for good)" is prose too
		{"-require-verified-registration", ""}, // flag's backquote quirk renames its type to "flockdeck remote enable"
	} {
		f, ok := flags[c.name]
		if !ok {
			t.Errorf("%s not parsed", c.name)
		} else if f.Default != c.def {
			t.Errorf("%s default = %q, want %q", c.name, f.Default, c.def)
		}
	}
	if got := flags["-require-verified-registration"].Envs; len(got) != 1 || got[0] != "FLOCKDECK_RELAY_REQUIRE_VERIFIED_REGISTRATION" {
		// The POSTMARK_TOKEN it "needs" is not an env var that does the same thing as the flag.
		t.Errorf("-require-verified-registration envs = %v", got)
	}
	if got := flags["-vapid-private-key"].Envs; len(got) != 1 || got[0] != "FLOCKDECK_RELAY_VAPID_PRIVATE_KEY" {
		t.Errorf("-vapid-private-key envs = %v (\"(or, better, X)\")", got)
	}
	if got := flags["-oidc-client-secret-file"].Envs; len(got) != 1 {
		t.Errorf("-oidc-client-secret-file envs = %v", got)
	}
}

func TestParseRelayHelpFailsLoudly(t *testing.T) {
	if _, err := parseRelayHelp("Usage: flockdeck-relay\n  -addr string\n    \tx\n"); err == nil {
		t.Error("help without a 'Flags for serve:' heading parsed; a format change would pass silently")
	}
	if _, err := parseRelayHelp("Flags for serve:\n"); err == nil {
		t.Error("help with no flags parsed")
	}
}

const sampleDoc = "# Configuration\n" +
	"\n" +
	"| Flag | Env | Default | What it does |\n" +
	"| --- | --- | --- | --- |\n" +
	"| `-addr` | | `:8080` (`:443` with `-acme-domain`) | Address. |\n" +
	"| `-tls-cert` / `-tls-key` | | | Your own certificate. |\n" +
	"| `-database` | `FLOCKDECK_RELAY_DATABASE_URL` | | MySQL. |\n" +
	"| `-session-ttl` | | 720h (30 days) | Idle. |\n" +
	"| `-oidc-client-secret-file` | `FLOCKDECK_RELAY_OIDC_CLIENT_SECRET` | | There's no `-oidc-client-secret` flag. |\n" +
	"\n" +
	"Use `-audit` and `flockdeck-relay invite` with `-note x`.\n" +
	"\n" +
	"```\n" +
	"flockdeck-relay serve \\\n" +
	"  -public-url https://relay.example.com \\\n" +
	"  -trust-proxy\n" +
	"```\n" +
	"\n" +
	"```\n" +
	"flockdeck-relay invite -note hello\n" +
	"```\n"

func TestParseConfigDoc(t *testing.T) {
	d := parseConfigDoc(sampleDoc)
	for _, want := range []string{"-addr", "-tls-cert", "-tls-key", "-database", "-session-ttl", "-audit", "-public-url", "-trust-proxy", "-oidc-client-secret-file"} {
		if !d.Mentions[want] {
			t.Errorf("%s not seen", want)
		}
	}
	for _, not := range []string{"-oidc-client-secret", "-note", "-acme-domain"} {
		if d.Mentions[not] && not != "-acme-domain" {
			t.Errorf("%s counted as a claim that the relay has it", not)
		}
	}
	if r := d.Rows["-addr"]; !r.HasDefault || !strings.Contains(r.Default, ":8080") {
		t.Errorf("-addr row = %+v", r)
	}
	if r := d.Rows["-database"]; len(r.Envs) != 1 || !r.HasEnv {
		t.Errorf("-database row = %+v", r)
	}
}

func TestCanonDefault(t *testing.T) {
	same := [][2]string{
		{"720h (30 days)", "720h0m0s"},
		{"10m", "10m0s"},
		{"`:8080` (`:443` with `-acme-domain`)", ":8080"},
		{"off", "0"},
		{"0 (for good)", "0s"},
		{"`info`", `info`},
	}
	for _, c := range same {
		a, _ := canonDefault(c[0])
		b, _ := canonDefault(c[1])
		if a != b {
			t.Errorf("canonDefault(%q)=%q and (%q)=%q, want equal", c[0], a, c[1], b)
		}
	}
	if _, ok := canonDefault("  "); ok {
		t.Error("a blank cell states no default, so there is nothing to compare")
	}
	if a, _ := canonDefault("10"); a == "0" {
		t.Error("10 must not look like zero")
	}
}

func flagsFor(names ...string) map[string]relayFlag {
	m := map[string]relayFlag{}
	for _, n := range names {
		m[n] = relayFlag{Name: n}
	}
	return m
}

func docFor(rows map[string]docRow, mentions ...string) docFlags {
	d := docFlags{Rows: rows, Mentions: map[string]bool{}}
	for _, m := range mentions {
		d.Mentions[m] = true
	}
	for m := range rows {
		d.Mentions[m] = true
	}
	return d
}

func wantFinding(t *testing.T, got []string, substr string) {
	t.Helper()
	for _, g := range got {
		if strings.Contains(g, substr) {
			return
		}
	}
	t.Errorf("no finding containing %q in %q", substr, got)
}

func TestCompareRelayCatchesUndocumentedFlag(t *testing.T) {
	got := compareRelay(flagsFor("-a", "-b"), docFor(nil, "-a"), nil)
	if len(got) != 1 {
		t.Fatalf("got %q", got)
	}
	wantFinding(t, got, "relay flag -b is undocumented")
}

func TestCompareRelayCatchesRemovedFlag(t *testing.T) {
	got := compareRelay(flagsFor("-a"), docFor(nil, "-a", "-old"), nil)
	wantFinding(t, got, "docs mention -old, which the relay no longer has")
}

func TestCompareRelayCatchesWrongDefault(t *testing.T) {
	flags := map[string]relayFlag{"-max-hosts": {Name: "-max-hosts", Default: "10"}}
	doc := docFor(map[string]docRow{"-max-hosts": {Default: "5", HasDefault: true}})
	wantFinding(t, compareRelay(flags, doc, nil), `docs give default "5", the relay's is "10"`)

	doc = docFor(map[string]docRow{"-max-hosts": {Default: "10", HasDefault: true}})
	if got := compareRelay(flags, doc, nil); len(got) != 0 {
		t.Errorf("matching default reported: %q", got)
	}
	// A blank Default cell states nothing, so it is not compared.
	doc = docFor(map[string]docRow{"-max-hosts": {Default: "", HasDefault: true}})
	if got := compareRelay(flags, doc, nil); len(got) != 0 {
		t.Errorf("blank default reported: %q", got)
	}
	// "off" is a bool's zero value: the relay prints no default for it.
	flags = flagsFor("-audit")
	doc = docFor(map[string]docRow{"-audit": {Default: "off", HasDefault: true}})
	if got := compareRelay(flags, doc, nil); len(got) != 0 {
		t.Errorf("off vs no default reported: %q", got)
	}
	doc = docFor(map[string]docRow{"-audit": {Default: "on", HasDefault: true}})
	wantFinding(t, compareRelay(flags, doc, nil), "-audit: docs give default")
}

func TestCompareRelayEnv(t *testing.T) {
	flags := map[string]relayFlag{"-database": {Name: "-database", Usage: "DSN (or FLOCKDECK_RELAY_DATABASE_URL)", Envs: []string{"FLOCKDECK_RELAY_DATABASE_URL"}}}
	doc := docFor(map[string]docRow{"-database": {HasEnv: true}})
	wantFinding(t, compareRelay(flags, doc, nil), "Env column omits")
	doc = docFor(map[string]docRow{"-database": {HasEnv: true, Envs: []string{"FLOCKDECK_RELAY_DATABASE_URL", "FLOCKDECK_RELAY_NOPE"}}})
	wantFinding(t, compareRelay(flags, doc, nil), "FLOCKDECK_RELAY_NOPE, which the relay's help does not")
	// A table without an Env column makes no claim about env vars.
	doc = docFor(map[string]docRow{"-database": {}})
	if got := compareRelay(flags, doc, nil); len(got) != 0 {
		t.Errorf("table with no Env column reported: %q", got)
	}
}

func TestCompareRelayAllowlist(t *testing.T) {
	allow := map[string]string{"-billing-url": "hosted-only"}
	if got := compareRelay(flagsFor("-a", "-billing-url"), docFor(nil, "-a"), allow); len(got) != 0 {
		t.Errorf("allowlisted flag reported: %q", got)
	}
	// An entry that outlives its reason rots the file: fail on it.
	wantFinding(t, compareRelay(flagsFor("-a", "-billing-url"), docFor(nil, "-a", "-billing-url"), allow), "remove it from the allowlist")
	wantFinding(t, compareRelay(flagsFor("-a"), docFor(nil, "-a"), allow), "allowlist entry -billing-url names a flag the relay no longer has")
}

func TestCompareRelayRealHelpAgainstCompleteDoc(t *testing.T) {
	flags := realRelayFlags(t)
	rows := map[string]docRow{}
	var names []string
	for n := range flags {
		names = append(names, n)
	}
	// A doc that names every flag (no tables, so no defaults to disagree on)
	// is clean; dropping one flag is exactly one finding.
	if got := compareRelay(flags, docFor(rows, names...), nil); len(got) != 0 {
		t.Fatalf("complete doc reported: %q", got)
	}
	d := docFor(rows, names...)
	delete(d.Mentions, "-max-streams")
	got := compareRelay(flags, d, nil)
	if len(got) != 1 || !strings.Contains(got[0], "-max-streams is undocumented") {
		t.Errorf("got %q", got)
	}
}

func TestNewFindings(t *testing.T) {
	got := newFindings([]string{"a", "b", "c"}, []string{"b"})
	if len(got) != 2 || got[0] != "a" || got[1] != "c" {
		t.Errorf("got %q", got)
	}
}

func TestLoadAllowRequiresReason(t *testing.T) {
	dir := t.TempDir()
	write := func(s string) string {
		p := dir + "/allow.txt"
		if err := os.WriteFile(p, []byte(s), 0o644); err != nil {
			t.Fatal(err)
		}
		return p
	}
	if _, err := loadAllow(write("# c\n\n-x hosted only\n")); err != nil {
		t.Errorf("good file: %v", err)
	}
	if _, err := loadAllow(write("-x\n")); err == nil {
		t.Error("an entry with no reason was accepted")
	}
	if _, err := loadAllow(write("-x a\n-x b\n")); err == nil {
		t.Error("a duplicate entry was accepted")
	}
}
