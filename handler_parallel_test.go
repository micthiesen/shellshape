package shellshape

import "testing"

func TestParallel(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple with input args", "parallel gzip ::: file1 file2 file3", "parallel gzip ::: <arg>+"},
		{"command only", "parallel gzip", "parallel gzip"},
		{"jobs flag short", "parallel -j4 gzip ::: a b", "parallel -j N gzip ::: <arg>+"},
		{"jobs flag long", "parallel --jobs 8 gzip ::: a b c", "parallel --jobs N gzip ::: <arg>+"},

		// Flags with arguments
		{"sshlogin", "parallel -S host1,host2 uptime ::: a b", "parallel -S <str> uptime ::: <arg>+"},
		{"colsep", "parallel --colsep '\\t' echo {1} {2} ::: 'a\\tb' 'c\\td'", "parallel --colsep <str> echo {1} {2} ::: <arg>+"},
		{"arg-file", "parallel -a input.txt gzip", "parallel -a <path> gzip"},
		{"block flag", "parallel --pipe --block 1M sort", "parallel --pipe --block <str> sort"},
		{"max-args", "parallel -n 2 echo ::: a b c d", "parallel -n N echo ::: <arg>+"},
		{"max-lines", "parallel -L 1 wc -l ::: a b", "parallel -L N wc -l ::: <arg>+"},
		{"workdir", "parallel --workdir /tmp echo ::: a", "parallel --workdir <path> echo ::: <arg>"},
		{"results", "parallel --results outdir echo ::: a b", "parallel --results <str> echo ::: <arg>+"},
		{"joblog", "parallel --joblog log.txt echo ::: a b", "parallel --joblog <path> echo ::: <arg>+"},
		{"timeout", "parallel --timeout 60 echo ::: a b", "parallel --timeout N echo ::: <arg>+"},
		{"retries", "parallel --retries 3 echo ::: a", "parallel --retries N echo ::: <arg>"},
		{"delay", "parallel --delay 0.5 echo ::: a b", "parallel --delay N echo ::: <arg>+"},
		{"tagstring", "parallel --tagstring '{#}' echo ::: a b", "parallel --tagstring <str> echo ::: <arg>+"},

		// Boolean flags
		{"keep-order", "parallel -k echo ::: a b c", "parallel -k echo ::: <arg>+"},
		{"pipe", "parallel --pipe sort", "parallel --pipe sort"},
		{"dry-run", "parallel --dry-run echo ::: a b", "parallel --dry-run echo ::: <arg>+"},
		{"progress", "parallel --progress echo ::: a b", "parallel --progress echo ::: <arg>+"},
		{"null input", "parallel -0 echo ::: a b", "parallel -0 echo ::: <arg>+"},

		// Multiple ::: groups
		{"two arg groups", "parallel echo ::: A B C ::: 1 2 3", "parallel echo ::: <arg>+ ::: <arg>+"},
		{"three arg groups", "parallel echo ::: a ::: b ::: c", "parallel echo ::: <arg> ::: <arg> ::: <arg>"},

		// :::: (arg file source)
		{"argfile source", "parallel echo :::: input.txt", "parallel echo :::: <path>"},

		// Multi-token command template
		{"multi-token cmd", "parallel convert {} {.}.png ::: img1.jpg img2.jpg", "parallel convert {} {.}.png ::: <arg>+"},

		// Redirect
		{"redirect", "parallel echo ::: a b > out.txt", "parallel echo ::: <arg>+ > <path>"},
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
	t.Run("different args collide", func(t *testing.T) {
		a := Normalize("parallel gzip ::: fileA.dat fileB.dat")
		b := Normalize("parallel gzip ::: x.log y.log z.log")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different hosts collide", func(t *testing.T) {
		a := Normalize("parallel -S host1,host2 uptime ::: a")
		b := Normalize("parallel -S server3,server4 uptime ::: a")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("parallel echo ::: literal-arg")
		subshell := Normalize("parallel echo ::: $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
