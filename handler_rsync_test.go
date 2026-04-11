package shellshape

import "testing"

func TestRsync(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"local copy", "rsync src/ dest/", "rsync <path>+"},
		{"archive verbose compress", "rsync -avz /data/ /backup/", "rsync -avz <path>+"},
		{"remote destination", "rsync -avz src/ user@host:/remote/dir/", "rsync -avz <path>+"},
		{"remote source", "rsync -avz server:/var/log/ /tmp/logs/", "rsync -avz <path>+"},
		{"delete flag", "rsync -avz --delete /local/ remote:/backup/", "rsync -avz --delete <path>+"},
		{"progress", "rsync --progress -r /data/ /backup/", "rsync --progress -r <path>+"},

		// Flags with arguments
		{"rsh flag short", "rsync -avz -e ssh src/ dest/", "rsync -avz -e <rsh> <path>+"},
		{"rsh flag with options", "rsync -avz -e 'ssh -p 2222' src/ dest/", "rsync -avz -e <rsh> <path>+"},
		{"filter flag short", "rsync -avz -f 'merge .gitignore' src/ dest/", "rsync -avz -f <filter> <path>+"},
		{"exclude with equals", "rsync -avz --exclude=*.log src/ dest/", "rsync -avz --exclude=<val> <path>+"},
		{"exclude space-separated", "rsync -avz --exclude '*.log' src/ dest/", "rsync -avz --exclude <filter> <path>+"},
		{"include space-separated", "rsync -avz --include '*.txt' --exclude '*' src/ dest/", "rsync -avz --include <filter> --exclude <filter> <path>+"},
		{"bwlimit equals", "rsync -avz --bwlimit=1000 src/ dest/", "rsync -avz --bwlimit=<val> <path>+"},
		{"multiple sources", "rsync -avz file1.txt file2.txt dest/", "rsync -avz <path>+"},

		// Edge cases
		{"no flags", "rsync file.txt /backup/", "rsync <path>+"},
		{"dry run", "rsync -avzn --delete src/ dest/", "rsync -avzn --delete <path>+"},
		{"rsync protocol url", "rsync -avz rsync://mirror.example.com/pub/ /local/", "rsync -avz <rsync-uri> <path>"},
		{"single path", "rsync --list-only server:/var/log/", "rsync --list-only <path>"},
		{"exclude-from flag", "rsync -avz --exclude-from /home/user/.rsync-exclude src/ dest/", "rsync -avz --exclude-from <path>+"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different paths collide", func(t *testing.T) {
		a := Normalize("rsync -avz /home/alice/docs/ /mnt/backup/alice/")
		b := Normalize("rsync -avz /home/bob/photos/ /mnt/backup/bob/")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different exclude patterns collide", func(t *testing.T) {
		a := Normalize("rsync -avz --exclude '*.log' src/ dest/")
		b := Normalize("rsync -avz --exclude '*.tmp' src/ dest/")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("rsync -avz /data/ /backup/")
		subshell := Normalize("rsync -avz $(dangerous-command) /backup/")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
