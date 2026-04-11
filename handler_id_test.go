package shellshape

import "testing"

func TestId(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare id", "id", "id"},
		{"user flag", "id -u", "id -u"},
		{"group flag", "id -g", "id -g"},
		{"groups flag", "id -G", "id -G"},
		{"real flag", "id -u -r", "id -u -r"},
		{"name flag", "id -gn", "id -gn"},
		{"combined flags", "id -Gn", "id -Gn"},
		{"human readable", "id -p", "id -p"},

		// With username positional
		{"username positional", "id bob", "id <user>"},
		{"username with flag", "id -u root", "id -u <user>"},
		{"username with combined flags", "id -Gn alice", "id -Gn <user>"},
		{"multiple flags then user", "id -u -r bob", "id -u -r <user>"},

		// Numeric UID as positional
		{"numeric uid", "id 501", "id <user>"},
		{"numeric uid with flag", "id -u 0", "id -u <user>"},

		// With redirect
		{"with redirect", "id -u > /tmp/out", "id -u > <path>"},

		// In pipeline
		{"in pipeline", "id -u && whoami", "id -u && whoami"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS: different data values → same shape
	t.Run("different users collide", func(t *testing.T) {
		a := Normalize("id alice")
		b := Normalize("id bob")
		c := Normalize("id root")
		if a != b || b != c {
			t.Errorf("expected all same: %q, %q, %q", a, b, c)
		}
	})

	t.Run("different uids collide", func(t *testing.T) {
		a := Normalize("id 0")
		b := Normalize("id 501")
		c := Normalize("id 1000")
		if a != b || b != c {
			t.Errorf("expected all same: %q, %q, %q", a, b, c)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("id bob")
		subshell := Normalize("id $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
