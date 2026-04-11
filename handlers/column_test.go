package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestColumn(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare column", "column file.txt", "column <path>"},
		{"table mode", "column -t file.txt", "column -t <path>"},
		{"stdin no args", "column -t", "column -t"},

		// Flags with arguments
		{"delimiter", "column -t -s ',' data.csv", "column -t -s <str> <path>"},
		{"column width", "column -c 80 file.txt", "column -c N <path>"},
		{"output separator", "column -t -o ' | ' file.txt", "column -t -o <str> <path>"},
		{"table columns", "column -t -N 'name,age,city' data.csv", "column -t -N <str> <path>"},
		{"table right align", "column -t -R 1,3 data.csv", "column -t -R <str> <path>"},
		{"table hide", "column -t -H 2 data.csv", "column -t -H <str> <path>"},

		// Long-form flags
		{"long table-columns", "column -t --table-columns 'a,b,c' data.csv", "column -t --table-columns <str> <path>"},
		{"long table-right", "column -t --table-right 1 data.csv", "column -t --table-right <str> <path>"},
		{"long table-hide", "column -t --table-hide 2 data.csv", "column -t --table-hide <str> <path>"},

		// Multiple positionals
		{"multiple files", "column -t file1.txt file2.txt", "column -t <path>+"},

		// Boolean flags
		{"fill columns first", "column -x file.txt", "column -x <path>"},
		{"multiple booleans", "column -t -x -n file.txt", "column -t -x -n <path>"},
		{"bool e and L", "column -t -e -L file.txt", "column -t -e -L <path>"},

		// Combined flags and values
		{"delimiter and output sep", "column -t -s ':' -o ' | ' /etc/passwd", "column -t -s <str> -o <str> <path>"},

		// Redirect
		{"with redirect", "column -t > output.txt", "column -t > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different delimiters collide", func(t *testing.T) {
		a := shellshape.Normalize("column -t -s ',' file.csv")
		b := shellshape.Normalize("column -t -s ':' file.csv")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different files collide", func(t *testing.T) {
		a := shellshape.Normalize("column -t data.csv")
		b := shellshape.Normalize("column -t report.tsv")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("column literal-arg")
		subshell := shellshape.Normalize("column $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
