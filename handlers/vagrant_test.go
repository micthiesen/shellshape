package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestVagrant(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic subcommands (no args)
		{"init bare", "vagrant init", "vagrant init"},
		{"up bare", "vagrant up", "vagrant up"},
		{"halt bare", "vagrant halt", "vagrant halt"},
		{"destroy bare", "vagrant destroy", "vagrant destroy"},
		{"suspend", "vagrant suspend", "vagrant suspend"},
		{"resume", "vagrant resume", "vagrant resume"},
		{"status", "vagrant status", "vagrant status"},
		{"ssh bare", "vagrant ssh", "vagrant ssh"},
		{"reload bare", "vagrant reload", "vagrant reload"},
		{"provision bare", "vagrant provision", "vagrant provision"},
		{"global-status", "vagrant global-status", "vagrant global-status"},
		{"validate", "vagrant validate", "vagrant validate"},

		// Subcommand with VM name positional
		{"up with name", "vagrant up web", "vagrant up web"},
		{"halt with name", "vagrant halt db", "vagrant halt db"},
		{"ssh with name", "vagrant ssh web", "vagrant ssh web"},
		{"destroy with name", "vagrant destroy worker", "vagrant destroy worker"},

		// Boolean flags
		{"destroy force", "vagrant destroy -f", "vagrant destroy -f"},
		{"destroy force long", "vagrant destroy --force", "vagrant destroy --force"},
		{"destroy graceful", "vagrant destroy --graceful", "vagrant destroy --graceful"},
		{"reload provision", "vagrant reload --provision", "vagrant reload --provision"},
		{"up no-provision", "vagrant up --no-provision", "vagrant up --no-provision"},
		{"up parallel", "vagrant up --no-parallel", "vagrant up --no-parallel"},

		// Flags with arguments
		{"up with provider", "vagrant up --provider virtualbox", "vagrant up --provider <val>"},
		{"up with provision-with", "vagrant up --provision-with shell", "vagrant up --provision-with <val>"},
		{"ssh with command", "vagrant ssh -c 'ls -la'", "vagrant ssh -c <val>"},
		{"ssh with command long", "vagrant ssh --command 'echo hello'", "vagrant ssh --command <val>"},
		{"package with output", "vagrant package --output mybox.box", "vagrant package --output <path>"},

		// box subcommand (second-level subcommand kept verbatim, args collapsed)
		{"box add", "vagrant box add hashicorp/precise32", "vagrant box add <val>"},
		{"box add with provider", "vagrant box add mybox --provider vmware", "vagrant box add <val> --provider <val>"},
		{"box add with checksum", "vagrant box add mybox --checksum abc123 --checksum-type sha256", "vagrant box add <val> --checksum <val> --checksum-type <val>"},
		{"box add with version", "vagrant box add mybox --box-version 1.0", "vagrant box add <val> --box-version <val>"},
		{"box add force", "vagrant box add --force mybox", "vagrant box add --force <val>"},
		{"box list", "vagrant box list", "vagrant box list"},
		{"box remove", "vagrant box remove hashicorp/precise32", "vagrant box remove <val>"},
		{"box remove all", "vagrant box remove mybox --all", "vagrant box remove <val> --all"},
		{"box update", "vagrant box update --box mybox", "vagrant box update --box <val>"},
		{"box outdated global", "vagrant box outdated --global", "vagrant box outdated --global"},
		{"box prune dry-run", "vagrant box prune --dry-run", "vagrant box prune --dry-run"},

		// snapshot subcommand
		{"snapshot save", "vagrant snapshot save mysnap", "vagrant snapshot save <val>"},
		{"snapshot restore", "vagrant snapshot restore mysnap", "vagrant snapshot restore <val>"},
		{"snapshot restore no-start", "vagrant snapshot restore mysnap --no-start", "vagrant snapshot restore <val> --no-start"},
		{"snapshot delete", "vagrant snapshot delete mysnap", "vagrant snapshot delete <val>"},
		{"snapshot list", "vagrant snapshot list", "vagrant snapshot list"},
		{"snapshot push", "vagrant snapshot push", "vagrant snapshot push"},
		{"snapshot pop", "vagrant snapshot pop", "vagrant snapshot pop"},

		// plugin subcommand
		{"plugin install", "vagrant plugin install vagrant-vbguest", "vagrant plugin install <val>"},
		{"plugin uninstall", "vagrant plugin uninstall vagrant-vbguest", "vagrant plugin uninstall <val>"},
		{"plugin list", "vagrant plugin list", "vagrant plugin list"},

		// init with box name
		{"init with box", "vagrant init ubuntu/focal64", "vagrant init <val>"},
		{"init with box and url", "vagrant init mybox https://example.com/box.box", "vagrant init <val>+"},

		// Combined flags and positionals
		{"up provider and name", "vagrant up web --provider vmware_fusion", "vagrant up web --provider <val>"},
		{"destroy force with name", "vagrant destroy -f web", "vagrant destroy -f web"},

		// Redirect
		{"status with redirect", "vagrant status > /tmp/status.txt", "vagrant status > <path>"},
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
	t.Run("different box names collide", func(t *testing.T) {
		a := shellshape.Normalize("vagrant box add hashicorp/precise32")
		b := shellshape.Normalize("vagrant box add ubuntu/focal64")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different provider values collide", func(t *testing.T) {
		a := shellshape.Normalize("vagrant up --provider virtualbox")
		b := shellshape.Normalize("vagrant up --provider vmware_fusion")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different snapshot names collide", func(t *testing.T) {
		a := shellshape.Normalize("vagrant snapshot save before-update")
		b := shellshape.Normalize("vagrant snapshot save clean-state")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("vagrant up myvm")
		subshell := shellshape.Normalize("vagrant up $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
