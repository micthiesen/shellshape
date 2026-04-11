package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestUseradd(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"simple useradd", "useradd john", "useradd john"},
		{"simple userdel", "userdel jane", "userdel jane"},
		{"simple usermod", "usermod bob", "usermod bob"},

		// Flags with path arguments
		{"home dir", "useradd -d /home/custom john", "useradd -d <path> john"},
		{"shell", "useradd --shell /bin/bash john", "useradd --shell <path> john"},
		{"skel dir", "useradd -k /etc/skel.custom -m john", "useradd -k <path> -m john"},
		{"root dir", "useradd -R /mnt/chroot john", "useradd -R <path> john"},
		{"base dir long", "useradd --base-dir /home john", "useradd --base-dir <path> john"},

		// Flags with value arguments
		{"uid", "useradd -u 1001 john", "useradd -u N john"},
		{"uid long", "useradd --uid 1001 john", "useradd --uid N john"},
		{"gid", "useradd -g developers john", "useradd -g <group> john"},
		{"groups", "useradd -G wheel,docker john", "useradd -G <groups> john"},
		{"comment", "useradd -c 'John Doe' john", "useradd -c <comment> john"},
		{"inactive days", "useradd -f 30 john", "useradd -f N john"},
		{"expire date", "useradd -e 2025-12-31 john", "useradd -e <value> john"},
		{"password", "useradd -p encrypted_pw john", "useradd -p <value> john"},
		{"login rename", "usermod -l newname oldname", "usermod -l <value> oldname"},
		{"key override", "useradd -K PASS_MAX_DAYS=99 john", "useradd -K <value> john"},
		{"selinux user", "useradd -Z user_u john", "useradd -Z <value> john"},

		// Boolean flags
		{"create home", "useradd --create-home john", "useradd --create-home john"},
		{"system user", "useradd -r -s /usr/sbin/nologin sysuser", "useradd -r -s <path> sysuser"},
		{"remove userdel", "userdel --remove jane", "userdel --remove jane"},
		{"force userdel", "userdel -rf jane", "userdel -rf jane"},

		// Combined flags
		{"full useradd", "useradd -m -u 1001 -g staff -G wheel,docker -s /bin/zsh -d /home/deploy deploy",
			"useradd -m -u N -g <group> -G <groups> -s <path> -d <path> deploy"},
		{"usermod append groups", "usermod -aG sudo bob", "usermod -aG <groups> bob"},
		{"usermod move home", "usermod --move-home --home /new/home bob", "usermod --move-home --home <path> bob"},

		// Redirect
		{"with redirect", "useradd john 2>/dev/null", "useradd john 2>/dev/null"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TEST: different usernames stay distinct
	t.Run("different usernames distinct", func(t *testing.T) {
		a := shellshape.Normalize("useradd alice")
		b := shellshape.Normalize("useradd bob")
		if a == b {
			t.Errorf("different usernames should produce different shapes, both got %q", a)
		}
	})

	// COLLISION TEST: different paths collapse to same shape
	t.Run("different paths collide", func(t *testing.T) {
		a := shellshape.Normalize("useradd -d /home/alice alice")
		b := shellshape.Normalize("useradd -d /opt/deploy deploy")
		if a == b {
			t.Errorf("expected different shapes for different usernames")
		}
		// But same username with different paths should collide
		c := shellshape.Normalize("useradd -d /home/custom john")
		d := shellshape.Normalize("useradd -d /opt/custom john")
		if c != d {
			t.Errorf("expected %q == %q", c, d)
		}
	})

	// COLLISION TEST: different UIDs collapse
	t.Run("different uids collide", func(t *testing.T) {
		a := shellshape.Normalize("useradd -u 1001 john")
		b := shellshape.Normalize("useradd -u 5000 john")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("useradd literal-user")
		subshell := shellshape.Normalize("useradd $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
