// docdrift checks this repository's hand-maintained docs against what the
// released Flockdeck artifacts really do. It reads captured `-h` output, so
// it never runs or downloads anything itself; the scripts beside it fetch
// the artifacts, and the workflows decide when to run them.
//
//	docdrift relay -help relay-help.txt -docs content/self-hosting/configuration.md -allow allow-relay-flags.txt
//	docdrift cli   -help desktop-help.txt -docs app/cli.html -allow allow-desktop-commands.txt
//
// With -baseline-docs, the same check also runs against the docs as the pull
// request's base had them, and only findings the head has and the base did
// not are failures: a PR is not red for drift it did not cause, and a
// finding it introduces or leaves in a file it changed is.
//
// Exit status: 0 nothing (new) to report, 1 findings, 2 could not run.
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: docdrift relay|cli -help FILE -docs FILE -allow FILE [-baseline-docs FILE]")
		os.Exit(2)
	}
	mode := os.Args[1]
	fs := flag.NewFlagSet(mode, flag.ExitOnError)
	helpFile := fs.String("help", "", "captured -h output of the released binary")
	docsFile := fs.String("docs", "", "the docs page to check")
	allowFile := fs.String("allow", "", "allowlist of intentional omissions, each with a reason")
	baseFile := fs.String("baseline-docs", "", "the docs page as the base branch has it; only findings not present there fail")
	fs.Parse(os.Args[2:])

	var check func(docs string) ([]string, error)
	switch mode {
	case "relay":
		check = func(docs string) ([]string, error) { return runRelay(*helpFile, docs, *allowFile) }
	case "cli":
		check = func(docs string) ([]string, error) { return runCLI(*helpFile, docs, *allowFile) }
	default:
		fmt.Fprintf(os.Stderr, "docdrift: unknown check %q\n", mode)
		os.Exit(2)
	}

	findings, err := check(*docsFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "docdrift:", err)
		os.Exit(2)
	}
	report := findings
	if *baseFile != "" {
		base, err := check(*baseFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "docdrift: baseline:", err)
			os.Exit(2)
		}
		report = newFindings(findings, base)
		for _, f := range findings {
			if !contains(report, f) {
				fmt.Printf("note (already so on the base branch, not failing this PR): %s\n", f)
			}
		}
	}
	for _, f := range report {
		fmt.Printf("- %s\n", f)
	}
	if len(report) > 0 {
		os.Exit(1)
	}
	fmt.Printf("%s: no new findings (%d already on the base branch)\n", mode, len(findings))
}

func newFindings(head, base []string) []string {
	var out []string
	for _, f := range head {
		if !contains(base, f) {
			out = append(out, f)
		}
	}
	return out
}

func runRelay(helpFile, docsFile, allowFile string) ([]string, error) {
	help, docs, allow, err := load(helpFile, docsFile, allowFile)
	if err != nil {
		return nil, err
	}
	flags, err := parseRelayHelp(help)
	if err != nil {
		return nil, err
	}
	return compareRelay(flags, parseConfigDoc(docs), allow), nil
}

func runCLI(helpFile, docsFile, allowFile string) ([]string, error) {
	help, docs, allow, err := load(helpFile, docsFile, allowFile)
	if err != nil {
		return nil, err
	}
	cmds, err := parseDesktopHelp(help)
	if err != nil {
		return nil, err
	}
	return compareCLI(cmds, docs, allow), nil
}

func load(helpFile, docsFile, allowFile string) (help, docs string, allow map[string]string, err error) {
	if helpFile == "" || docsFile == "" || allowFile == "" {
		return "", "", nil, fmt.Errorf("-help, -docs and -allow are all required")
	}
	h, err := os.ReadFile(helpFile)
	if err != nil {
		return "", "", nil, err
	}
	d, err := os.ReadFile(docsFile)
	if err != nil {
		return "", "", nil, err
	}
	allow, err = loadAllow(allowFile)
	if err != nil {
		return "", "", nil, err
	}
	return string(h), string(d), allow, nil
}
