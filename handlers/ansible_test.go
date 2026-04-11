package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestAnsible(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic ansible usage
		{"ping all", "ansible all -m ping", "ansible <target> -m <module>"},
		{"shell module with args", "ansible webservers -m shell -a 'uptime'", "ansible <target> -m <module> -a <val>"},
		{"command module with args", "ansible dbservers -m command -a 'df -h'", "ansible <target> -m <module> -a <val>"},
		{"with inventory", "ansible all -i hosts.ini -m ping", "ansible <target> -i <path> -m <module>"},
		{"with user", "ansible all -u deploy -m setup", "ansible <target> -u <val> -m <module>"},
		{"become flags", "ansible all --become --ask-become-pass -m command -a 'whoami'", "ansible <target> --become --ask-become-pass -m <module> -a <val>"},
		{"with limit", "ansible all -l webservers -m ping", "ansible <target> -l <target> -m <module>"},
		{"with extra vars", "ansible all -e 'foo=bar' -m debug -a 'var=foo'", "ansible <target> -e <val> -m <module> -a <val>"},
		{"with forks", "ansible all -f 20 -m ping", "ansible <target> -f N -m <module>"},
		{"with timeout", "ansible all -t 30 -m ping", "ansible <target> -t N -m <module>"},
		{"with private key", "ansible all --private-key /home/user/.ssh/id_rsa -m ping", "ansible <target> --private-key <path> -m <module>"},
		{"verbose", "ansible all -vvv -m ping", "ansible <target> -vvv -m <module>"},
		{"list hosts", "ansible webservers --list-hosts", "ansible <target> --list-hosts"},

		// ansible-playbook usage
		{"playbook simple", "ansible-playbook site.yml", "ansible-playbook <path>"},
		{"playbook with inventory", "ansible-playbook -i inventory/prod deploy.yml", "ansible-playbook -i <path>+"},
		{"playbook with extra vars", "ansible-playbook playbook.yml -e 'version=1.2'", "ansible-playbook <path> -e <val>"},
		{"playbook with tags", "ansible-playbook playbook.yml --tags deploy,setup", "ansible-playbook <path> --tags <val>"},
		{"playbook with skip tags", "ansible-playbook playbook.yml --skip-tags slow", "ansible-playbook <path> --skip-tags <val>"},
		{"playbook with limit", "ansible-playbook playbook.yml --limit staging", "ansible-playbook <path> --limit <target>"},
		{"playbook with start-at", "ansible-playbook playbook.yml --start-at 'Install nginx'", "ansible-playbook <path> --start-at <val>"},
		{"playbook check mode", "ansible-playbook playbook.yml -C -D", "ansible-playbook <path> -C -D"},
		{"playbook with forks and user", "ansible-playbook playbook.yml -f 10 -u deploy", "ansible-playbook <path> -f N -u <val>"},
		{"playbook with vault", "ansible-playbook playbook.yml --vault-password-file /tmp/vault.txt", "ansible-playbook <path> --vault-password-file <path>"},
		{"playbook with become user", "ansible-playbook playbook.yml -b --become-user root", "ansible-playbook <path> -b --become-user <val>"},
		{"playbook extra vars json file", "ansible-playbook playbook.yml -e @variables.json", "ansible-playbook <path> -e <val>"},
		{"playbook multiple plays", "ansible-playbook site.yml deploy.yml", "ansible-playbook <path>+"},
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
	t.Run("different hosts collide", func(t *testing.T) {
		a := shellshape.Normalize("ansible webservers -m ping")
		b := shellshape.Normalize("ansible dbservers -m ping")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different playbooks collide", func(t *testing.T) {
		a := shellshape.Normalize("ansible-playbook site.yml -e 'env=prod'")
		b := shellshape.Normalize("ansible-playbook deploy.yml -e 'env=staging'")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("ansible all -m ping")
		subshell := shellshape.Normalize("ansible $(cat hosts) -m ping")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
