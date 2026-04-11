package shellshape

import "testing"

func TestCrontab(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"list", "crontab -l", "crontab -l"},
		{"edit", "crontab -e", "crontab -e"},
		{"remove", "crontab -r", "crontab -r"},
		{"interactive remove", "crontab -ri", "crontab -ri"},

		// Flag with argument
		{"user list", "crontab -u root -l", "crontab -u <user> -l"},
		{"user edit", "crontab -u admin -e", "crontab -u <user> -e"},
		{"user remove", "crontab -u deploy -r", "crontab -u <user> -r"},

		// Positional (file to install)
		{"install file", "crontab /etc/cron.d/myjob", "crontab <path>"},
		{"install relative", "crontab mycron.txt", "crontab <path>"},
		{"user install file", "crontab -u www-data /tmp/cron.txt", "crontab -u <user> <path>"},

		// With redirect
		{"list to file", "crontab -l > /tmp/backup.txt", "crontab -l > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS: different data values → same shape
	t.Run("different users collide", func(t *testing.T) {
		a := Normalize("crontab -u alice -l")
		b := Normalize("crontab -u bob -l")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different files collide", func(t *testing.T) {
		a := Normalize("crontab /etc/cron.d/job1")
		b := Normalize("crontab /tmp/mycron.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("crontab mycron.txt")
		subshell := Normalize("crontab $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
