package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestComm(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"two files", "comm file1.txt file2.txt", "comm <path>+"},
		{"absolute paths", "comm /tmp/sorted1.txt /tmp/sorted2.txt", "comm <path>+"},

		// Digit flags (not caught by isFlagToken)
		{"-1 suppress col1", "comm -1 file1.txt file2.txt", "comm -1 <path>+"},
		{"-2 suppress col2", "comm -2 file1.txt file2.txt", "comm -2 <path>+"},
		{"-3 suppress col3", "comm -3 file1.txt file2.txt", "comm -3 <path>+"},
		{"-12 fused", "comm -12 file1.txt file2.txt", "comm -12 <path>+"},
		{"-23 fused", "comm -23 file1.txt file2.txt", "comm -23 <path>+"},
		{"-123 fused", "comm -123 file1.txt file2.txt", "comm -123 <path>+"},

		// Boolean letter flags
		{"-i case insensitive", "comm -i file1.txt file2.txt", "comm -i <path>+"},

		// Mixed digit and letter flags
		{"-23i fused", "comm -23i file1.txt file2.txt", "comm -23i <path>+"},
		{"-1 -2 separate", "comm -1 -2 file1.txt file2.txt", "comm -1 -2 <path>+"},
		{"-1 -2 -3 -i all separate", "comm -1 -2 -3 -i a.txt b.txt", "comm -1 -2 -3 -i <path>+"},

		// Stdin dash
		{"stdin dash first", "comm - file2.txt", "comm - <path>"},
		{"stdin dash second", "comm file1.txt -", "comm <path> -"},

		// Redirects
		{"redirect output", "comm -12 file1.txt file2.txt > out.txt", "comm -12 <path>+ > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different files collide", func(t *testing.T) {
		a := shellshape.Normalize("comm -12 names.txt ids.txt")
		b := shellshape.Normalize("comm -12 config.yaml settings.yaml")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("comm file1.txt file2.txt")
		subshell := shellshape.Normalize("comm $(dangerous-command) file2.txt")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
