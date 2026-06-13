// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package manifest

import (
	"context"
	"testing"
)

func TestCollectContainsDocsFetchAndDryRunFlag(t *testing.T) {
	got, err := Collect(context.Background())
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	cmd := findManifestCommand(got, "docs +fetch")
	if cmd == nil {
		t.Fatalf("docs +fetch not found")
	}
	if !cmd.Runnable {
		t.Fatalf("docs +fetch should be runnable")
	}
	if findManifestFlag(cmd, "dry-run") == nil {
		t.Fatalf("docs +fetch should expose --dry-run")
	}
	if cmd.Source != SourceShortcut {
		t.Fatalf("docs +fetch source = %q, want shortcut", cmd.Source)
	}
}

func TestCollectExcludesGeneratedServiceCommands(t *testing.T) {
	got, err := Collect(context.Background())
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	for _, cmd := range got.Commands {
		if cmd.Source == SourceService || cmd.Generated {
			t.Fatalf("quality-gate manifest should not include generated service command: %#v", cmd)
		}
	}
}

func TestCollectCommandIndexIncludesEmbeddedServiceCommand(t *testing.T) {
	got, err := CollectCommandIndex(context.Background())
	if err != nil {
		t.Fatalf("CollectCommandIndex() error = %v", err)
	}
	cmd := findManifestCommand(got, "drive file.comments create_v2")
	if cmd == nil {
		t.Fatalf("drive file.comments create_v2 not found")
	}
	if cmd.Source != SourceService {
		t.Fatalf("source = %q, want service", cmd.Source)
	}
	if !cmd.Generated {
		t.Fatalf("service command should be marked generated")
	}
	if !cmd.Runnable {
		t.Fatalf("service method command should be runnable")
	}
	for _, name := range []string{"file-token", "params", "data", "dry-run"} {
		if findManifestFlag(cmd, name) == nil {
			t.Fatalf("drive file.comments create_v2 should expose --%s", name)
		}
	}
}

func TestCollectCommandIndexDoesNotInheritGeneratedFromServiceParentForShortcut(t *testing.T) {
	got, err := CollectCommandIndex(context.Background())
	if err != nil {
		t.Fatalf("CollectCommandIndex() error = %v", err)
	}
	cmd := findManifestCommand(got, "docs +fetch")
	if cmd == nil {
		t.Fatalf("docs +fetch not found")
	}
	if cmd.Source != SourceShortcut {
		t.Fatalf("docs +fetch source = %q, want shortcut", cmd.Source)
	}
	if cmd.Generated {
		t.Fatalf("shortcut under service parent must not inherit generated=true")
	}
}

func TestCollectDoesNotInheritGeneratedFromServiceParentForShortcut(t *testing.T) {
	got, err := Collect(context.Background())
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	cmd := findManifestCommand(got, "docs +fetch")
	if cmd == nil {
		t.Fatalf("docs +fetch not found")
	}
	if cmd.Source != SourceShortcut {
		t.Fatalf("docs +fetch source = %q, want shortcut", cmd.Source)
	}
	if cmd.Generated {
		t.Fatalf("shortcut under service parent must not inherit generated=true")
	}
}

func TestCollectIgnoresRuntimeStrictMode(t *testing.T) {
	t.Setenv("LARKSUITE_CLI_CONFIG_DIR", t.TempDir())
	t.Setenv("LARKSUITE_CLI_APP_ID", "dry-run")
	t.Setenv("LARKSUITE_CLI_APP_SECRET", "dry-run")
	t.Setenv("LARKSUITE_CLI_BRAND", "feishu")

	got, err := Collect(context.Background())
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}
	if cmd := findManifestCommand(got, "contact +search-user"); cmd == nil {
		t.Fatal("user-only shortcut missing; manifest collection should not apply runtime strict mode")
	}
}

func findManifestCommand(m *Manifest, path string) *Command {
	for i := range m.Commands {
		if m.Commands[i].Path == path {
			return &m.Commands[i]
		}
	}
	return nil
}

func findManifestFlag(cmd *Command, name string) *Flag {
	for i := range cmd.Flags {
		if cmd.Flags[i].Name == name {
			return &cmd.Flags[i]
		}
	}
	return nil
}
