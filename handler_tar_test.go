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

		// Value flags (--format, --owner, --group, etc.)
		{"format flag", "tar -czf archive.tar.gz --format pax src/", "tar -czf <archive> --format <val> <path>"},
		{"owner flag", "tar -czf archive.tar.gz --owner root src/", "tar -czf <archive> --owner <val> <path>"},
		{"group flag", "tar -czf archive.tar.gz --group staff src/", "tar -czf <archive> --group <val> <path>"},
		{"val flag at end", "tar -czf archive.tar.gz --format", "tar -czf <archive> --format"},

		// Subshell as archive after -f/--file
		{"subshell after -f", "tar -xf $(get-archive)", "tar -xf $(get-archive)"},
		{"subshell after --file", "tar --extract --file $(get-archive)", "tar --extract --file $(get-archive)"},

		// Subshell after path flag
		{"subshell after -C", "tar -xf archive.tar -C $(get-dir)", "tar -xf <archive> -C $(get-dir)"},
		{"subshell after --directory", "tar -xf archive.tar --directory $(get-dir)", "tar -xf <archive> --directory $(get-dir)"},

		// Subshell after pattern flag
		{"subshell after --exclude", "tar -czf archive.tar.gz --exclude $(get-pattern) src/", "tar -czf <archive> --exclude $(get-pattern) <path>"},
		{"subshell after --include", "tar -czf archive.tar.gz --include $(get-pattern) src/", "tar -czf <archive> --include $(get-pattern) <path>"},

		// Subshell as positional token in main loop
		{"subshell positional", "tar -xf archive.tar $(get-files)", "tar -xf <archive> $(get-files)"},

		// Subshell after bundle f
		{"subshell after bundle f", "tar czf $(gen-name) src/", "tar czf $(gen-name) <path>"},

		// Non-bundle first tokens (exercise isTarBundle false paths)
		{"long flag first", "tar --verbose -xf archive.tar", "tar --verbose -xf <archive>"},
		{"single char first", "tar v", "tar v"},
		{"digits in bundle", "tar x2f archive.tar", "tar x2f <archive>"},
		{"no mode letter bundle", "tar vz archive.tar", "tar vz <path>"},

		// Block size numeric flag
		{"block-size flag", "tar -xf archive.tar --block-size 512", "tar -xf <archive> --block-size N"},
		{"b flag", "tar -xf archive.tar -b 20", "tar -xf <archive> -b N"},

		// Redirects
		{"redirect output", "tar -cf - src/ > backup.tar", "tar -cf <archive> <path> > <path>"},
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
