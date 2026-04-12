package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestRclone(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands
		{"copy local to remote", "rclone copy /local/path remote:bucket/path", "rclone copy <path>+"},
		{"sync", "rclone sync /local/path gdrive:backup", "rclone sync <path>+"},
		{"move", "rclone move /tmp/files s3:mybucket/folder", "rclone move <path>+"},
		{"ls", "rclone ls remote:path", "rclone ls <path>"},
		{"lsd", "rclone lsd remote:", "rclone lsd <path>"},
		{"mkdir", "rclone mkdir remote:newdir", "rclone mkdir <path>"},
		{"delete", "rclone delete remote:path/to/file", "rclone delete <path>"},
		{"config", "rclone config", "rclone config"},
		{"mount", "rclone mount remote:path /mnt/remote", "rclone mount <path>+"},

		// Backend path syntax
		{"backend path", "rclone copy :s3:bucket/path /local", "rclone copy <path>+"},

		// Filter flags
		{"include filter", "rclone copy remote:src /local --include '*.txt'", "rclone copy <path>+ --include <pattern>"},
		{"exclude filter", "rclone sync /local remote:dst --exclude '*.log'", "rclone sync <path>+ --exclude <pattern>"},
		{"filter flag", "rclone ls remote: --filter '+ *.jpg'", "rclone ls <path> --filter <pattern>"},

		// Numeric flags
		{"transfers", "rclone copy /src remote:dst --transfers 8", "rclone copy <path>+ --transfers N"},
		{"checkers", "rclone sync /src remote:dst --checkers 16", "rclone sync <path>+ --checkers N"},

		// Value flags
		{"bwlimit", "rclone copy /src remote:dst --bwlimit 10M", "rclone copy <path>+ --bwlimit <val>"},

		// Boolean flags
		{"dry run", "rclone sync /src remote:dst --dry-run", "rclone sync <path>+ --dry-run"},
		{"verbose", "rclone copy /src remote:dst -v", "rclone copy <path>+ -v"},
		{"progress", "rclone copy /src remote:dst --progress", "rclone copy <path>+ --progress"},

		// Path flags
		{"log file", "rclone sync /src remote:dst --log-file /var/log/rclone.log", "rclone sync <path>+ --log-file <path>"},
		{"config file", "rclone ls remote: --config /home/user/.config/rclone/rclone.conf", "rclone ls <path> --config <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different remotes collide", func(t *testing.T) {
		a := shellshape.Normalize("rclone copy /local gdrive:photos")
		b := shellshape.Normalize("rclone copy /local s3:mybucket")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different filters collide", func(t *testing.T) {
		a := shellshape.Normalize("rclone ls remote: --include '*.jpg'")
		b := shellshape.Normalize("rclone ls remote: --include '*.png'")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("rclone copy /local remote:path")
		subshell := shellshape.Normalize("rclone copy $(dangerous) remote:path")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
