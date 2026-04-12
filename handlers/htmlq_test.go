package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestHtmlq(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"selector only", "htmlq 'div.content'", "htmlq <val>"},
		{"text mode", "htmlq -t h1", "htmlq -t <val>"},
		{"text mode long", "htmlq --text h1", "htmlq --text <val>"},
		{"pretty print", "htmlq -p body", "htmlq -p <val>"},

		// Attribute flag (structural - name kept verbatim)
		{"attribute href", "htmlq -a href 'a.link'", "htmlq -a href <val>"},
		{"attribute src", "htmlq --attribute src img", "htmlq --attribute src <val>"},

		// Value flags (flag arg + selector positional collapse together)
		{"remove nodes", "htmlq -r script body", "htmlq -r <val>+"},
		{"base url", "htmlq -B https://example.com 'a.link'", "htmlq -B <val>+"},
		{"base url long", "htmlq --base https://example.com 'a.link'", "htmlq --base <val>+"},

		// Path flags
		{"filename", "htmlq -f page.html title", "htmlq -f <path> <val>"},
		{"filename long", "htmlq --filename /tmp/page.html title", "htmlq --filename <path> <val>"},

		// Combined
		{"text with remove", "htmlq -t -r script 'div.main'", "htmlq -t -r <val>+"},
		{"attribute with pretty", "htmlq -p -a href a", "htmlq -p -a href <val>"},

		// Redirect (piped usage)
		{"with pipe input redirect", "htmlq title < page.html", "htmlq <val> < <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different selectors collide", func(t *testing.T) {
		a := shellshape.Normalize("htmlq 'div.content'")
		b := shellshape.Normalize("htmlq 'span.title'")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("htmlq div")
		subshell := shellshape.Normalize("htmlq $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
