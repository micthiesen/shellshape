package shellshape

import "testing"

func TestFile(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"single file", "file foo.txt", "file <path>"},
		{"multiple files", "file foo.txt bar.jpg /tmp/baz", "file <path>+"},
		{"absolute path", "file /usr/bin/ls", "file <path>"},

		// Boolean flags
		{"brief mode", "file -b foo.txt", "file -b <path>"},
		{"mime type", "file --mime-type image.jpg", "file --mime-type <path>"},
		{"dereference", "file -L /tmp/link", "file -L <path>"},
		{"combined bool flags", "file -bi foo.txt", "file -bi <path>"},
		{"keep going", "file -k foo.txt", "file -k <path>"},
		{"compressed", "file -z archive.tar.gz", "file -z <path>"},

		// Flags with arguments
		{"magic file", "file -m /usr/share/magic foo.txt", "file -m <path>+"},
		{"long magic file", "file --magic-file /custom/magic foo.txt", "file --magic-file <path>+"},
		{"files from", "file -f filelist.txt", "file -f <path>"},
		{"long files from", "file --files-from names.txt", "file --files-from <path>"},
		{"separator", "file -F : foo.txt", "file -F <val> <path>"},
		{"long separator", "file --separator @ foo.txt", "file --separator <val> <path>"},
		{"exclude", "file -e soft foo.txt", "file -e <val> <path>"},
		{"long exclude", "file --exclude compress foo.txt", "file --exclude <val> <path>"},
		{"capital M", "file -M /custom/magic foo.txt", "file -M <path>+"},
		{"P flag", "file -P bytes=1024 foo.txt", "file -P <val> <path>"},

		// Mixed flags and files
		{"mime with magic", "file -i -m /custom/magic photo.png", "file -i -m <path>+"},
		{"multiple excludes", "file -e soft -e compress foo.txt", "file -e <val> -e <val> <path>"},

		// Redirects
		{"with redirect", "file foo.txt > /tmp/out", "file <path> > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision tests
	t.Run("different files collide", func(t *testing.T) {
		a := Normalize("file -i /tmp/photo.png")
		b := Normalize("file -i /var/data/report.csv")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different separators collide", func(t *testing.T) {
		a := Normalize("file -F : foo.txt")
		b := Normalize("file -F @ bar.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("file literal-arg")
		subshell := Normalize("file $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
