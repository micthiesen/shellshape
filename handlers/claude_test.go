package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestClaude(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"no args", "claude", "claude"},
		{"print mode", "claude -p 'fix the bug in main.go'", "claude -p <str>"},
		{"long print", "claude --print 'explain this code'", "claude --print <str>"},
		{"model flag", "claude --model sonnet 'hello'", "claude --model sonnet <str>"},
		{"max turns", "claude --max-turns 5 'do the thing'", "claude --max-turns N <str>"},
		{"continue", "claude -c", "claude -c"},
		{"resume session", "claude --resume abc123def456 'continue working'", "claude --resume <val> <str>"},
		{"allowed tools", "claude --allowedTools Bash,Read 'fix it'", "claude --allowedTools <val> <str>"},
		{"mcp subcommand", "claude mcp serve", "claude mcp serve"},
		{"config subcommand", "claude config list", "claude config list"},
		{"print with model", "claude -p --model haiku 'summarize this'", "claude -p --model haiku <str>"},
		{"output format", "claude -p --output-format json 'list files'", "claude -p --output-format <val> <str>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different prompts collide", func(t *testing.T) {
		a := shellshape.Normalize("claude -p 'fix the bug'")
		b := shellshape.Normalize("claude -p 'add a feature'")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("claude -p 'hello'")
		subshell := shellshape.Normalize("claude -p $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
