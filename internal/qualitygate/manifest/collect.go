// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package manifest

import (
	"context"
	"io"
	"sort"
	"strings"

	rootcmd "github.com/larksuite/cli/cmd"
	"github.com/larksuite/cli/internal/cmdmeta"
	"github.com/larksuite/cli/internal/cmdutil"
	"github.com/larksuite/cli/internal/registry"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Collect(ctx context.Context) (*Manifest, error) {
	root := rootcmd.Build(ctx, cmdutil.InvocationContext{},
		rootcmd.WithIO(strings.NewReader(""), io.Discard, io.Discard),
		rootcmd.WithoutPlugins(),
		rootcmd.WithoutStrictMode(),
		rootcmd.WithoutServiceCommands(),
	)

	return collectFromRoot(root), nil
}

func CollectCommandIndex(ctx context.Context) (*Manifest, error) {
	root := rootcmd.Build(ctx, cmdutil.InvocationContext{},
		rootcmd.WithIO(strings.NewReader(""), io.Discard, io.Discard),
		rootcmd.WithoutPlugins(),
		rootcmd.WithoutStrictMode(),
		rootcmd.WithServiceCatalog(registry.EmbeddedCatalog()),
	)

	return collectFromRoot(root), nil
}

func collectFromRoot(root *cobra.Command) *Manifest {
	var commands []Command
	walkCommands(root, func(c *cobra.Command) {
		if c == root {
			return
		}
		commands = append(commands, commandFromCobra(c, nil))
	})
	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Path < commands[j].Path
	})

	return &Manifest{SchemaVersion: 1, Commands: commands}
}

func walkCommands(root *cobra.Command, visit func(*cobra.Command)) {
	visit(root)
	for _, child := range root.Commands() {
		walkCommands(child, visit)
	}
}

func commandFromCobra(c *cobra.Command, defaultFields map[string][]string) Command {
	path := strings.TrimPrefix(c.CommandPath(), "lark-cli ")
	source := SourceBuiltin
	if s, ok := cmdmeta.SourceOf(c); ok {
		source = Source(s)
	}
	entry := Command{
		Path:          path,
		CanonicalPath: canonicalCommandPath(path),
		Domain:        commandDomain(c, path, source),
		Use:           c.Use,
		Short:         c.Short,
		Example:       c.Example,
		Hidden:        c.Hidden,
		Runnable:      c.Runnable(),
		Source:        source,
		Generated:     cmdmeta.Generated(c),
		Identities:    cmdmeta.Identities(c),
		DefaultFields: defaultFields[path],
	}
	if risk, ok := cmdmeta.Risk(c); ok {
		entry.Risk = risk
	}

	c.Flags().VisitAll(func(f *pflag.Flag) {
		entry.Flags = append(entry.Flags, flagFromPFlag(f))
	})
	c.InheritedFlags().VisitAll(func(f *pflag.Flag) {
		if findFlag(entry.Flags, f.Name) == nil {
			entry.Flags = append(entry.Flags, flagFromPFlag(f))
		}
	})
	sort.Slice(entry.Flags, func(i, j int) bool {
		return entry.Flags[i].Name < entry.Flags[j].Name
	})
	return entry
}

func commandDomain(c *cobra.Command, path string, source Source) string {
	if domain := cmdmeta.Domain(c); domain != "" {
		return domain
	}
	if source == SourceService {
		if first, _, ok := strings.Cut(path, " "); ok {
			return first
		}
		return path
	}
	return ""
}

func canonicalCommandPath(path string) string {
	parts := strings.Fields(path)
	for i, part := range parts {
		prefix := ""
		if strings.HasPrefix(part, "+") {
			prefix = "+"
			part = strings.TrimPrefix(part, "+")
		}
		part = strings.ReplaceAll(part, ".", "-")
		part = strings.ReplaceAll(part, "_", "-")
		parts[i] = prefix + part
	}
	return strings.Join(parts, " ")
}

func flagFromPFlag(f *pflag.Flag) Flag {
	return Flag{
		Name:        f.Name,
		Shorthand:   f.Shorthand,
		Usage:       f.Usage,
		Hidden:      f.Hidden,
		Required:    hasAnnotation(f, cobra.BashCompOneRequiredFlag),
		TakesValue:  f.NoOptDefVal == "",
		DefValue:    f.DefValue,
		NoOptValue:  f.NoOptDefVal,
		Annotations: cloneAnnotations(f.Annotations),
	}
}

func findFlag(flags []Flag, name string) *Flag {
	for i := range flags {
		if flags[i].Name == name {
			return &flags[i]
		}
	}
	return nil
}

func hasAnnotation(f *pflag.Flag, key string) bool {
	if f.Annotations == nil {
		return false
	}
	values, ok := f.Annotations[key]
	return ok && len(values) > 0
}

func cloneAnnotations(in map[string][]string) map[string][]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string][]string, len(in))
	for key, values := range in {
		out[key] = append([]string(nil), values...)
	}
	return out
}
