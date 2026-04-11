package shellshape

import "testing"

func TestTar(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"create gzipped", "tar -czf archive.tar.gz src/", "tar -czf <archive> <path>"},
		{"extract gzipped", "tar -xzf archive.tar.gz", "tar -xzf <archive>"},
		{"list archive", "tar -tf archive.tar", "tar -tf <archive>"},
		{"bundled no dash", "tar xzf archive.tar.gz", "tar xzf <archive>"},
		{"create bundled no dash", "tar czf backup.tar.gz /home/user", "tar czf <archive> <path>"},

		// Long form flags
		{"long create", "tar --create --gzip --file backup.tar.gz /home", "tar --create --gzip --file <archive> <path>"},
		{"long extract", "tar --extract --file archive.tar", "tar --extract --file <archive>"},

		// Flags with arguments
		{"change dir", "tar -xf archive.tar.gz -C /tmp/dest", "tar -xf <archive> -C <path>"},
		{"exclude pattern", "tar -czf archive.tar.gz --exclude '*.log' src/", "tar -czf <archive> --exclude <pattern> <path>"},
		{"strip components", "tar -xf archive.tar --strip-components 2", "tar -xf <archive> --strip-components N"},
		{"files-from", "tar -czf archive.tar.gz -T filelist.txt", "tar -czf <archive> -T <path>"},
		{"exclude-from", "tar -czf archive.tar.gz -X excludes.txt src/", "tar -czf <archive> -X <path>+"},

		// Multiple positionals
		{"multiple files", "tar -cf archive.tar file1.txt file2.txt file3.txt", "tar -cf <archive> <path>+"},
		{"extract specific files", "tar -xf archive.tar path/to/file1 path/to/file2", "tar -xf <archive> <path>+"},

		// Boolean flags preserved
		{"verbose", "tar -xvf archive.tar.gz", "tar -xvf <archive>"},
		{"exclude vcs", "tar -czf archive.tar.gz --exclude-vcs .", "tar -czf <archive> --exclude-vcs ."},

		// Stdout as archive
		{"to stdout", "tar -cf - dir/", "tar -cf <archive> <path>"},

		// Multiple excludes
		{"multiple excludes", "tar -czf archive.tar.gz --exclude '*.log' --exclude '*.tmp' src/", "tar -czf <archive> --exclude <pattern> --exclude <pattern> <path>"},
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
	t.Run("different archives collide", func(t *testing.T) {
		a := Normalize("tar -xzf project-v1.tar.gz")
		b := Normalize("tar -xzf backup-2024.tar.gz")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different directories collide", func(t *testing.T) {
		a := Normalize("tar -czf archive.tar.gz /home/alice/docs")
		b := Normalize("tar -czf archive.tar.gz /var/log/nginx")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("tar -xf archive.tar")
		subshell := Normalize("tar -xf $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
