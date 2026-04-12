package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestSystemdAnalyze(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands
		{"blame", "systemd-analyze blame", "systemd-analyze blame"},
		{"critical-chain", "systemd-analyze critical-chain", "systemd-analyze critical-chain"},
		{"critical-chain unit", "systemd-analyze critical-chain sshd.service", "systemd-analyze critical-chain <unit>"},
		{"plot", "systemd-analyze plot", "systemd-analyze plot"},
		{"dot", "systemd-analyze dot", "systemd-analyze dot"},
		{"unit-paths", "systemd-analyze unit-paths", "systemd-analyze unit-paths"},
		{"security", "systemd-analyze security", "systemd-analyze security"},
		{"security unit", "systemd-analyze security nginx.service", "systemd-analyze security <unit>"},

		// Time expression subcommands
		{"calendar", "systemd-analyze calendar 'Mon *-*-* 00:00:00'", "systemd-analyze calendar <val>"},
		{"timestamp", "systemd-analyze timestamp '2023-01-01 12:00:00'", "systemd-analyze timestamp <val>"},
		{"timespan", "systemd-analyze timespan 5h30min", "systemd-analyze timespan <val>"},

		// Verify takes paths
		{"verify", "systemd-analyze verify /etc/systemd/system/foo.service", "systemd-analyze verify <path>"},

		// Boolean flags
		{"no-pager", "systemd-analyze --no-pager blame", "systemd-analyze --no-pager blame"},
		{"user flag", "systemd-analyze --user security", "systemd-analyze --user security"},
		{"order flag", "systemd-analyze dot --order", "systemd-analyze dot --order"},
		{"require flag", "systemd-analyze dot --require", "systemd-analyze dot --require"},

		// Flags with values
		{"host flag", "systemd-analyze -H root@server blame", "systemd-analyze -H <val> blame"},
		{"root flag", "systemd-analyze --root /mnt/sysroot verify /mnt/sysroot/etc/systemd/system/foo.service", "systemd-analyze --root <path> verify <path>"},

		// Redirect
		{"plot redirect", "systemd-analyze plot > boot.svg", "systemd-analyze plot > <path>"},
		{"dot pipe", "systemd-analyze dot | dot -Tsvg > output.svg", "systemd-analyze dot | dot -Tsvg > <path>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different unit names collide", func(t *testing.T) {
		a := shellshape.Normalize("systemd-analyze security nginx.service")
		b := shellshape.Normalize("systemd-analyze security sshd.service")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different time expressions collide", func(t *testing.T) {
		a := shellshape.Normalize("systemd-analyze calendar 'Mon *-*-* 00:00:00'")
		b := shellshape.Normalize("systemd-analyze calendar 'Fri *-*-* 12:00:00'")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("systemd-analyze security literal-arg")
		subshell := shellshape.Normalize("systemd-analyze security $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
