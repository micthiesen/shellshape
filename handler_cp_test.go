package shellshape

import "testing"

func TestCp(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare cp", "cp", "cp"},
		{"two files", "cp file.txt backup.txt", "cp <path>+"},
		{"absolute paths", "cp /etc/hosts /tmp/hosts", "cp <path>+"},
		{"relative path", "cp ./config.yaml /tmp/", "cp <path>+"},
		{"tilde path", "cp ~/file.txt /tmp/", "cp <path>+"},

		// Flags (all boolean)
		{"recursive", "cp -r src/ /tmp/backup/", "cp -r <path>+"},
		{"recursive uppercase", "cp -R src/ /tmp/backup/", "cp -R <path>+"},
		{"archive", "cp -a /home/user/ /backup/", "cp -a <path>+"},
		{"force", "cp -f file.txt /tmp/", "cp -f <path>+"},
		{"interactive", "cp -i file.txt /tmp/", "cp -i <path>+"},
		{"no-clobber", "cp -n file.txt /tmp/", "cp -n <path>+"},
		{"verbose recursive", "cp -rv dir1/ /dest/", "cp -rv <path>+"},
		{"preserve", "cp -p file.txt /tmp/", "cp -p <path>+"},
		{"bundled flags", "cp -rpf src/ /tmp/", "cp -rpf <path>+"},

		// Multiple source files
		{"three files", "cp a.go b.go /tmp/", "cp <path>+"},
		{"many files", "cp a.txt b.txt c.txt d.txt /dest/", "cp <path>+"},

		// Redirects
		{"with redirect", "cp file.txt /tmp/ 2>/dev/null", "cp <path>+ 2>/dev/null"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS: different file paths → same shape
	t.Run("different files collide", func(t *testing.T) {
		a := Normalize("cp /etc/hosts /tmp/hosts")
		b := Normalize("cp /var/log/syslog /backup/syslog")
		c := Normalize("cp ~/.bashrc /tmp/bashrc")
		if a != b || b != c {
			t.Errorf("expected all same: %q, %q, %q", a, b, c)
		}
	})

	t.Run("different flags same paths collide", func(t *testing.T) {
		a := Normalize("cp -r /home/alice/ /backup/")
		b := Normalize("cp -r /home/bob/ /archive/")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("cp literal-arg")
		subshell := Normalize("cp $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
