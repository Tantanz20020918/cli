// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package manifest

type Source string

const (
	SourceBuiltin  Source = "builtin"
	SourceShortcut Source = "shortcut"
	SourceService  Source = "service"
)

type Manifest struct {
	SchemaVersion int       `json:"schema_version"`
	Commands      []Command `json:"commands"`
}

type Command struct {
	Path          string   `json:"path"`
	CanonicalPath string   `json:"canonical_path,omitempty"`
	Domain        string   `json:"domain,omitempty"`
	Use           string   `json:"use"`
	Short         string   `json:"short,omitempty"`
	Example       string   `json:"example,omitempty"`
	Hidden        bool     `json:"hidden,omitempty"`
	Runnable      bool     `json:"runnable"`
	Source        Source   `json:"source"`
	Generated     bool     `json:"generated,omitempty"`
	Risk          string   `json:"risk,omitempty"`
	Identities    []string `json:"identities,omitempty"`
	Flags         []Flag   `json:"flags,omitempty"`
	DefaultFields []string `json:"default_fields,omitempty"`
}

type Flag struct {
	Name        string              `json:"name"`
	Shorthand   string              `json:"shorthand,omitempty"`
	Usage       string              `json:"usage,omitempty"`
	Hidden      bool                `json:"hidden,omitempty"`
	Required    bool                `json:"required,omitempty"`
	TakesValue  bool                `json:"takes_value"`
	DefValue    string              `json:"default,omitempty"`
	NoOptValue  string              `json:"no_opt_value,omitempty"`
	Annotations map[string][]string `json:"annotations,omitempty"`
}
