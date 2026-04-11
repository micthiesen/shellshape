package shellshape

import "testing"

func TestDiff(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"two files", "diff file1.txt file2.txt", "diff <path>+"},
		{"unified", "diff -u file1.txt file2.txt", "diff -u <path>+"},
		{"recursive dirs", "diff -ruN dir1/ dir2/", "diff -ruN <path>+"},
		{"context format", "diff -c old.txt new.txt", "diff -c <path>+"},
		{"brief", "diff -q dir1/ dir2/", "diff -q <path>+"},
		{"side by side", "diff -y left.txt right.txt", "diff -y <path>+"},

		// Flags with arguments
		{"-C numeric", "diff -C 3 file1.txt file2.txt", "diff -C N <path>+"},
		{"-U numeric", "diff -U 5 old.txt new.txt", "diff -U N <path>+"},
		{"--context long", "diff --context 10 a.txt b.txt", "diff --context N <path>+"},
		{"--unified long", "diff --unified 8 a.txt b.txt", "diff --unified N <path>+"},
		{"-I pattern", "diff -I '^#' file1.go file2.go", "diff -I <pattern> <path>+"},
		{"-F pattern", "diff -F '^func' file1.go file2.go", "diff -F <pattern> <path>+"},
		{"-L label", "diff -L original -L modified file1.txt file2.txt", "diff -L <str> -L <str> <path>+"},
		{"-D ifdef", "diff -D NEWVERSION old.c new.c", "diff -D <str> <path>+"},
		{"-x exclude", "diff -x '*.o' dir1/ dir2/", "diff -x <pattern> <path>+"},
		{"-X exclude-from", "diff -X ./exclude.list dir1/ dir2/", "diff -X <path>+"},
		{"-S starting file", "diff -S ./start.txt dir1/ dir2/", "diff -S <path>+"},
		{"--tabsize", "diff --tabsize 4 a.txt b.txt", "diff --tabsize N <path>+"},
		{"-W width", "diff -W 120 -y a.txt b.txt", "diff -W N -y <path>+"},
		{"--color=always", "diff --color=always a.txt b.txt", "diff --color=<val> <path>+"},
		{"-A algorithm", "diff -A patience a.txt b.txt", "diff -A patience <path>+"},
		{"--algorithm", "diff --algorithm myers a.txt b.txt", "diff --algorithm myers <path>+"},
		{"--changed-group-format", "diff --changed-group-format '%>' a.txt b.txt", "diff --changed-group-format <str> <path>+"},

		// Fused numeric
		{"-C fused", "diff -C3 old.txt new.txt", "diff -C N <path>+"},
		{"-U fused", "diff -U5 old.txt new.txt", "diff -U N <path>+"},

		// Edge cases
		{"stdin dash", "diff - file.txt", "diff - <path>"},
		{"absolute paths", "diff /etc/hosts /tmp/hosts", "diff <path>+"},
		{"redirect", "diff a.txt b.txt > output.txt", "diff <path>+ > <path>"},
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
	t.Run("different files collide", func(t *testing.T) {
		a := Normalize("diff -u config.yaml settings.yaml")
		b := Normalize("diff -u main.go util.go")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different context numbers collide", func(t *testing.T) {
		a := Normalize("diff -C 3 file1.txt file2.txt")
		b := Normalize("diff -C 10 old.txt new.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("diff file1.txt file2.txt")
		subshell := Normalize("diff $(dangerous-command) file2.txt")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
