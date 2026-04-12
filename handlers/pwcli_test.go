package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestPwCli(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands
		{"info", "pw-cli info 42", "pw-cli info <val>"},
		{"list nodes", "pw-cli list Node", "pw-cli list Node"},
		{"list ports", "pw-cli list Port", "pw-cli list Port"},
		{"dump", "pw-cli dump", "pw-cli dump"},
		{"dump type", "pw-cli dump Node", "pw-cli dump Node"},

		// Connection management
		{"connect", "pw-cli connect 35 42", "pw-cli connect <val>+"},
		{"disconnect", "pw-cli disconnect 17", "pw-cli disconnect <val>"},

		// Node management
		{"destroy-node", "pw-cli destroy-node 42", "pw-cli destroy-node <val>"},
		{"create-node", "pw-cli create-node adapter node.name=my-source", "pw-cli create-node adapter <val>"},

		// Param operations
		{"set-param", "pw-cli set-param 42 Props {audio.volume:0.5}", "pw-cli set-param <val> Props <val>"},
		{"enum-params", "pw-cli enum-params 42 Props", "pw-cli enum-params <val> Props"},
		{"enum-params route", "pw-cli enum-params 42 Route", "pw-cli enum-params <val> Route"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision test
	t.Run("different node IDs collide", func(t *testing.T) {
		a := shellshape.Normalize("pw-cli info 42")
		b := shellshape.Normalize("pw-cli info 99")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("pw-cli info 42")
		subshell := shellshape.Normalize("pw-cli info $(evil-cmd)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
