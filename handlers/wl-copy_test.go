package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestWlCopy(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare command", "wl-copy", "wl-copy"},
		{"copy text", `wl-copy "hello world"`, "wl-copy <str>"},
		{"copy with clear", "wl-copy --clear", "wl-copy --clear"},
		{"copy with paste-once", "wl-copy --paste-once something", "wl-copy --paste-once <str>"},
		{"copy with short paste-once", "wl-copy -o something", "wl-copy -o <str>"},
		{"copy with primary", "wl-copy -p text", "wl-copy -p <str>"},
		{"copy with type", "wl-copy -t text/plain something", "wl-copy -t <val> <str>"},
		{"copy with long type", "wl-copy --type text/html content", "wl-copy --type <val> <str>"},
		{"copy with foreground", "wl-copy -f data", "wl-copy -f <str>"},
		{"copy with no-newline", "wl-copy -n text", "wl-copy -n <str>"},
		{"multiple flags", "wl-copy -p -n -o text", "wl-copy -p -n -o <str>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different text collides", func(t *testing.T) {
		a := shellshape.Normalize("wl-copy hello")
		b := shellshape.Normalize("wl-copy goodbye")
		if a != b {
			t.Errorf("expected same shape: %q vs %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("wl-copy hello")
		subshell := shellshape.Normalize("wl-copy $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}

func TestWlPaste(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare command", "wl-paste", "wl-paste"},
		{"with no-newline", "wl-paste -n", "wl-paste -n"},
		{"with list types", "wl-paste -l", "wl-paste -l"},
		{"with primary", "wl-paste -p", "wl-paste -p"},
		{"with type", "wl-paste -t text/plain", "wl-paste -t <val>"},
		{"with long type", "wl-paste --type image/png", "wl-paste --type <val>"},
		{"with watch mode", "wl-paste -w xclip -selection clipboard", "wl-paste -w xclip -selection clipboard"},
		{"with long watch", "wl-paste --watch notify-send clipboard", "wl-paste --watch notify-send clipboard"},
		{"watch with flags", "wl-paste -p -w cat", "wl-paste -p -w cat"},
		{"multiple flags", "wl-paste -n -p -t text/html", "wl-paste -n -p -t <val>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different types collide", func(t *testing.T) {
		a := shellshape.Normalize("wl-paste -t text/plain")
		b := shellshape.Normalize("wl-paste -t image/png")
		if a != b {
			t.Errorf("expected same shape: %q vs %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("wl-paste -n")
		subshell := shellshape.Normalize("wl-paste $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
