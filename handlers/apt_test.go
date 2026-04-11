package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestApt(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// install subcommand
		{"install single package", "apt install nginx", "apt install nginx"},
		{"install multiple packages", "apt install nginx curl wget", "apt install nginx curl wget"},
		{"install with -y", "apt install -y python3", "apt install -y python3"},
		{"install with target release", "apt install -t stable nginx", "apt install -t <val> nginx"},
		{"install with option", "apt install -o Dpkg::Options::=--force-confold nginx", "apt install -o <val> nginx"},
		{"install with --no-install-recommends", "apt install --no-install-recommends build-essential", "apt install --no-install-recommends build-essential"},

		// remove / purge subcommands
		{"remove package", "apt remove nginx", "apt remove nginx"},
		{"remove with purge flag", "apt remove --purge nginx", "apt remove --purge nginx"},
		{"purge package", "apt purge nginx", "apt purge nginx"},

		// update / upgrade
		{"update", "apt update", "apt update"},
		{"upgrade", "apt upgrade", "apt upgrade"},
		{"upgrade with -y", "apt upgrade -y", "apt upgrade -y"},
		{"dist-upgrade", "apt dist-upgrade", "apt dist-upgrade"},
		{"full-upgrade", "apt full-upgrade", "apt full-upgrade"},

		// search subcommand
		{"search query", "apt search postgres", "apt search <query>"},
		{"search with flags", "apt search --names-only redis", "apt search --names-only <query>"},

		// show subcommand
		{"show package", "apt show nginx", "apt show nginx"},
		{"show multiple", "apt show nginx curl", "apt show nginx curl"},

		// list subcommand
		{"list all", "apt list", "apt list"},
		{"list installed", "apt list --installed", "apt list --installed"},
		{"list upgradable", "apt list --upgradable", "apt list --upgradable"},

		// autoremove / clean
		{"autoremove", "apt autoremove", "apt autoremove"},
		{"autoremove with -y", "apt autoremove -y", "apt autoremove -y"},
		{"autoclean", "apt autoclean", "apt autoclean"},
		{"clean", "apt clean", "apt clean"},

		// depends / rdepends
		{"depends", "apt depends nginx", "apt depends nginx"},
		{"rdepends", "apt rdepends nginx", "apt rdepends nginx"},

		// policy (apt-cache)
		{"policy", "apt-cache policy nginx", "apt-cache policy nginx"},
		{"showpkg", "apt-cache showpkg curl", "apt-cache showpkg curl"},

		// apt-get alias
		{"apt-get install", "apt-get install -y nginx", "apt-get install -y nginx"},
		{"apt-get update", "apt-get update", "apt-get update"},
		{"apt-get upgrade", "apt-get upgrade", "apt-get upgrade"},
		{"apt-get dist-upgrade", "apt-get dist-upgrade", "apt-get dist-upgrade"},
		{"apt-get remove", "apt-get remove nginx", "apt-get remove nginx"},
		{"apt-get purge", "apt-get purge --auto-remove nginx", "apt-get purge --auto-remove nginx"},
		{"apt-get autoremove", "apt-get autoremove", "apt-get autoremove"},
		{"apt-get clean", "apt-get clean", "apt-get clean"},
		{"apt-get download", "apt-get download firefox", "apt-get download firefox"},

		// apt-cache alias
		{"apt-cache search", "apt-cache search redis", "apt-cache search <query>"},
		{"apt-cache show", "apt-cache show nginx", "apt-cache show nginx"},
		{"apt-cache depends", "apt-cache depends curl", "apt-cache depends curl"},
		{"apt-cache rdepends", "apt-cache rdepends curl", "apt-cache rdepends curl"},

		// option flag with long form
		{"long option flag", "apt install --option Dpkg::Options::=--force-confold nginx", "apt install --option <val> nginx"},

		// boolean flags preserved
		{"fix-broken", "apt install --fix-broken", "apt install --fix-broken"},
		{"simulate", "apt install -s nginx", "apt install -s nginx"},
		{"quiet", "apt install -qq nginx", "apt install -qq nginx"},
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
	t.Run("different search terms collide", func(t *testing.T) {
		a := shellshape.Normalize("apt search postgres")
		b := shellshape.Normalize("apt search redis")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different search terms collide apt-cache", func(t *testing.T) {
		a := shellshape.Normalize("apt-cache search nginx")
		b := shellshape.Normalize("apt-cache search memcached")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("apt install literal-arg")
		subshell := shellshape.Normalize("apt install $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
