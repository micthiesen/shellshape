package shellshape

import "testing"

func TestDmesg(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"bare", "dmesg", "dmesg"},
		{"human readable", "dmesg -H", "dmesg -H"},
		{"timestamps", "dmesg -T", "dmesg -T"},
		{"follow", "dmesg -w", "dmesg -w"},
		{"combined boolean flags", "dmesg -TH", "dmesg -TH"},
		{"clear", "dmesg -C", "dmesg -C"},
		{"read-clear", "dmesg -c", "dmesg -c"},

		// Flags with arguments
		{"level", "dmesg --level err", "dmesg --level <level>"},
		{"level short", "dmesg -l warn,err", "dmesg -l <level>"},
		{"facility", "dmesg -f kern,daemon", "dmesg -f <facility>"},
		{"facility long", "dmesg --facility user", "dmesg --facility <facility>"},
		{"buffer size", "dmesg -s 16384", "dmesg -s N"},
		{"buffer size long", "dmesg --buffer-size 32768", "dmesg --buffer-size N"},
		{"console level", "dmesg -n 1", "dmesg -n <level>"},
		{"console level long", "dmesg --console-level warn", "dmesg --console-level <level>"},
		{"file", "dmesg -F /var/log/dmesg.log", "dmesg -F <path>"},
		{"file long", "dmesg --file /tmp/kern.log", "dmesg --file <path>"},
		{"kmsg file", "dmesg -K /dev/kmsg", "dmesg -K <path>"},
		{"since", "dmesg --since '1 hour ago'", "dmesg --since <time>"},
		{"until", "dmesg --until '2024-01-01'", "dmesg --until <time>"},
		{"time format", "dmesg --time-format iso", "dmesg --time-format <fmt>"},

		// macOS flags
		{"macos core", "dmesg -M /var/crash/core", "dmesg -M <path>"},
		{"macos system", "dmesg -N /mach_kernel", "dmesg -N <path>"},

		// Combined flags with arguments
		{"timestamp and level", "dmesg -T --level err,warn", "dmesg -T --level <level>"},
		{"human with facility", "dmesg -H -f kern", "dmesg -H -f <facility>"},
		{"kernel json", "dmesg -kJ", "dmesg -kJ"},

		// Redirect
		{"redirect", "dmesg > /tmp/kern.log", "dmesg > <path>"},
		{"pipe grep", "dmesg -T | grep error", "dmesg -T | grep <pattern>"},
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
	t.Run("different levels collide", func(t *testing.T) {
		a := Normalize("dmesg --level err")
		b := Normalize("dmesg --level warn")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different facilities collide", func(t *testing.T) {
		a := Normalize("dmesg -f kern")
		b := Normalize("dmesg -f daemon,user")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different buffer sizes collide", func(t *testing.T) {
		a := Normalize("dmesg -s 8192")
		b := Normalize("dmesg -s 65536")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("dmesg -F /var/log/kern.log")
		subshell := Normalize("dmesg -F $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
