package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestScp(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"local to remote", "scp file.txt user@host:/tmp/", "scp <path> <remote>"},
		{"remote to local", "scp user@host:/tmp/file.txt ./dest/", "scp <remote> <path>"},
		{"remote shorthand", "scp host:file.txt /tmp/", "scp <remote> <path>"},
		{"recursive copy", "scp -r /local/dir user@host:/remote/dir", "scp -r <path> <remote>"},
		{"multiple sources", "scp file1.txt file2.txt user@host:/tmp/", "scp <path>+ <remote>"},

		// Flags with arguments
		{"port flag", "scp -P 2222 file.txt user@host:/tmp/", "scp -P N <path> <remote>"},
		{"identity flag", "scp -i ~/.ssh/id_rsa file.txt user@host:/tmp/", "scp -i <path>+ <remote>"},
		{"cipher flag", "scp -c aes256-ctr file.txt user@host:/tmp/", "scp -c <val> <path> <remote>"},
		{"ssh option", "scp -o StrictHostKeyChecking=no file.txt user@host:", "scp -o <val> <path> <remote>"},
		{"jump host", "scp -J jumphost file.txt user@host:/tmp/", "scp -J <host> <path> <remote>"},
		{"bandwidth limit", "scp -l 1000 file.txt user@host:/tmp/", "scp -l N <path> <remote>"},
		{"config file", "scp -F /etc/ssh/config file.txt user@host:/tmp/", "scp -F <path>+ <remote>"},
		{"sftp server", "scp -D /usr/lib/sftp-server file.txt user@host:/tmp/", "scp -D <path>+ <remote>"},
		{"sftp program", "scp -S /usr/bin/ssh file.txt user@host:/tmp/", "scp -S <path>+ <remote>"},

		// Boolean flags
		{"compress and quiet", "scp -C -q user@host:/tmp/file.txt /local/", "scp -C -q <remote> <path>"},
		{"preserve and verbose", "scp -pv file.txt user@host:/tmp/", "scp -pv <path> <remote>"},

		// URI form
		{"scp uri", "scp scp://user@host/file.txt /tmp/", "scp <scp-uri> <path>"},

		// Edge cases
		{"bare remote colon", "scp user@host: /tmp/file.txt", "scp <remote> <path>"},
		{"redirect", "scp file.txt user@host:/tmp/ 2>/dev/null", "scp <path> <remote> 2>/dev/null"},
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
	t.Run("different hosts collide", func(t *testing.T) {
		a := shellshape.Normalize("scp file.txt alice@prod-server:/opt/app/")
		b := shellshape.Normalize("scp file.txt bob@staging-host:/var/data/")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different ports collide", func(t *testing.T) {
		a := shellshape.Normalize("scp -P 22 file.txt user@host:/tmp/")
		b := shellshape.Normalize("scp -P 8022 file.txt user@host:/tmp/")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("scp file.txt user@host:/tmp/")
		subshell := shellshape.Normalize("scp $(dangerous-command) user@host:/tmp/")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
