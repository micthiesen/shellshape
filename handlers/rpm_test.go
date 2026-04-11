package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestRpm(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"install rpm file", "rpm -ivh package-1.0.rpm", "rpm -ivh <path>"},
		{"upgrade rpm file", "rpm -Uvh package-2.0.rpm", "rpm -Uvh <path>"},
		{"erase package", "rpm -e nginx", "rpm -e nginx"},
		{"list all", "rpm -qa", "rpm -qa"},
		{"query package files", "rpm -ql nginx", "rpm -ql nginx"},
		{"query which package owns file", "rpm -qf /usr/bin/vim", "rpm -qf <path>"},
		{"query rpm file", "rpm -qpl package.rpm", "rpm -qpl <path>"},
		{"verify package", "rpm -V nginx", "rpm -V nginx"},

		// Long flags
		{"long query all", "rpm --query --all", "rpm --query --all"},
		{"long erase", "rpm --erase mypackage", "rpm --erase mypackage"},
		{"long upgrade", "rpm --upgrade package.rpm", "rpm --upgrade <path>"},
		{"long install", "rpm --install package.rpm", "rpm --install <path>"},

		// Flags with arguments
		{"queryformat", "rpm --queryformat '%{NAME}' -qa", "rpm --queryformat <val> -qa"},
		{"qf shorthand", "rpm --qf '%{NAME}-%{VERSION}' -qa", "rpm --qf <val> -qa"},
		{"root flag", "rpm --root /mnt/sysimage -qa", "rpm --root <path> -qa"},
		{"dbpath flag", "rpm --dbpath /var/lib/rpm.bak -qa", "rpm --dbpath <path> -qa"},

		// Checksig
		{"checksig", "rpm -K package.rpm", "rpm -K <path>"},

		// Multiple packages
		{"erase multiple", "rpm -e foo bar baz", "rpm -e foo bar baz"},
		{"install multiple", "rpm -ivh a.rpm b.rpm", "rpm -ivh <path>+"},

		// Boolean flags
		{"force install", "rpm -ivh --force --nodeps package.rpm", "rpm -ivh --force --nodeps <path>"},
		{"whatrequires", "rpm -q --whatrequires openssl", "rpm -q --whatrequires openssl"},

		// Separate -q and -f flags (second-pass fileMode detection)
		{"query file separate flags", "rpm -q -f /usr/bin/vim", "rpm -q -f <path>"},
		// Separate -q and -p flags
		{"query package separate flags", "rpm -q -p package.rpm", "rpm -q -p <path>"},
		// Build mode detection
		{"build spec", "rpm -bb mypackage.spec", "rpm -bb <path>"},
		// subshell as path-flag argument
		{"root subshell", "rpm --root $(get-root) -qa", "rpm --root $(get-root) -qa"},
		// subshell as val-flag argument
		{"queryformat subshell", "rpm --queryformat $(get-fmt) -qa", "rpm --queryformat $(get-fmt) -qa"},

		// Edge cases
		{"no args", "rpm", "rpm"},
		{"define macro", "rpm --define 'dist .el8' -ba foo.spec", "rpm --define <val> -ba <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS
	t.Run("different rpm files collide", func(t *testing.T) {
		a := shellshape.Normalize("rpm -ivh foo-1.0-1.x86_64.rpm")
		b := shellshape.Normalize("rpm -ivh bar-2.3-4.noarch.rpm")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different queryformat values collide", func(t *testing.T) {
		a := shellshape.Normalize("rpm --queryformat '%{NAME}' -qa")
		b := shellshape.Normalize("rpm --queryformat '%{VERSION}-%{RELEASE}' -qa")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("rpm -ivh package.rpm")
		subshell := shellshape.Normalize("rpm -ivh $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
