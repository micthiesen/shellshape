package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestDf(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"no args", "df", "df"},
		{"human readable", "df -h", "df -h"},
		{"single path", "df /home", "df <path>"},
		{"specific device", "df /dev/sda1", "df <path>"},
		{"human readable with path", "df -h /var/log", "df -h <path>"},

		// Boolean flags
		{"inode flag", "df -i", "df -i"},
		{"local only", "df -l", "df -l"},
		{"posix output", "df -P", "df -P"},
		{"si units", "df -H", "df -H"},
		{"combined flags", "df -ahY", "df -ahY"},
		{"all local inodes", "df -ali", "df -ali"},

		// Flags with arguments
		{"type flag short", "df -T nfs", "df -T <type>"},
		{"type flag with csv", "df -T nfs,tmpfs", "df -T <type>"},
		{"type with no prefix", "df -T nonfs,mfs", "df -T <type>"},
		{"legacy type flag", "df -t ext4", "df -t <type>"},
		{"exclude type linux", "df -x tmpfs", "df -x <type>"},
		{"exclude type long", "df --exclude-type tmpfs", "df --exclude-type <type>"},
		{"type long with equals", "df --type=ext4", "df --type=<val>"},

		// Multiple positionals
		{"multiple paths", "df /home /tmp /var", "df <path>+"},
		{"flags and multi paths", "df -h /home /tmp /var", "df -h <path>+"},

		// Edge cases
		{"current dir", "df .", "df ."},
		{"redirect", "df -h > output.txt", "df -h > <path>"},
		{"combined flags with type", "df -lT ext4", "df -lT <type>"},
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
	t.Run("different paths collide", func(t *testing.T) {
		a := shellshape.Normalize("df -h /home/alice")
		b := shellshape.Normalize("df -h /home/bob")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different types collide", func(t *testing.T) {
		a := shellshape.Normalize("df -T ext4")
		b := shellshape.Normalize("df -T nfs")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("df /home/user")
		subshell := shellshape.Normalize("df $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
