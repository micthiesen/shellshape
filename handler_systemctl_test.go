package shellshape

import "testing"

func TestSystemctl(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands with unit names
		{"start service", "systemctl start nginx.service", "systemctl start <unit>"},
		{"stop service", "systemctl stop sshd", "systemctl stop <unit>"},
		{"restart service", "systemctl restart postgresql", "systemctl restart <unit>"},
		{"reload service", "systemctl reload apache2.service", "systemctl reload <unit>"},
		{"status service", "systemctl status docker.service", "systemctl status <unit>"},
		{"enable service", "systemctl enable nginx.service", "systemctl enable <unit>"},
		{"disable service", "systemctl disable bluetooth.service", "systemctl disable <unit>"},
		{"is-active", "systemctl is-active ufw", "systemctl is-active <unit>"},
		{"is-enabled", "systemctl is-enabled sshd.service", "systemctl is-enabled <unit>"},
		{"is-failed", "systemctl is-failed foo.service", "systemctl is-failed <unit>"},

		// Multiple unit names collapse to single <unit>
		{"multiple units", "systemctl restart foo.service bar.service baz.service", "systemctl restart <unit>"},

		// No-argument subcommands
		{"daemon-reload", "systemctl daemon-reload", "systemctl daemon-reload"},
		{"daemon-reexec", "systemctl daemon-reexec", "systemctl daemon-reexec"},
		{"poweroff", "systemctl poweroff", "systemctl poweroff"},
		{"reboot", "systemctl reboot", "systemctl reboot"},
		{"suspend", "systemctl suspend", "systemctl suspend"},

		// Boolean flags
		{"enable with now", "systemctl enable --now docker.service", "systemctl enable --now <unit>"},
		{"user flag", "systemctl --user start my-app.service", "systemctl --user start <unit>"},
		{"failed flag", "systemctl --failed", "systemctl --failed"},
		{"list with all", "systemctl list-units --all", "systemctl list-units --all"},
		{"quiet flag", "systemctl -q is-active nginx", "systemctl -q is-active <unit>"},
		{"no-pager", "systemctl --no-pager status sshd", "systemctl --no-pager status <unit>"},

		// Flags with verbatim values (type, state, output)
		{"list-units type", "systemctl list-units --type service", "systemctl list-units --type service"},
		{"list-units short type", "systemctl list-units -t service", "systemctl list-units -t service"},
		{"list-units state", "systemctl list-units --state running", "systemctl list-units --state running"},
		{"list-units type and state", "systemctl list-units -t service --state running", "systemctl list-units -t service --state running"},
		{"list-unit-files type", "systemctl list-unit-files -a -t service", "systemctl list-unit-files -a -t service"},

		// Flags with collapsed values
		{"show property", "systemctl show -p ActiveState nginx", "systemctl show -p <val> <unit>"},
		{"show property long", "systemctl show --property ActiveState sshd", "systemctl show --property <val> <unit>"},
		{"show value flag", "systemctl show -p ActiveState --value ufw", "systemctl show -p <val> --value <unit>"},
		{"host flag", "systemctl -H root@server1 status nginx", "systemctl -H <val> status <unit>"},
		{"lines flag", "systemctl status -n 50 nginx", "systemctl status -n N <unit>"},
		{"root flag", "systemctl --root /mnt/sysroot enable sshd", "systemctl --root <path> enable <unit>"},

		// list-dependencies
		{"list-dependencies", "systemctl list-dependencies foo.service", "systemctl list-dependencies <unit>"},

		// edit
		{"edit service", "systemctl edit nginx.service", "systemctl edit <unit>"},

		// cat
		{"cat service", "systemctl cat sshd.service", "systemctl cat <unit>"},

		// mask/unmask
		{"mask", "systemctl mask cups.service", "systemctl mask <unit>"},
		{"unmask", "systemctl unmask cups.service", "systemctl unmask <unit>"},

		// isolate
		{"isolate target", "systemctl isolate multi-user.target", "systemctl isolate <unit>"},
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
	t.Run("different unit names collide", func(t *testing.T) {
		a := Normalize("systemctl restart nginx.service")
		b := Normalize("systemctl restart postgresql.service")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different property values collide", func(t *testing.T) {
		a := Normalize("systemctl show -p ActiveState nginx")
		b := Normalize("systemctl show -p SubState sshd")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different hosts collide", func(t *testing.T) {
		a := Normalize("systemctl -H root@server1 status nginx")
		b := Normalize("systemctl -H admin@server2 status sshd")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("systemctl start literal-arg")
		subshell := Normalize("systemctl start $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
