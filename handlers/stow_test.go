package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestStow(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage - package names collapse to <pkg>
		{"single package", "stow zsh", "stow <pkg>"},
		{"multiple packages", "stow zsh vim claude", "stow <pkg>+"},

		// Boolean flags
		{"delete flag", "stow -D zsh", "stow -D <pkg>"},
		{"restow flag", "stow -R zsh vim", "stow -R <pkg>+"},
		{"simulate flag", "stow -n zsh", "stow -n <pkg>"},
		{"adopt flag", "stow --adopt zsh", "stow --adopt <pkg>"},
		{"no-folding flag", "stow --no-folding zsh", "stow --no-folding <pkg>"},
		{"dotfiles flag", "stow --dotfiles zsh", "stow --dotfiles <pkg>"},

		// Flags with path arguments
		{"dir flag short", "stow -d ~/.dotfiles zsh", "stow -d <path> <pkg>"},
		{"dir flag long", "stow --dir=/usr/local/stow zsh", "stow --dir=<val> <pkg>"},
		{"target flag short", "stow -t ~ zsh", "stow -t <path> <pkg>"},
		{"target flag long", "stow --target=/usr/local zsh", "stow --target=<val> <pkg>"},
		{"both dir and target", "stow -d ~/dotfiles -t ~ zsh vim", "stow -d <path> -t <path> <pkg>+"},

		// Flags with value arguments
		{"ignore flag", "stow --ignore='.git' zsh", "stow --ignore=<val> <pkg>"},
		{"defer flag", "stow --defer=man zsh", "stow --defer=<val> <pkg>"},
		{"override flag", "stow --override=man zsh", "stow --override=<val> <pkg>"},
		{"ignore separate", "stow --ignore .git zsh", "stow --ignore <val> <pkg>"},

		// Combined flags
		{"verbose and simulate", "stow -n -v -d ~/dotfiles zsh", "stow -n -v -d <path> <pkg>"},
		{"restow with target", "stow -R -t /usr/local zsh vim", "stow -R -t <path> <pkg>+"},

		// Numeric --flag=value normalizes the value
		{"verbose flag=value", "stow --verbose=5 zsh", "stow --verbose=N <pkg>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision test: different paths → same shape
	t.Run("different paths collide", func(t *testing.T) {
		a := shellshape.Normalize("stow -d ~/.dotfiles -t /home/user zsh")
		b := shellshape.Normalize("stow -d /opt/stow -t /usr/local zsh")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Collision test: different ignore patterns → same shape
	t.Run("different ignore patterns collide", func(t *testing.T) {
		a := shellshape.Normalize("stow --ignore='.git' zsh")
		b := shellshape.Normalize("stow --ignore='.svn' zsh")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("stow zsh")
		subshell := shellshape.Normalize("stow $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
