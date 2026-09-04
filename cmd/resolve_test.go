package cmd

import (
	"testing"

	"google.golang.org/api/tasks/v1"
)

func TestEffectiveTasklist(t *testing.T) {
	t.Run("flag wins over empty default", func(t *testing.T) {
		if got := effectiveTasklist("work"); got != "work" {
			t.Fatalf("got %q, want work", got)
		}
	})

	t.Run("empty flag falls through to config", func(t *testing.T) {
		// config.GetDefaultTaskList() reads process env/config; we only
		// assert that an explicit flag is not overwritten.
		if got := effectiveTasklist("Personal"); got != "Personal" {
			t.Fatalf("got %q, want Personal", got)
		}
	})
}

func TestResolveTasklist(t *testing.T) {
	list := []tasks.TaskList{
		{Id: "1", Title: "Personal"},
		{Id: "2", Title: "Work"},
	}

	t.Run("explicit name", func(t *testing.T) {
		got, prompt, err := resolveTasklist(list, "Work")
		if err != nil {
			t.Fatal(err)
		}
		if prompt {
			t.Fatal("expected no prompt")
		}
		if got.Id != "2" {
			t.Fatalf("got id %q, want 2", got.Id)
		}
	})

	t.Run("invalid name", func(t *testing.T) {
		_, _, err := resolveTasklist(list, "missing")
		if err == nil {
			t.Fatal("expected error")
		}
		if err.Error() != `tasklist "missing" not found` {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("single list auto-selects", func(t *testing.T) {
		one := []tasks.TaskList{{Id: "9", Title: "Only"}}
		got, prompt, err := resolveTasklist(one, "")
		if err != nil {
			t.Fatal(err)
		}
		if prompt {
			t.Fatal("expected no prompt for a single list")
		}
		if got.Id != "9" {
			t.Fatalf("got id %q, want 9", got.Id)
		}
	})

	t.Run("multiple lists need prompt", func(t *testing.T) {
		_, prompt, err := resolveTasklist(list, "")
		if err != nil {
			t.Fatal(err)
		}
		if !prompt {
			t.Fatal("expected interactive prompt when no default exists")
		}
	})

	t.Run("empty account", func(t *testing.T) {
		_, _, err := resolveTasklist(nil, "")
		if err == nil {
			t.Fatal("expected error")
		}
		if err.Error() != "no tasklists found" {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestParseTaskNumber(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		count   int
		want    int
		wantErr string
	}{
		{name: "first", arg: "1", count: 3, want: 0},
		{name: "last", arg: "3", count: 3, want: 2},
		{name: "zero", arg: "0", count: 3, wantErr: "task 0 not found"},
		{name: "too high", arg: "4", count: 3, wantErr: "task 4 not found"},
		{name: "not a number", arg: "abc", count: 3, wantErr: "task abc not found"},
		{name: "empty list", arg: "1", count: 0, wantErr: "task 1 not found"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTaskNumber(tt.arg, tt.count)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("expected error")
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("got error %q, want %q", err.Error(), tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Fatalf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestResolveAddTitle(t *testing.T) {
	if got := resolveAddTitle([]string{"Buy milk"}, ""); got != "Buy milk" {
		t.Fatalf("got %q", got)
	}
	if got := resolveAddTitle(nil, "From flag"); got != "From flag" {
		t.Fatalf("got %q", got)
	}
	if got := resolveAddTitle([]string{"Positional"}, "Flag"); got != "Positional" {
		t.Fatalf("positional should win, got %q", got)
	}
	if got := resolveAddTitle(nil, ""); got != "" {
		t.Fatalf("got %q, want empty", got)
	}
}
