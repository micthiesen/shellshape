package shellshape

import "testing"

func TestPasswd(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"no args", "passwd", "passwd"},
		{"with username", "passwd alice", "passwd <user>"},
		{"with root user", "passwd root", "passwd <user>"},

		// Boolean flags
		{"status flag", "passwd -S alice", "passwd -S <user>"},
		{"status all", "passwd -S -a", "passwd -S -a"},
		{"lock user", "passwd -l bob", "passwd -l <user>"},
		{"unlock user", "passwd -u charlie", "passwd -u <user>"},
		{"delete password", "passwd -d testuser", "passwd -d <user>"},
		{"expire password", "passwd -e deploy", "passwd -e <user>"},
		{"quiet mode", "passwd -q alice", "passwd -q <user>"},
		{"keep tokens", "passwd -k", "passwd -k"},
		{"stdin flag", "passwd --stdin alice", "passwd --stdin <user>"},

		// Flags with numeric arguments
		{"maxdays", "passwd -x 90 alice", "passwd -x N <user>"},
		{"mindays", "passwd -n 7 bob", "passwd -n N <user>"},
		{"warndays", "passwd -w 14 charlie", "passwd -w N <user>"},
		{"inactive", "passwd -i 30 deploy", "passwd -i N <user>"},
		{"long maxdays", "passwd --maxdays 90 alice", "passwd --maxdays N <user>"},
		{"long mindays", "passwd --mindays 7 bob", "passwd --mindays N <user>"},

		// Flags with path arguments
		{"chroot dir", "passwd -R /mnt/sysimage alice", "passwd -R <path> <user>"},
		{"prefix dir", "passwd -P /srv/chroot bob", "passwd -P <path> <user>"},
		{"long root", "passwd --root /mnt/sysimage alice", "passwd --root <path> <user>"},

		// Repository flag (kept verbatim)
		{"repository", "passwd -r ldap alice", "passwd -r ldap <user>"},
		{"long repository", "passwd --repository nis bob", "passwd --repository nis <user>"},

		// Combined flags
		{"multiple numeric", "passwd -x 90 -n 7 -w 14 bob", "passwd -x N -n N -w N <user>"},
		{"lock with quiet", "passwd -q -l alice", "passwd -q -l <user>"},
		{"complex combo", "passwd -x 365 -w 30 -i 10 --root /mnt/alt deploy", "passwd -x N -w N -i N --root <path> <user>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TEST: different usernames collapse to same shape
	t.Run("different users collide", func(t *testing.T) {
		a := Normalize("passwd -l alice")
		b := Normalize("passwd -l bob")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// COLLISION TEST: different numeric values collapse
	t.Run("different days collide", func(t *testing.T) {
		a := Normalize("passwd -x 90 alice")
		b := Normalize("passwd -x 365 bob")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("passwd alice")
		subshell := Normalize("passwd $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
