package shellshape

import "testing"

func TestJournalctl(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"follow", "journalctl -f", "journalctl -f"},
		{"unit filter", "journalctl -u nginx", "journalctl -u <val>"},
		{"boot errors", "journalctl -b -p err", "journalctl -b -p <val>"},
		{"since until", "journalctl --since '2024-01-01' --until '2024-01-02 23:59:59'", "journalctl --since <val> --until <val>"},
		{"output format", "journalctl -o json", "journalctl -o <val>"},
		{"lines limit", "journalctl -n 100", "journalctl -n N"},
		{"identifier", "journalctl -t systemd-resolved", "journalctl -t <val>"},
		{"grep", "journalctl -g 'error.*timeout'", "journalctl -g <pattern>"},
		{"kernel messages", "journalctl -k", "journalctl -k"},
		{"reverse with unit", "journalctl -r -u sshd", "journalctl -r -u <val>"},

		// Flags with arguments
		{"directory", "journalctl -D /var/log/journal/remote", "journalctl -D <path>"},
		{"file flag", "journalctl --file /run/log/journal/abc/system.journal", "journalctl --file <path>"},
		{"cursor", "journalctl --cursor s=abc123", "journalctl --cursor <val>"},
		{"boot offset", "journalctl -b -1", "journalctl -b <val>"},
		{"vacuum time", "journalctl --vacuum-time=2d", "journalctl --vacuum-time=<val>"},
		{"vacuum size", "journalctl --vacuum-size=500M", "journalctl --vacuum-size=<val>"},
		{"priority number", "journalctl -p 3", "journalctl -p <val>"},
		{"facility", "journalctl --facility=kern", "journalctl --facility=<val>"},

		// Positionals: field matches and executable paths
		{"pid match", "journalctl _PID=123", "journalctl <field-match>"},
		{"comm match", "journalctl _COMM=sshd", "journalctl <field-match>"},
		{"uid match", "journalctl _UID=1000", "journalctl <field-match>"},
		{"exe path", "journalctl /usr/bin/dbus-daemon", "journalctl <path>"},
		{"field match with flags", "journalctl _COMM=crond --since '10:00' --until '11:00'", "journalctl <field-match> --since <val> --until <val>"},

		// Multiple positionals
		{"multiple field matches", "journalctl _COMM=sshd _PID=456", "journalctl <field-match>+"},

		// Boolean flags
		{"list boots", "journalctl --list-boots", "journalctl --list-boots"},
		{"disk usage", "journalctl --disk-usage", "journalctl --disk-usage"},
		{"no pager", "journalctl --no-pager -u docker", "journalctl --no-pager -u <val>"},

		// Long flag with =
		{"output equals", "journalctl --output=json-pretty", "journalctl --output=<val>"},
		{"unit equals", "journalctl --unit=sshd", "journalctl --unit=<val>"},
		{"priority equals", "journalctl --priority=warning", "journalctl --priority=<val>"},

		// Combined
		{"follow unit no-pager", "journalctl -f -u nginx --no-pager", "journalctl -f -u <val> --no-pager"},
		{"complex", "journalctl -b -p err -u sshd -o json --no-pager -n 50", "journalctl -b -p <val> -u <val> -o <val> --no-pager -n N"},
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
	t.Run("different units collide", func(t *testing.T) {
		a := Normalize("journalctl -u nginx")
		b := Normalize("journalctl -u sshd")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different field matches collide", func(t *testing.T) {
		a := Normalize("journalctl _PID=123")
		b := Normalize("journalctl _COMM=sshd")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different timestamps collide", func(t *testing.T) {
		a := Normalize("journalctl --since '2024-01-01' --until '2024-01-02'")
		b := Normalize("journalctl --since yesterday --until today")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("journalctl -u nginx")
		subshell := Normalize("journalctl -u $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
