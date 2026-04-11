package shellshape

import "testing"

func TestHtpasswd(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"add user", "htpasswd /etc/htpasswd admin", "htpasswd <path> <user>"},
		{"create file", "htpasswd -c /etc/htpasswd admin", "htpasswd -c <path> <user>"},
		{"delete user", "htpasswd -D /etc/htpasswd olduser", "htpasswd -D <path> <user>"},
		{"verify user", "htpasswd -v /etc/htpasswd admin", "htpasswd -v <path> <user>"},

		// Batch mode with password
		{"batch mode", "htpasswd -b /etc/htpasswd admin secret123", "htpasswd -b <path> <user> <str>"},
		{"create batch", "htpasswd -cb /etc/htpasswd admin mypassword", "htpasswd -cb <path> <user> <str>"},

		// Display to stdout (-n, no file)
		{"stdout mode", "htpasswd -n admin", "htpasswd -n <user>"},
		{"stdout batch", "htpasswd -nb admin secret123", "htpasswd -nb <user> <str>"},
		{"stdout md5 batch", "htpasswd -nbm admin secret123", "htpasswd -nbm <user> <str>"},

		// Bcrypt cost flag
		{"bcrypt cost", "htpasswd -B -C 12 /etc/htpasswd admin", "htpasswd -B -C N <path> <user>"},
		{"bcrypt cost stdout", "htpasswd -nB -C 10 admin", "htpasswd -nB -C N <user>"},

		// Algorithm flags
		{"md5", "htpasswd -m /etc/htpasswd admin", "htpasswd -m <path> <user>"},
		{"bcrypt", "htpasswd -B /etc/htpasswd admin", "htpasswd -B <path> <user>"},
		{"sha", "htpasswd -s /etc/htpasswd admin", "htpasswd -s <path> <user>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different users collide", func(t *testing.T) {
		a := Normalize("htpasswd /etc/htpasswd alice")
		b := Normalize("htpasswd /etc/htpasswd bob")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different files collide", func(t *testing.T) {
		a := Normalize("htpasswd -c /var/www/.htpasswd admin")
		b := Normalize("htpasswd -c /etc/apache2/.htpasswd admin")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("htpasswd /etc/htpasswd admin")
		subshell := Normalize("htpasswd $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
