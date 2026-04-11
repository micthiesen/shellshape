package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestLsblk(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"no args", "lsblk", "lsblk"},
		{"single device", "lsblk /dev/sda", "lsblk <path>"},
		{"filesystem info", "lsblk -f", "lsblk -f"},
		{"all devices", "lsblk -a", "lsblk -a"},
		{"human readable bytes", "lsblk -b", "lsblk -b"},

		// Boolean flags
		{"json output", "lsblk -J", "lsblk -J"},
		{"list format", "lsblk -l", "lsblk -l"},
		{"scsi devices", "lsblk --scsi", "lsblk --scsi"},
		{"no headings", "lsblk -n", "lsblk -n"},
		{"full paths", "lsblk -p", "lsblk -p"},
		{"combined flags", "lsblk -fnp", "lsblk -fnp"},
		{"topology", "lsblk -t", "lsblk -t"},
		{"discard info", "lsblk --discard", "lsblk --discard"},

		// Flags with arguments
		{"output columns short", "lsblk -o NAME,SIZE,UUID", "lsblk -o <columns>"},
		{"output columns long", "lsblk --output NAME,SIZE,FSTYPE", "lsblk --output <columns>"},
		{"exclude majors", "lsblk -e 7", "lsblk -e N"},
		{"exclude majors csv", "lsblk -e 7,1", "lsblk -e <columns>"},
		{"exclude long", "lsblk --exclude 7", "lsblk --exclude N"},
		{"include majors", "lsblk -I 8", "lsblk -I N"},
		{"include long", "lsblk --include 8,259", "lsblk --include <columns>"},
		{"sort column", "lsblk -x SIZE", "lsblk -x <column>"},
		{"sort long", "lsblk --sort NAME", "lsblk --sort <column>"},
		{"width", "lsblk -w 80", "lsblk -w N"},
		{"width long", "lsblk --width 120", "lsblk --width N"},
		{"output with equals", "lsblk --output=NAME,SIZE", "lsblk --output=<val>"},

		// Multiple positionals
		{"multiple devices", "lsblk /dev/sda /dev/sdb", "lsblk <path>+"},
		{"flags and devices", "lsblk -f /dev/sda /dev/nvme0n1", "lsblk -f <path>+"},

		// Edge cases
		{"redirect", "lsblk -f > output.txt", "lsblk -f > <path>"},
		{"columns and device", "lsblk -o NAME,SIZE /dev/sda", "lsblk -o <columns> <path>"},
		{"pipe to grep", "lsblk -f | grep ext4", "lsblk -f | grep <pattern>"},
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
	t.Run("different devices collide", func(t *testing.T) {
		a := shellshape.Normalize("lsblk /dev/sda")
		b := shellshape.Normalize("lsblk /dev/nvme0n1")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different columns collide", func(t *testing.T) {
		a := shellshape.Normalize("lsblk -o NAME,SIZE")
		b := shellshape.Normalize("lsblk -o UUID,LABEL,MOUNTPOINT")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// SAFETY TEST
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("lsblk /dev/sda")
		subshell := shellshape.Normalize("lsblk $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
