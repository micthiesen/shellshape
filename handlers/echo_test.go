package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestEcho(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"bare echo", "echo", "echo"},
		{"single word", "echo hi", "echo <str>"},
		{"multiple words", "echo hello world", "echo <str>"},
		{"quoted", `echo "hello world"`, "echo <str>"},
		{"separator style", "echo --- foo ---", "echo <str>"},
		{"with flag", "echo -n foo", "echo -n <str>"},
		{"with redirect", "echo hi > /tmp/out", "echo <str> > <path>"},
		{"printf", `printf "%s\n" foo`, "printf <str>"},
		{"in pipeline", "echo hi && echo bye", "echo <str> && echo <str>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different contents collide", func(t *testing.T) {
		a := shellshape.Normalize("echo backup created")
		b := shellshape.Normalize("echo rerere enabled")
		c := shellshape.Normalize("echo --- main changes ---")
		if a != b || b != c {
			t.Errorf("expected all same: %q, %q, %q", a, b, c)
		}
	})

	t.Run("subshell not collapsed to str", func(t *testing.T) {
		benign := shellshape.Normalize("echo hello")
		subshell := shellshape.Normalize("echo $(rm -rf /tmp/foo)")
		if benign == subshell {
			t.Error("echo with subshell must not collapse to same shape as echo with literal")
		}
		if !contains(subshell, "$(") {
			t.Errorf("expected subshell marker in %q", subshell)
		}
	})
}
