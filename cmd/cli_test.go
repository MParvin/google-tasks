package cmd

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func mustFind(t *testing.T, path ...string) *cobra.Command {
	t.Helper()
	cmd, _, err := rootCmd.Find(path)
	if err != nil {
		t.Fatalf("command %v not found: %v", path, err)
	}
	return cmd
}

func resetCommandFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		_ = f.Value.Set(f.DefValue)
		f.Changed = false
	})
	for _, child := range cmd.Commands() {
		resetCommandFlags(child)
	}
}

func helpOutput(t *testing.T, args ...string) string {
	t.Helper()
	resetCommandFlags(rootCmd)
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("help %v failed: %v\n%s", args, err, buf.String())
	}
	return buf.String()
}

func TestRootHelpShowsNewCommands(t *testing.T) {
	out := helpOutput(t, "--help")

	want := []string{
		"add",
		"clear",
		"done",
		"info",
		"ls",
		"rm",
		"tasklists",
		"undo",
		"update",
		"login",
		"logout",
		"skills",
	}
	for _, name := range want {
		if !strings.Contains(out, name) {
			t.Errorf("root help missing command %q\n%s", name, out)
		}
	}

	if strings.Contains(out, "  tasks ") {
		t.Errorf("root help should hide deprecated tasks namespace\n%s", out)
	}
}

func TestRootHelpDescriptions(t *testing.T) {
	out := helpOutput(t, "--help")
	checks := map[string]string{
		"ls":        "List tasks",
		"add":       "Add a task",
		"done":      "Mark a task as done",
		"undo":      "Mark a completed task as incomplete",
		"update":    "Update an existing task",
		"info":      "View detailed information about a task",
		"rm":        "Delete a task",
		"clear":     "Hide all completed tasks",
		"tasklists": "Manage tasklists",
	}
	for name, desc := range checks {
		if !strings.Contains(out, desc) {
			t.Errorf("root help missing %q description %q\n%s", name, desc, out)
		}
	}
}

func TestLsHelpAndFlags(t *testing.T) {
	out := helpOutput(t, "ls", "--help")
	if !strings.Contains(out, "List tasks") {
		t.Fatalf("ls help missing summary\n%s", out)
	}

	cmd := mustFind(t, "ls")
	for _, name := range []string{"tasklist", "include-completed", "completed", "sort", "format", "max"} {
		if cmd.Flags().Lookup(name) == nil {
			t.Errorf("ls missing --%s", name)
		}
	}
	if cmd.Flags().ShorthandLookup("l") == nil {
		t.Error("ls missing -l shorthand for --tasklist")
	}
}

func TestAddHelpAndFlags(t *testing.T) {
	out := helpOutput(t, "add", "--help")
	if !strings.Contains(out, "Add a task") {
		t.Fatalf("add help missing summary\n%s", out)
	}

	cmd := mustFind(t, "add")
	for _, name := range []string{"tasklist", "title", "note", "due", "repeat", "repeat-count", "repeat-until"} {
		if cmd.Flags().Lookup(name) == nil {
			t.Errorf("add missing --%s", name)
		}
	}
}

func TestTaskCommandFlags(t *testing.T) {
	cases := []struct {
		path  []string
		flags []string
		short map[string]string
	}{
		{path: []string{"done"}, flags: []string{"tasklist"}, short: map[string]string{"l": "tasklist"}},
		{path: []string{"undo"}, flags: []string{"tasklist"}, short: map[string]string{"l": "tasklist"}},
		{path: []string{"update"}, flags: []string{"tasklist", "title", "note", "due"}, short: map[string]string{"l": "tasklist", "t": "title"}},
		{path: []string{"info"}, flags: []string{"tasklist", "include-completed"}, short: map[string]string{"l": "tasklist", "i": "include-completed"}},
		{path: []string{"rm"}, flags: []string{"tasklist"}, short: map[string]string{"l": "tasklist"}},
		{path: []string{"clear"}, flags: []string{"tasklist", "force"}, short: map[string]string{"l": "tasklist", "f": "force"}},
	}

	for _, tt := range cases {
		cmd := mustFind(t, tt.path...)
		for _, name := range tt.flags {
			if cmd.Flags().Lookup(name) == nil && cmd.InheritedFlags().Lookup(name) == nil {
				t.Errorf("%s missing --%s", strings.Join(tt.path, " "), name)
			}
		}
		for shorthand, name := range tt.short {
			flag := cmd.Flags().ShorthandLookup(shorthand)
			if flag == nil {
				flag = cmd.InheritedFlags().ShorthandLookup(shorthand)
			}
			if flag == nil || flag.Name != name {
				t.Errorf("%s missing -%s (--%s)", strings.Join(tt.path, " "), shorthand, name)
			}
		}
	}
}

func TestTasklistsHelp(t *testing.T) {
	out := helpOutput(t, "tasklists", "--help")
	for _, name := range []string{"add", "rm", "update"} {
		if !strings.Contains(out, name) {
			t.Errorf("tasklists help missing %q\n%s", name, out)
		}
	}
}

func TestTasklistsSubcommands(t *testing.T) {
	if mustFind(t, "tasklists").Run == nil {
		t.Fatal("gtasks tasklists should list tasklists without a view subcommand")
	}

	add := mustFind(t, "tasklists", "add")
	if add.Flags().Lookup("title") == nil {
		t.Fatal("tasklists add missing --title")
	}
	if add.Flags().ShorthandLookup("t") == nil {
		t.Fatal("tasklists add missing -t")
	}

	update := mustFind(t, "tasklists", "update")
	if update.Flags().Lookup("title") == nil {
		t.Fatal("tasklists update missing --title")
	}

	mustFind(t, "tasklists", "rm")
	view := mustFind(t, "tasklists", "view")
	if view.Deprecated == "" {
		t.Fatal("tasklists view should be marked deprecated")
	}
}

func TestLegacyTasksAliases(t *testing.T) {
	mapping := map[string]string{
		"view":   "gtasks ls",
		"add":    "gtasks add",
		"done":   "gtasks done",
		"undo":   "gtasks undo",
		"update": "gtasks update",
		"info":   "gtasks info",
		"rm":     "gtasks rm",
		"clear":  "gtasks clear",
	}

	parent := mustFind(t, "tasks")
	if !parent.Hidden {
		t.Fatal("tasks namespace should be hidden from root help")
	}

	for sub, instead := range mapping {
		cmd := mustFind(t, "tasks", sub)
		if cmd.Deprecated == "" {
			t.Errorf("tasks %s should be deprecated", sub)
		}
		if !strings.Contains(cmd.Deprecated, instead) {
			t.Errorf("tasks %s deprecation %q should mention %s", sub, cmd.Deprecated, instead)
		}
	}

	legacyView := mustFind(t, "tasks", "view")
	if legacyView.Flags().Lookup("include-completed") == nil {
		t.Error("legacy tasks view should keep --include-completed")
	}
	if parent.PersistentFlags().Lookup("tasklist") == nil {
		t.Error("legacy tasks should keep persistent --tasklist")
	}
}

func TestRejectsExtraArgs(t *testing.T) {
	cases := [][]string{
		{"ls", "extra"},
		{"clear", "extra"},
		{"done", "1", "2"},
		{"undo", "1", "2"},
		{"info", "1", "2"},
		{"update", "1", "2"},
		{"rm", "1", "2"},
		{"add", "one", "two"},
		{"tasklists", "add", "one", "two"},
		{"tasklists", "rm", "one", "two"},
		{"tasklists", "update", "one", "two"},
	}

	for _, args := range cases {
		resetCommandFlags(rootCmd)
		rootCmd.SetOut(io.Discard)
		rootCmd.SetErr(io.Discard)
		rootCmd.SetArgs(args)
		if err := rootCmd.Execute(); err == nil {
			t.Errorf("expected error for extra args: %v", args)
		}
	}
}

func TestCompletionIncludesNewCommands(t *testing.T) {
	var buf bytes.Buffer
	if err := rootCmd.GenBashCompletion(&buf); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, name := range []string{"ls", "add", "done", "undo", "update", "info", "rm", "clear", "tasklists"} {
		if !strings.Contains(out, name) {
			t.Errorf("bash completion missing %q", name)
		}
	}
}

func TestCommandsShareImplementationsWithLegacy(t *testing.T) {
	pairs := []struct {
		modern []string
		legacy []string
	}{
		{[]string{"ls"}, []string{"tasks", "view"}},
		{[]string{"add"}, []string{"tasks", "add"}},
		{[]string{"done"}, []string{"tasks", "done"}},
		{[]string{"undo"}, []string{"tasks", "undo"}},
		{[]string{"update"}, []string{"tasks", "update"}},
		{[]string{"info"}, []string{"tasks", "info"}},
		{[]string{"rm"}, []string{"tasks", "rm"}},
		{[]string{"clear"}, []string{"tasks", "clear"}},
	}

	for _, tt := range pairs {
		modern := mustFind(t, tt.modern...)
		legacy := mustFind(t, tt.legacy...)
		if modern.Run == nil || legacy.Run == nil {
			t.Fatalf("%s / %s missing Run", strings.Join(tt.modern, " "), strings.Join(tt.legacy, " "))
		}
	}
}
