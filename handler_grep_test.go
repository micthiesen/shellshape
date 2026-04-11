package shellshape

import "testing"

func TestGrep(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple", "grep foo file.txt", "grep <pattern> <path>"},
		{"extended", "grep -E 'foo|bar' file.txt", "grep -E <pattern> <path>"},
		{"recursive bundled", "grep -rn pattern /path", "grep -rn <pattern> <path>"},
		{"multiple patterns", "grep -e foo -e bar file.txt", "grep -e <pattern> -e <pattern> <path>"},
		{"pattern only", "grep pattern", "grep <pattern>"},
		{"egrep", "egrep 'foo|bar' file.txt", "egrep <pattern> <path>"},
		{"rg", "rg 'todo' src/", "rg <pattern> <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different patterns collide", func(t *testing.T) {
		a := Normalize("grep -E '(foo|bar)' file")
		b := Normalize("grep -E '(baz|qux)' file")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed to pattern", func(t *testing.T) {
		benign := Normalize("grep foo file.txt")
		subshell := Normalize("grep $(cat pattern.txt) file.txt")
		if benign == subshell {
			t.Error("grep with subshell must not collapse to same shape as grep with literal")
		}
		if !contains(subshell, "$(") {
			t.Errorf("expected subshell marker in %q", subshell)
		}
	})
}

func TestGrepContext(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"-A numeric", "grep -A 5 pattern file.txt", "grep -A N <pattern> <path>"},
		{"-B numeric", "grep -B 3 pattern file.txt", "grep -B N <pattern> <path>"},
		{"-C numeric", "grep -C 10 pattern file.txt", "grep -C N <pattern> <path>"},
		{"-A fused", "grep -A3 pattern file.txt", "grep -A N <pattern> <path>"},
		{"-m flag", "grep -m 5 pattern file.txt", "grep -m N <pattern> <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different numbers collide", func(t *testing.T) {
		a := Normalize("grep -A 5 foo file.txt")
		b := Normalize("grep -A 20 bar file.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})
}
