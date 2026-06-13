// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package main

import (
	"os"
	"testing"
)

func TestConfigureQualityGateEnvironmentForcesDeterministicRegistry(t *testing.T) {
	t.Setenv("LARKSUITE_CLI_REMOTE_META", "on")
	t.Setenv("LARKSUITE_CLI_CONFIG_DIR", "")

	configureQualityGateEnvironment()

	if got := os.Getenv("LARKSUITE_CLI_REMOTE_META"); got != "off" {
		t.Fatalf("LARKSUITE_CLI_REMOTE_META = %q, want off", got)
	}
	if got := os.Getenv("LARKSUITE_CLI_CONFIG_DIR"); got == "" {
		t.Fatal("LARKSUITE_CLI_CONFIG_DIR was not set")
	}
}
