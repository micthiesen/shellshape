package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestStrace(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"trace command", "strace ls -la /tmp", "strace <cmd>+"},
		{"trace with pid", "strace -p 12345", "strace -p N"},
		{"output to file", "strace -o trace.log ls", "strace -o <path> <cmd>"},
		{"filter syscalls", "strace -e trace=open,read ls", "strace -e <expr> <cmd>"},
		{"follow forks", "strace -f ./myapp arg1 arg2", "strace -f <cmd>+"},
		{"summary mode", "strace -c ls /tmp", "strace -c <cmd>+"},
		{"string size", "strace -s 256 ls", "strace -s N <cmd>"},
		{"trace path", "strace -P /etc/passwd ls", "strace -P <path> <cmd>"},
		{"combined flags", "strace -f -e trace=network -o net.log curl http://example.com", "strace -f -e <expr> -o <path> <cmd>+"},
		{"sort flag", "strace -c -S calls ls", "strace -c -S <val> <cmd>"},
		{"alignment column", "strace -a 40 ls", "strace -a N <cmd>"},
		{"long flag trace", "strace --trace=open ls", "strace --trace=<expr> <cmd>"},
		{"long flag output", "strace --output=trace.log ls", "strace --output=<path> <cmd>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different commands collide", func(t *testing.T) {
		a := shellshape.Normalize("strace -f ./myapp --foo")
		b := shellshape.Normalize("strace -f /usr/bin/other --bar baz")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different pids collide", func(t *testing.T) {
		a := shellshape.Normalize("strace -p 1234")
		b := shellshape.Normalize("strace -p 5678")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("strace ls")
		subshell := shellshape.Normalize("strace $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}

func TestDtrace(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"run script", "dtrace -s script.d", "dtrace -s <path>"},
		{"trace probe", "dtrace -n 'syscall::read:entry'", "dtrace -n <probe>"},
		{"trace command", "dtrace -c ./myapp", "dtrace -c <cmd>"},
		{"trace pid", "dtrace -p 12345", "dtrace -p N"},
		{"buffer size", "dtrace -b 4m", "dtrace -b <val>"},
		{"output file", "dtrace -o trace.out", "dtrace -o <path>"},
		{"trace function", "dtrace -f 'malloc'", "dtrace -f <probe>"},
		{"define macro", "dtrace -D DEBUG=1", "dtrace -D <val>"},
		{"set option", "dtrace -x bufsize=4m", "dtrace -x <val>"},
		{"combined", "dtrace -n 'syscall:::entry' -o out.log -p 999", "dtrace -n <probe> -o <path> -p N"},
		{"list probes", "dtrace -l -n 'syscall::*'", "dtrace -l -n <probe>"},
		{"module", "dtrace -m mymodule", "dtrace -m <probe>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different probes collide", func(t *testing.T) {
		a := shellshape.Normalize("dtrace -n 'syscall::read:entry'")
		b := shellshape.Normalize("dtrace -n 'io:::start'")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("dtrace -c ./myapp")
		subshell := shellshape.Normalize("dtrace -c $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
