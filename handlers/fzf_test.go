package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestFzf(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare", "fzf", "fzf"},
		{"multi flag", "fzf -m", "fzf -m"},
		{"query flag short", "fzf -q hello", "fzf -q <val>"},
		{"query flag long", "fzf --query hello", "fzf --query <val>"},
		{"filter flag", "fzf -f pattern", "fzf -f <val>"},
		{"preview flag", `fzf --preview "cat {}"`, "fzf --preview <val>"},
		{"bind flag", `fzf --bind "ctrl-a:select-all"`, "fzf --bind <val>"},
		{"height flag", "fzf --height 40%", "fzf --height <val>"},
		{"layout structural", "fzf --layout reverse", "fzf --layout reverse"},
		{"delimiter flag", "fzf -d :", "fzf -d <val>"},
		{"with-nth flag", "fzf --with-nth 1..3", "fzf --with-nth <val>"},
		{"prompt flag", `fzf --prompt "Pick> "`, "fzf --prompt <val>"},
		{"header flag", `fzf --header "Select a file"`, "fzf --header <val>"},

		// Combined flags
		{"multi with query", "fzf -m -q search", "fzf -m -q <val>"},
		{"complex invocation", `fzf --height 50% --layout reverse --preview "head -20 {}" --bind "ctrl-/:toggle-preview"`,
			"fzf --height <val> --layout reverse --preview <val> --bind <val>"},

		// Fused --flag=value syntax
		{"fused query", "fzf --query=hello", "fzf --query=<val>"},
		{"fused height", "fzf --height=40%", "fzf --height=<val>"},
		{"fused layout", "fzf --layout=reverse", "fzf --layout=reverse"},

		// Boolean flags
		{"boolean flags", "fzf --ansi --exact --no-sort", "fzf --ansi --exact --no-sort"},

		// Positional collapses to <val>
		{"positional arg", "fzf somefile", "fzf <val>"},

		// Redirect
		{"with redirect", "fzf > output.txt", "fzf > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different queries collide", func(t *testing.T) {
		a := shellshape.Normalize("fzf -q hello")
		b := shellshape.Normalize("fzf -q world")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("fzf somefile")
		subshell := shellshape.Normalize("fzf $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
