package shellshape

import "testing"

func TestXargs(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare xargs", "xargs", "xargs"},
		{"utility only", "xargs rm", "xargs rm"},
		{"utility with flags", "xargs rm -rf", "xargs rm -rf"},
		{"utility with args", "xargs grep foo", "xargs grep <arg>"},

		// Flags with arguments
		{"-n numeric", "xargs -n 1 rm", "xargs -n N rm"},
		{"-P numeric", "xargs -P 4 curl", "xargs -P N curl"},
		{"-L numeric", "xargs -L 1 wc -l", "xargs -L N wc -l"},
		{"-s numeric", "xargs -s 1024 echo", "xargs -s N echo"},
		{"-I replstr", "xargs -I {} mv {} /dest", "xargs -I <str> mv <arg>+"},
		{"-J replstr", "xargs -J % cp -Rp % destdir", "xargs -J <str> cp -Rp <arg>+"},
		{"-E eofstr", "xargs -E EOF cat", "xargs -E <str> cat"},
		{"-R numeric", "xargs -I {} -R 3 echo {}", "xargs -I <str> -R N echo <arg>"},
		{"-S numeric", "xargs -I {} -S 512 echo {}", "xargs -I <str> -S N echo <arg>"},

		// Boolean flags
		{"-0 null separator", "xargs -0 rm", "xargs -0 rm"},
		{"-t verbose", "xargs -t rm -rf", "xargs -t rm -rf"},
		{"-p interactive", "xargs -p rm", "xargs -p rm"},
		{"-r no-run-if-empty", "xargs -r echo", "xargs -r echo"},

		// Combined flags
		{"-0 -n 1", "xargs -0 -n 1 rm", "xargs -0 -n N rm"},
		{"-P and -n", "xargs -P 4 -n 1 curl https://example.com", "xargs -P N -n N curl <arg>"},

		// Long flags
		{"--max-args", "xargs --max-args=5 rm", "xargs --max-args=<val> rm"},
		{"--null", "xargs --null rm", "xargs --null rm"},

		// Utility args classified
		{"utility with path arg", "xargs cat /tmp/foo.txt", "xargs cat <arg>"},

		// Redirect
		{"with redirect", "xargs echo > /tmp/out.txt", "xargs echo > <path>"},
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
	t.Run("different utility args collide", func(t *testing.T) {
		a := Normalize("xargs -I {} cp {} /backup/dir1")
		b := Normalize("xargs -I {} cp {} /backup/dir2")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different numeric args collide", func(t *testing.T) {
		a := Normalize("xargs -n 1 rm")
		b := Normalize("xargs -n 10 rm")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("xargs echo hello")
		subshell := Normalize("xargs echo $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
