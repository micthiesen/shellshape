package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestKopia(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Snapshot operations
		{"snapshot create", "kopia snapshot create /home/user/documents", "kopia snapshot create <path>"},
		{"snapshot list", "kopia snapshot list", "kopia snapshot list"},
		{"snapshot restore", "kopia snapshot restore abc123def456 /tmp/restore", "kopia snapshot restore <val> <path>"},

		// Repository operations
		{"repo create filesystem", "kopia repository create filesystem --path /mnt/backup", "kopia repository create filesystem --path <path>"},
		{"repo connect s3", "kopia repository connect s3 --bucket my-backup-bucket --access-key AKIA1234 --secret-access-key secretkey123", "kopia repository connect s3 --bucket <val> --access-key <val> --secret-access-key <val>"},
		{"repo status", "kopia repository status", "kopia repository status"},

		// Policy operations
		{"policy set global", "kopia policy set --global --keep-latest 10 --compression zstd", "kopia policy set --global --keep-latest N --compression <val>"},
		{"policy set add-ignore", "kopia policy set --global --add-ignore node_modules", "kopia policy set --global --add-ignore <val>"},
		{"policy list", "kopia policy list", "kopia policy list"},

		// Server and mount
		{"server start", "kopia server start --address 0.0.0.0:51515", "kopia server start --address <val>"},
		{"mount all", "kopia mount <val> /mnt/kopia", "kopia mount <val> <path>"},

		// Multiple paths
		{"snapshot create multiple", "kopia snapshot create /home/user /etc /var/log", "kopia snapshot create <path>+"},
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
	t.Run("different paths collide", func(t *testing.T) {
		a := shellshape.Normalize("kopia snapshot create /home/alice")
		b := shellshape.Normalize("kopia snapshot create /home/bob")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different snapshot ids collide", func(t *testing.T) {
		a := shellshape.Normalize("kopia snapshot restore abc123def456 /tmp/a")
		b := shellshape.Normalize("kopia snapshot restore fff999aaa111 /tmp/b")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("kopia snapshot create /home/user")
		subshell := shellshape.Normalize("kopia snapshot create $(dangerous)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
