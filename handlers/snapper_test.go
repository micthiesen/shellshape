package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestSnapper(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands
		{"list", "snapper list", "snapper list"},
		{"create simple", "snapper create", "snapper create"},
		{"create with desc", "snapper create -d before-update", "snapper create -d <val>"},
		{"create typed", "snapper create -t single -d backup", "snapper create -t single -d <val>"},
		{"delete by id", "snapper delete 42", "snapper delete N"},
		{"delete range", "snapper delete 10-20", "snapper delete N"},
		{"diff two ids", "snapper diff 10..20", "snapper diff N"},
		{"status", "snapper status 10..20", "snapper status N"},
		{"rollback", "snapper rollback 15", "snapper rollback N"},
		{"undochange", "snapper undochange 10..20", "snapper undochange N"},

		// Global flags before subcommand
		{"config before sub", "snapper -c home list", "snapper -c <val> list"},
		{"config long before sub", "snapper --config home create -d test", "snapper --config <val> create -d <val>"},
		{"verbose", "snapper -v list", "snapper -v list"},
		{"no-dbus", "snapper --no-dbus list", "snapper --no-dbus list"},

		// Output flags
		{"csvout", "snapper --csvout list", "snapper --csvout list"},
		{"jsonout", "snapper --jsonout list", "snapper --jsonout list"},

		// Cleanup
		{"cleanup", "snapper cleanup number", "snapper cleanup <val>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different configs collide", func(t *testing.T) {
		a := shellshape.Normalize("snapper -c home list")
		b := shellshape.Normalize("snapper -c root list")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different snapshot ids collide", func(t *testing.T) {
		a := shellshape.Normalize("snapper delete 42")
		b := shellshape.Normalize("snapper delete 99")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("snapper delete 42")
		subshell := shellshape.Normalize("snapper delete $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
