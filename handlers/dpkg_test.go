package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestDpkg(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"install deb", "dpkg -i ./foo_1.0_amd64.deb", "dpkg -i <path>"},
		{"install long", "dpkg --install /tmp/package.deb", "dpkg --install <path>"},
		{"remove", "dpkg -r nginx", "dpkg -r nginx"},
		{"remove long", "dpkg --remove curl", "dpkg --remove curl"},
		{"purge", "dpkg -P nginx", "dpkg -P nginx"},
		{"purge long", "dpkg --purge nginx", "dpkg --purge nginx"},

		// List and search (patterns collapse)
		{"list all", "dpkg -l", "dpkg -l"},
		{"list pattern", "dpkg -l 'lib*'", "dpkg -l <pattern>"},
		{"list long", "dpkg --list 'python*'", "dpkg --list <pattern>"},
		{"search file", "dpkg -S /usr/bin/awk", "dpkg -S <pattern>"},
		{"search long", "dpkg --search /etc/nginx/nginx.conf", "dpkg --search <pattern>"},

		// Structural package names
		{"listfiles", "dpkg -L curl", "dpkg -L curl"},
		{"listfiles long", "dpkg --listfiles bash", "dpkg --listfiles bash"},
		{"status", "dpkg -s bash", "dpkg -s bash"},
		{"status long", "dpkg --status nginx", "dpkg --status nginx"},
		{"print-avail", "dpkg -p coreutils", "dpkg -p coreutils"},
		{"configure", "dpkg --configure -a", "dpkg --configure -a"},
		{"configure pkg", "dpkg --configure nginx", "dpkg --configure nginx"},

		// Contents of a deb file
		{"contents", "dpkg -c ./package.deb", "dpkg -c <path>"},
		{"contents long", "dpkg --contents /tmp/foo.deb", "dpkg --contents <path>"},

		// Unpack and build (paths)
		{"unpack", "dpkg --unpack ./pkg.deb", "dpkg --unpack <path>"},
		{"build", "dpkg -b ./build-dir", "dpkg -b <path>"},

		// Selections
		{"get-selections", "dpkg --get-selections", "dpkg --get-selections"},
		{"set-selections", "dpkg --set-selections", "dpkg --set-selections"},

		// Multiple packages
		{"remove multi", "dpkg -r nginx redis curl", "dpkg -r nginx redis curl"},
		{"purge multi", "dpkg --purge nginx redis", "dpkg --purge nginx redis"},

		// Path-consuming flags
		{"admindir", "dpkg --admindir /var/lib/dpkg -l", "dpkg --admindir <path> -l"},
		{"root", "dpkg --root /mnt -i ./pkg.deb", "dpkg --root <path> -i <path>"},
		{"log flag", "dpkg --log /var/log/dpkg.log -i ./pkg.deb", "dpkg --log <path> -i <path>"},

		// Boolean flags preserved
		{"force", "dpkg -i --force-overwrite ./pkg.deb", "dpkg -i --force-overwrite <path>"},
		{"recursive", "dpkg -i -R /var/cache/debs/", "dpkg -i -R <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TESTS: different data values -> same shape
	t.Run("different deb files collide", func(t *testing.T) {
		a := shellshape.Normalize("dpkg -i ./foo_1.0_amd64.deb")
		b := shellshape.Normalize("dpkg -i /tmp/bar_2.3_all.deb")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different search patterns collide", func(t *testing.T) {
		a := shellshape.Normalize("dpkg -S /usr/bin/awk")
		b := shellshape.Normalize("dpkg -S /usr/local/bin/python3")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("dpkg -i ./package.deb")
		subshell := shellshape.Normalize("dpkg -i $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
