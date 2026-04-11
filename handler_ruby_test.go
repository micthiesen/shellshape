package shellshape

import "testing"

func TestRuby(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"script only", "ruby script.rb", "ruby <script>"},
		{"script with path", "ruby lib/app.rb", "ruby <script>"},
		{"no args", "ruby", "ruby"},
		{"version", "ruby -v", "ruby -v"},
		{"syntax check", "ruby -c script.rb", "ruby -c <script>"},

		// Inline code with -e
		{"eval", "ruby -e 'puts 1'", "ruby -e <code>"},
		{"eval with args", "ruby -e 'puts ARGV' foo bar", "ruby -e <code> <arg>+"},
		{"multiple eval", "ruby -e 'x=1' -e 'puts x'", "ruby -e <code> -e <code>"},

		// One-liner flags (-n, -p, -l, -a combined with -e)
		{"n and e", "ruby -ne 'puts $_' file.txt", "ruby -ne <code> <path>"},
		{"p and e", "ruby -pe 'gsub(/foo/,\"bar\")' file.txt", "ruby -pe <code> <path>"},
		{"l and n and e", "ruby -lne 'puts $_' file.txt", "ruby -lne <code> <path>"},

		// Include path flags
		{"include short", "ruby -I lib script.rb", "ruby -I <path> <script>"},
		{"include fused", "ruby -Ilib script.rb", "ruby -I <path> <script>"},
		{"require short", "ruby -r json script.rb", "ruby -r <val> <script>"},
		{"require fused", "ruby -rjson script.rb", "ruby -r <val> <script>"},

		// Directory change
		{"chdir short", "ruby -C /tmp script.rb", "ruby -C <path> <script>"},
		{"chdir fused", "ruby -C/tmp script.rb", "ruby -C <path> <script>"},

		// Encoding flag
		{"encoding short", "ruby -E utf-8 script.rb", "ruby -E <val> <script>"},
		{"encoding long", "ruby --encoding=utf-8 script.rb", "ruby --encoding=<val> <script>"},

		// Split pattern
		{"field separator", "ruby -F: -ane 'puts $F[0]' file.txt", "ruby -F <val> -ane <code> <path>"},
		{"field separator fused", "ruby -F, -ane 'puts $F[0]' file.txt", "ruby -F <val> -ane <code> <path>"},

		// In-place edit
		{"in-place", "ruby -i -pe 'gsub(/a/,\"b\")' file.txt", "ruby -i -pe <code> <path>"},
		{"in-place with backup", "ruby -i.bak -pe 'gsub(/a/,\"b\")' file.txt", "ruby -i.bak -pe <code> <path>"},

		// Script with arguments
		{"script with args", "ruby app.rb --port 3000 foo", "ruby <script> <arg>+"},

		// Double dash
		{"double dash", "ruby script.rb -- --not-a-flag", "ruby <script> -- <arg>"},

		// Boolean flags
		{"debug", "ruby --debug script.rb", "ruby --debug <script>"},
		{"warnings", "ruby -w script.rb", "ruby -w <script>"},
		{"verbose", "ruby --verbose script.rb", "ruby --verbose <script>"},

		// JIT flags with = syntax
		{"jit-verbose", "ruby --jit-verbose=1 script.rb", "ruby --jit-verbose=<val> <script>"},
		{"jit-max-cache", "ruby --jit-max-cache=100 script.rb", "ruby --jit-max-cache=<val> <script>"},

		// Enable/disable
		{"enable", "ruby --enable=frozen-string-literal script.rb", "ruby --enable=<val> <script>"},
		{"disable", "ruby --disable=gems script.rb", "ruby --disable=<val> <script>"},

		// Complex real-world
		{"complex", "ruby -I lib -r helper -e 'puts 1'", "ruby -I <path> -r <val> -e <code>"},
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
	t.Run("different scripts collide", func(t *testing.T) {
		a := Normalize("ruby server.rb --port 3000")
		b := Normalize("ruby worker.rb --port 8080")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different eval code collides", func(t *testing.T) {
		a := Normalize("ruby -e 'puts 1'")
		b := Normalize("ruby -e 'puts 2'")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("ruby script.rb")
		subshell := Normalize("ruby $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
