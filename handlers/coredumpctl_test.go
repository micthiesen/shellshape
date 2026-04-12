package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestCoredumpctl(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands
		{"list", "coredumpctl list", "coredumpctl list"},
		{"info", "coredumpctl info", "coredumpctl info"},
		{"dump", "coredumpctl dump", "coredumpctl dump"},
		{"debug", "coredumpctl debug", "coredumpctl debug"},
		{"gdb", "coredumpctl gdb", "coredumpctl gdb"},

		// Positionals: PIDs become N, executable names become <val>
		{"list with program", "coredumpctl list myapp", "coredumpctl list <val>"},
		{"info with PID", "coredumpctl info 1234", "coredumpctl info N"},
		{"debug with program", "coredumpctl debug firefox", "coredumpctl debug <val>"},
		{"gdb with PID", "coredumpctl gdb 5678", "coredumpctl gdb N"},
		{"dump with program", "coredumpctl dump mycrash", "coredumpctl dump <val>"},

		// Flags with arguments
		{"output flag", "coredumpctl dump --output /tmp/core.dump myapp", "coredumpctl dump --output <path> <val>"},
		{"output short", "coredumpctl dump -o /tmp/core myapp", "coredumpctl dump -o <path> <val>"},
		{"since flag", "coredumpctl list --since yesterday", "coredumpctl list --since <val>"},
		{"until flag", "coredumpctl list --until '2024-01-01'", "coredumpctl list --until <val>"},
		{"directory flag", "coredumpctl list -D /var/log/journal", "coredumpctl list -D <path>"},
		{"debugger flag", "coredumpctl debug --debugger lldb myapp", "coredumpctl debug --debugger <val>+"},
		{"field flag", "coredumpctl list -F COREDUMP_EXE", "coredumpctl list -F <val>"},

		// Fused flags (--flag=value)
		{"fused output", "coredumpctl dump --output=/tmp/core myapp", "coredumpctl dump --output=<path> <val>"},

		// Boolean flags
		{"reverse flag", "coredumpctl list -r", "coredumpctl list -r"},
		{"quiet flag", "coredumpctl list -q", "coredumpctl list -q"},
		{"no-pager", "coredumpctl list --no-pager", "coredumpctl list --no-pager"},
		{"no-legend", "coredumpctl list --no-legend", "coredumpctl list --no-legend"},
		{"last entry", "coredumpctl info -1", "coredumpctl info -1"},

		// Multiple positionals
		{"multiple PIDs", "coredumpctl info 1234 5678", "coredumpctl info N N"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different programs collide", func(t *testing.T) {
		a := shellshape.Normalize("coredumpctl debug firefox")
		b := shellshape.Normalize("coredumpctl debug nginx")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different PIDs collide", func(t *testing.T) {
		a := shellshape.Normalize("coredumpctl info 1234")
		b := shellshape.Normalize("coredumpctl info 9999")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("coredumpctl debug myapp")
		subshell := shellshape.Normalize("coredumpctl debug $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
