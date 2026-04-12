package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestDelta(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage (typically piped, but can take file args)
		{"no args", "delta", "delta"},
		{"two files", "delta file_a.txt file_b.txt", "delta <path>+"},
		{"single file", "delta ./changes.patch", "delta <path>"},

		// Boolean flags
		{"side by side long", "delta --side-by-side", "delta --side-by-side"},
		{"side by side short", "delta -s", "delta -s"},
		{"line numbers long", "delta --line-numbers", "delta --line-numbers"},
		{"line numbers short", "delta -n", "delta -n"},
		{"diff so fancy", "delta --diff-so-fancy", "delta --diff-so-fancy"},
		{"color only", "delta --color-only", "delta --color-only"},

		// Flags with numeric arguments
		{"width", "delta --width 120", "delta --width N"},
		{"tabs", "delta --tabs 4", "delta --tabs N"},

		// Flags with value arguments
		{"paging", "delta --paging always", "delta --paging <val>"},

		// Fused flags
		{"width fused", "delta --width=120", "delta --width=N"},
		{"paging fused", "delta --paging=never", "delta --paging=<val>"},
		{"tabs fused", "delta --tabs=8", "delta --tabs=N"},

		// Combinations
		{"side by side with width", "delta -s --width 80", "delta -s --width N"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different files collide", func(t *testing.T) {
		a := shellshape.Normalize("delta /tmp/a.diff")
		b := shellshape.Normalize("delta /tmp/b.patch")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("delta file.txt")
		subshell := shellshape.Normalize("delta $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
