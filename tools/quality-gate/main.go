// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/larksuite/cli/internal/qualitygate/manifest"
	"github.com/larksuite/cli/internal/qualitygate/report"
	"github.com/larksuite/cli/internal/qualitygate/rules"
)

func main() {
	if len(os.Args) < 2 {
		usageAndExit(2)
	}

	switch os.Args[1] {
	case "manifest":
		configureQualityGateEnvironment()
		m, err := manifest.Collect(context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "quality-gate manifest: %v\n", err)
			os.Exit(2)
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(m); err != nil {
			fmt.Fprintf(os.Stderr, "quality-gate manifest: encode: %v\n", err)
			os.Exit(2)
		}
	case "command-index":
		configureQualityGateEnvironment()
		m, err := manifest.CollectCommandIndex(context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "quality-gate command-index: %v\n", err)
			os.Exit(2)
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(m); err != nil {
			fmt.Fprintf(os.Stderr, "quality-gate command-index: encode: %v\n", err)
			os.Exit(2)
		}
	case "check":
		os.Exit(runCheck(os.Args[2:]))
	default:
		usageAndExit(2)
	}
}

func runCheck(args []string) int {
	configureQualityGateEnvironment()

	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	opts := rules.Options{}
	var printLegacyCommandCandidates bool
	var printLegacyFlagCandidates bool
	fs.StringVar(&opts.Repo, "repo", ".", "repository root")
	fs.StringVar(&opts.CLIBin, "cli-bin", "./lark-cli", "lark-cli binary used for dry-run validation")
	fs.StringVar(&opts.ChangedFrom, "changed-from", "", "base revision for incremental checks")
	fs.StringVar(&opts.FactsOut, "facts-out", "", "write facts JSON to this path")
	fs.BoolVar(&printLegacyCommandCandidates, "print-legacy-command-candidates", false, "print current non-kebab-case hand-authored command candidates")
	fs.BoolVar(&printLegacyFlagCandidates, "print-legacy-flag-candidates", false, "print current non-kebab-case flag candidates")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "quality-gate check: %v\n", err)
		return 2
	}

	if printLegacyCommandCandidates || printLegacyFlagCandidates {
		m, err := manifest.Collect(context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "quality-gate check: %v\n", err)
			return 2
		}
		if printLegacyCommandCandidates {
			for _, line := range rules.LegacyCommandCandidates(*m) {
				fmt.Fprintln(os.Stdout, line)
			}
		}
		if printLegacyFlagCandidates {
			for _, line := range rules.LegacyFlagCandidates(*m) {
				fmt.Fprintln(os.Stdout, line)
			}
		}
		return 0
	}

	diags, facts, err := rules.Run(context.Background(), opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "quality-gate check: %v\n", err)
		return 2
	}
	report.Print(os.Stderr, diags)
	if opts.FactsOut != "" {
		if err := facts.WriteFile(opts.FactsOut); err != nil {
			fmt.Fprintf(os.Stderr, "quality-gate check: write facts: %v\n", err)
			return 2
		}
	}
	return report.ExitCode(diags)
}

func configureQualityGateEnvironment() {
	_ = os.Setenv("LARKSUITE_CLI_REMOTE_META", "off")
	if os.Getenv("LARKSUITE_CLI_CONFIG_DIR") == "" {
		_ = os.Setenv("LARKSUITE_CLI_CONFIG_DIR", filepath.Join(os.TempDir(), "quality-gate-cli-config"))
	}
}

func usageAndExit(code int) {
	fmt.Fprintln(os.Stderr, "usage: quality-gate <manifest|command-index|check>")
	os.Exit(code)
}
