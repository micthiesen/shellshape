package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestTop(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare", "top", "top"},
		{"batch mode", "top -b", "top -b"},
		{"show threads", "top -H", "top -H"},
		{"idle toggle", "top -i", "top -i"},

		// Flags with numeric arguments
		{"iterations -n", "top -n 5", "top -n N"},
		{"delay -d", "top -d 2", "top -d N"},
		{"width -w", "top -w 120", "top -w N"},
		{"macos samples -l", "top -l 10", "top -l N"},
		{"macos delay -s", "top -s 3", "top -s N"},

		// Flags with data arguments
		{"user -u", "top -u root", "top -u <user>"},
		{"user -U", "top -U michael", "top -U <user>"},
		{"pid -p", "top -p 1234", "top -p <pid>"},
		{"sort -o", "top -o %CPU", "top -o <field>"},

		// Combinations
		{"batch with iterations", "top -b -n 10", "top -b -n N"},
		{"batch user iterations", "top -b -u root -n 5", "top -b -u <user> -n N"},
		{"pid and delay", "top -p 42 -d 1", "top -p <pid> -d N"},
		{"sort and user", "top -o %MEM -u nobody", "top -o <field> -u <user>"},

		// Bundled flags with consuming last char
		{"bundled -Hp", "top -Hp 5678", "top -Hp <pid>"},
		{"bundled -bn", "top -bn 3", "top -bn N"},

		// Consuming flag at end of input (no value follows)
		{"user trailing", "top -u", "top -u"},
		{"pid trailing", "top -p", "top -p"},
		{"sort trailing", "top -o", "top -o"},

		// Numeric flag at end of input (no value follows)
		{"iterations trailing", "top -n", "top -n"},
		{"delay trailing", "top -d", "top -d"},
		{"width trailing", "top -w", "top -w"},

		// Bundled flags at end of input (no value follows)
		{"bundled Hp trailing", "top -Hp", "top -Hp"},
		{"bundled bn trailing", "top -bn", "top -bn"},

		// Consuming flag with subshell value
		{"user subshell", "top -u $(whoami)", "top -u $(whoami)"},
		{"pid subshell", "top -p $(pgrep nginx)", "top -p $(pgrep nginx)"},
		{"sort subshell", "top -o $(echo CPU)", "top -o $(echo <str>)"},

		// Numeric flag with subshell value
		{"iterations subshell", "top -n $(echo 5)", "top -n $(echo <str>)"},
		{"delay subshell", "top -d $(calc)", "top -d $(calc)"},

		// Bundled flag with subshell value
		{"bundled Hp subshell", "top -Hp $(pgrep x)", "top -Hp $(pgrep x)"},
		{"bundled bn subshell", "top -bn $(echo 3)", "top -bn $(echo <str>)"},

		// Subshell as standalone token
		{"subshell standalone", "top $(flags)", "top $(flags)"},

		// Redirect
		{"redirect output", "top -b -n 1 > output.txt", "top -b -n N > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS: different data values → same shape
	t.Run("different users collide", func(t *testing.T) {
		a := shellshape.Normalize("top -u root")
		b := shellshape.Normalize("top -u nobody")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different pids collide", func(t *testing.T) {
		a := shellshape.Normalize("top -p 1234")
		b := shellshape.Normalize("top -p 9876")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different sort fields collide", func(t *testing.T) {
		a := shellshape.Normalize("top -o %CPU")
		b := shellshape.Normalize("top -o %MEM")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("top -p 1234")
		subshell := shellshape.Normalize("top -p $(pgrep nginx)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
