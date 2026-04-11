package shellshape

import "testing"

func TestWatch(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple command", "watch ls", "watch ls"},
		{"command with args", "watch ls -l", "watch ls -l"},
		{"command with path", "watch ls -l /tmp/foo", "watch ls -l <path>"},

		// Interval flag
		{"-n flag", "watch -n 5 ls", "watch -n N ls"},
		{"--interval flag", "watch --interval 2 df -h", "watch --interval N df -h"},
		{"-n decimal", "watch -n 0.5 date", "watch -n N date"},

		// Boolean flags
		{"-d flag", "watch -d ls -l", "watch -d ls -l"},
		{"-t flag", "watch -t ls", "watch -t ls"},
		{"-b flag", "watch -b make test", "watch -b make test"},
		{"-e flag", "watch -e make build", "watch -e make build"},
		{"-p flag", "watch -p -n 10 ls", "watch -p -n N ls"},
		{"-x flag", "watch -x ls -l /tmp", "watch -x ls -l <path>"},

		// Multiple flags
		{"multiple flags", "watch -d -n 1 df -h", "watch -d -n N df -h"},
		{"multiple flags and path", "watch -n 10 -p -d cat /etc/hosts", "watch -n N -p -d cat <path>"},
		{"-t -n combined", "watch -t -n 0.5 date", "watch -t -n N date"},

		// Equexit flag with argument
		{"-q flag", "watch -q 5 ls", "watch -q N ls"},
		{"--equexit flag", "watch --equexit 3 ls", "watch --equexit N ls"},

		// Shotsdir flag with path
		{"-s flag", "watch -s /tmp/shots ls", "watch -s <path> ls"},
		{"--shotsdir flag", "watch --shotsdir /tmp/shots ls", "watch --shotsdir <path> ls"},

		// Command with sub-command paths
		{"tail -f with path", "watch -n 2 tail -f /var/log/syslog", "watch -n N tail -f <path>"},
		{"grep with path", "watch -n 5 grep error /var/log/app.log", "watch -n N grep error <path>"},

		// Redirect
		{"with redirect", "watch ls > /tmp/out.txt", "watch ls > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision tests
	t.Run("different intervals collide", func(t *testing.T) {
		a := Normalize("watch -n 1 ls")
		b := Normalize("watch -n 60 ls")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different paths collide", func(t *testing.T) {
		a := Normalize("watch -n 5 cat /etc/hosts")
		b := Normalize("watch -n 5 cat /var/log/syslog")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("watch ls")
		subshell := Normalize("watch $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
