package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestFind(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"name glob", "find /src -name '*.ts'", "find <path> -name <pattern>"},
		{"name wildcard", "find /src -name '*buyside*'", "find <path> -name <pattern>"},
		{"path glob", "find /src -path '*notifications*'", "find <path> -path <pattern>"},
		{"iname glob", "find /src -iname '*.Test.*'", "find <path> -iname <pattern>"},
		{"type preserved", "find /src -type f -name '*.ts'", "find <path> -type f -name <pattern>"},
		{"not path", "find /src -name 'task*' -not -path '*/node_modules/*'", "find <path> -name <pattern> -not -path <pattern>"},
		{"with pipe", "find /src -name '*.ts' | head -20", "find <path> -name <pattern> | head N"},
		{"delete", "find /src -name '*.pyc' -delete", "find <path> -name <pattern> -delete"},
		{"subshell preserved", "find $(pwd) -name '*.ts'", "find $(pwd) -name <pattern>"},

		// Bare directory name as search root collapses to <path>.
		{"bare dir root", "find packages -name '*.ts'", "find <path> -name <pattern>"},
		{"current dir root", "find . -type f", "find <path> -type f"},
		{"multiple bare roots", "find pkg1 pkg2 -name '*.go'", "find <path>+ -name <pattern>"},

		// -exec / -execdir / -ok / -okdir consume the body up to `;` or `+`.
		{"exec semicolon", `find /src -name '*.ts' -exec grep -l foo {} \;`, "find <path> -name <pattern> -exec <cmd> ;"},
		{"exec plus", `find /src -exec grep -l foo {} +`, "find <path> -exec <cmd> +"},
		{"execdir semicolon", `find /src -type f -execdir rm {} \;`, "find <path> -type f -execdir <cmd> ;"},
		{"ok semicolon", `find /src -name '*.pyc' -ok rm {} \;`, "find <path> -name <pattern> -ok <cmd> ;"},
		{"exec with pipe pattern", `find /src -type f -name '*.ts' -exec grep -l 'foo\|bar' {} \;`, "find <path> -type f -name <pattern> -exec <cmd> ;"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different globs collide", func(t *testing.T) {
		a := shellshape.Normalize("find /src -name '*buyside*' -o -name '*Buyside*' 2>/dev/null")
		b := shellshape.Normalize("find /src -name '*notifications*' -o -name '*modelReady*' 2>/dev/null")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("complex with or", func(t *testing.T) {
		got := shellshape.Normalize("find /src -path '*notifications*' -name '*modelReady*' -o -path '*notifications*' -name '*modelNotSupported*' | head -10")
		want := "find <path> -path <pattern> -name <pattern> -o -path <pattern> -name <pattern> | head N"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("different exec bodies collide", func(t *testing.T) {
		a := shellshape.Normalize(`find /src -name '*.ts' -exec grep -l foo {} \;`)
		b := shellshape.Normalize(`find /src -name '*.ts' -exec grep -l 'export.*bar\|export.*baz' {} \;`)
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("exec subshell preserved", func(t *testing.T) {
		// Note: -exec body is collapsed to <cmd>, but any subshell *outside*
		// the exec body (as a search root) must still be preserved.
		got := shellshape.Normalize(`find $(pwd) -name '*.ts' -exec rm {} \;`)
		want := "find $(pwd) -name <pattern> -exec <cmd> ;"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
