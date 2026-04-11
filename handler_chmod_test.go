package shellshape

import "testing"

func TestChmod(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage - octal modes
		{"octal simple", "chmod 755 script.sh", "chmod 755 <path>"},
		{"octal 644", "chmod 644 config.yml", "chmod 644 <path>"},
		{"octal 4-digit", "chmod 4755 /usr/bin/foo", "chmod 4755 <path>"},

		// Symbolic modes
		{"symbolic +x", "chmod +x deploy.sh", "chmod +x <path>"},
		{"symbolic u+rw", "chmod u+rw file.txt", "chmod u+rw <path>"},
		{"symbolic compound", "chmod u+rw,go-w file.txt", "chmod u+rw,go-w <path>"},
		{"symbolic go=", "chmod go= secret.txt", "chmod go= <path>"},
		{"symbolic a=rx", "chmod a=rx /usr/local/bin/tool", "chmod a=rx <path>"},

		// Flags
		{"recursive", "chmod -R 755 /var/www", "chmod -R 755 <path>"},
		{"verbose recursive", "chmod -Rv 755 src/", "chmod -Rv 755 <path>"},
		{"force", "chmod -f 600 ~/.ssh/key.pem", "chmod -f 600 <path>"},
		{"long recursive", "chmod --recursive a+r docs/", "chmod --recursive a+r <path>"},

		// Multiple files
		{"multiple files", "chmod 644 a.txt b.txt c.txt", "chmod 644 <path>+"},
		{"multiple with flag", "chmod -R 755 src/ lib/ bin/", "chmod -R 755 <path>+"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TEST: different file paths collapse to same shape
	t.Run("different paths collide", func(t *testing.T) {
		a := Normalize("chmod 755 /home/alice/script.sh")
		b := Normalize("chmod 755 /opt/deploy/run.sh")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// COLLISION TEST: different modes stay distinct
	t.Run("different modes distinct", func(t *testing.T) {
		a := Normalize("chmod 755 file.sh")
		b := Normalize("chmod 644 file.sh")
		if a == b {
			t.Errorf("different modes should produce different shapes, both got %q", a)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("chmod 755 script.sh")
		subshell := Normalize("chmod $(dangerous-command) script.sh")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
