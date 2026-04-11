package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestTask(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"single task", "task build", "task build"},
		{"multiple tasks", "task build test lint", "task build test lint"},
		{"list tasks", "task --list", "task --list"},
		{"list short", "task -l", "task -l"},
		{"list all", "task --list-all", "task --list-all"},

		// Boolean flags
		{"force", "task -f deploy", "task -f deploy"},
		{"watch", "task -w build", "task -w build"},
		{"dry run", "task --dry build", "task --dry build"},
		{"verbose", "task -v build", "task -v build"},
		{"silent", "task -s build", "task -s build"},
		{"parallel", "task -p build test", "task -p build test"},

		// Flags with arguments
		{"dir flag", "task -d /some/path build", "task -d <path> build"},
		{"dir long", "task --dir /some/path build", "task --dir <path> build"},
		{"taskfile flag", "task -t custom.yml build", "task -t <path> build"},
		{"taskfile long", "task --taskfile custom.yml build", "task --taskfile <path> build"},
		{"concurrency", "task -C 4 build test", "task -C N build test"},
		{"concurrency long", "task --concurrency 4 build test", "task --concurrency N build test"},
		{"output flag", "task -o group build", "task -o <val> build"},
		{"interval flag", "task -I 500ms build", "task -I <val> build"},
		{"color flag", "task -c false build", "task -c <val> build"},
		{"sort flag", "task --sort alphanumeric --list", "task --sort <val> --list"},

		// Variables (KEY=VALUE)
		{"variable", "task build FOO=bar", "task build <var>=<val>"},
		{"multiple variables", "task deploy ENV=prod VERSION=1.2.3", "task deploy <var>=<val> <var>=<val>"},

		// CLI_ARGS (after --)
		{"cli args", "task test -- -v -count=1", "task test -- <str>+"},
		{"cli args data", "task deploy -- arg1 arg2", "task deploy -- <str>+"},

		// Fused long flags
		{"fused dir", "task --dir=/some/path build", "task --dir=<path> build"},
		{"fused output", "task --output=group build", "task --output=<val> build"},
		{"fused concurrency", "task --concurrency=4 build", "task --concurrency=N build"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision tests
	t.Run("different dirs collide", func(t *testing.T) {
		a := shellshape.Normalize("task -d /home/alice/project build")
		b := shellshape.Normalize("task -d /home/bob/project build")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different variables collide", func(t *testing.T) {
		a := shellshape.Normalize("task deploy ENV=staging VERSION=1.0")
		b := shellshape.Normalize("task deploy ENV=prod VERSION=2.5")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("task build")
		subshell := shellshape.Normalize("task $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
