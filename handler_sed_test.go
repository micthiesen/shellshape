package shellshape

import "testing"

func TestSed(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"address range", "sed -n 100,230p", "sed -n <sed-expr>"},
		{"address range with file", "sed -n 50,80p file.txt", "sed -n <sed-expr> <path>"},
		{"e flag", "sed -e 's/foo/bar/g' file.txt", "sed -e <sed-expr> <path>"},
		{"positional script", "sed 's/foo/bar/g' file.txt", "sed <sed-expr> <path>"},
		{"in place", "sed -i 's/foo/bar/g' file.txt", "sed -i <sed-expr> <path>"},
		{"sed -n single address", "sed -n 5p", "sed -n <sed-expr>"},
		{"sed -n dollar address", "sed -n '$p'", "sed -n <sed-expr>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different scripts collide", func(t *testing.T) {
		a := Normalize("sed -n 100,230p foo.txt")
		b := Normalize("sed -n 50,80p bar.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("piped sed ranges collide", func(t *testing.T) {
		a := Normalize("git show HEAD:foo.ts | sed -n 1,80p")
		b := Normalize("git show HEAD:foo.ts | sed -n 100,230p")
		c := Normalize("git show HEAD:foo.ts | sed -n 250,361p")
		if a != b || b != c {
			t.Errorf("expected all same: %q, %q, %q", a, b, c)
		}
	})

	t.Run("continuation collapses", func(t *testing.T) {
		cmd := "sed -i '' \\\n  -e 's/foo/bar/g' \\\n  -e 's/baz/qux/g' file.txt"
		got := Normalize(cmd)
		want := "sed -i <sed-expr> -e <sed-expr> -e <sed-expr> <path>"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("different continuations collide", func(t *testing.T) {
		a := Normalize("sed -i '' \\\n  -e 's/guest/orphan/g' \\\n  -e 's/Guest/Orphan/g' file.txt")
		b := Normalize("sed -i '' \\\n  -e 's/foo/bar/g' \\\n  -e 's/baz/qux/g' file.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})
}
