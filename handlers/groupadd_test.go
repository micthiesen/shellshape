package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestGroupadd(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"basic group", "groupadd developers", "groupadd <group>"},
		{"system group", "groupadd --system syslog", "groupadd --system <group>"},
		{"system group short", "groupadd -r daemon", "groupadd -r <group>"},

		// Flags with arguments
		{"gid long", "groupadd --gid 1001 mygroup", "groupadd --gid N <group>"},
		{"gid short", "groupadd -g 500 staff", "groupadd -g N <group>"},
		{"force flag", "groupadd -f existing", "groupadd -f <group>"},
		{"non-unique gid", "groupadd -o -g 0 rootdup", "groupadd -o -g N <group>"},
		{"root chroot", "groupadd -R /mnt/chroot mygroup", "groupadd -R <path> <group>"},
		{"prefix dir", "groupadd -P /target mygroup", "groupadd -P <path> <group>"},
		{"password", "groupadd -p secret123 mygroup", "groupadd -p <str> <group>"},
		{"users list", "groupadd -U alice,bob,charlie mygroup", "groupadd -U <str> <group>"},
		{"key override", "groupadd -K GID_MIN=100 mygroup", "groupadd -K <str> <group>"},
		{"combined flags", "groupadd -r -g 999 -K GID_MAX=999 sysgroup", "groupadd -r -g N -K <str> <group>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TEST: different group names collapse to same shape
	t.Run("different groups collide", func(t *testing.T) {
		a := shellshape.Normalize("groupadd developers")
		b := shellshape.Normalize("groupadd marketing")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// COLLISION TEST: different GIDs collapse to same shape
	t.Run("different gids collide", func(t *testing.T) {
		a := shellshape.Normalize("groupadd -g 1001 devs")
		b := shellshape.Normalize("groupadd -g 9999 ops")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("groupadd mygroup")
		subshell := shellshape.Normalize("groupadd $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}

func TestGroupmod(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"basic rename", "groupmod --new-name newgrp oldgrp", "groupmod --new-name <group>+"},
		{"rename short", "groupmod -n newgrp oldgrp", "groupmod -n <group>+"},
		{"change gid", "groupmod --gid 2000 mygroup", "groupmod --gid N <group>"},
		{"change gid short", "groupmod -g 500 staff", "groupmod -g N <group>"},
		{"non-unique", "groupmod -o -g 0 rootlike", "groupmod -o -g N <group>"},
		{"password", "groupmod -p encrypted mygroup", "groupmod -p <str> <group>"},
		{"append users", "groupmod -a -U alice,bob mygroup", "groupmod -a -U <str> <group>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("groupmod -n newname oldgroup")
		subshell := shellshape.Normalize("groupmod -n $(dangerous-command) oldgroup")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}

func TestGroupdel(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"basic delete", "groupdel oldgroup", "groupdel <group>"},
		{"force delete", "groupdel --force badgroup", "groupdel --force <group>"},
		{"force short", "groupdel -f mygroup", "groupdel -f <group>"},
		{"root chroot", "groupdel -R /mnt/chroot mygroup", "groupdel -R <path> <group>"},
		{"prefix", "groupdel -P /target mygroup", "groupdel -P <path> <group>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// COLLISION TEST: different group names collapse
	t.Run("different groups collide", func(t *testing.T) {
		a := shellshape.Normalize("groupdel marketing")
		b := shellshape.Normalize("groupdel engineering")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST: subshell must not collapse
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("groupdel mygroup")
		subshell := shellshape.Normalize("groupdel $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
