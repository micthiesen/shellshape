package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestPerl(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"script file", "perl script.pl", "perl <script>"},
		{"script with args", "perl script.pl foo.txt bar.txt", "perl <script> <path>+"},
		{"version", "perl -v", "perl -v"},
		{"syntax check", "perl -c script.pl", "perl -c <script>"},
		{"debug mode", "perl -d script.pl", "perl -d <script>"},
		{"warnings", "perl -W script.pl", "perl -W <script>"},

		// Expression flags
		{"e flag", "perl -e 'print 42'", "perl -e <perl-expr>"},
		{"E flag", "perl -E 'say 42'", "perl -E <perl-expr>"},
		{"multiple e flags", "perl -e 'use strict;' -e 'print 1'", "perl -e <perl-expr> -e <perl-expr>"},

		// Common one-liner patterns
		{"ne one-liner", "perl -ne 'print if /foo/' file.txt", "perl -ne <perl-expr> <path>"},
		{"pe one-liner", "perl -pe 's/foo/bar/g' file.txt", "perl -pe <perl-expr> <path>"},
		{"lane one-liner", "perl -lane 'print $F[0]' data.csv", "perl -lane <perl-expr> <path>"},
		{"ane one-liner", "perl -ane 'print $F[2]' data.txt", "perl -ane <perl-expr> <path>"},
		{"nle one-liner", "perl -nle 'print if /pattern/' file.txt", "perl -nle <perl-expr> <path>"},

		// In-place editing
		{"in-place no backup", "perl -p -i -e 's/foo/bar/g' file.txt", "perl -p -i -e <perl-expr> <path>"},
		{"in-place fused backup", "perl -pi.bak -e 's/foo/bar/g' file.txt", "perl -pi.bak -e <perl-expr> <path>"},
		{"bare -i flag", "perl -i -pe 's/foo/bar/g' file.txt", "perl -i -pe <perl-expr> <path>"},

		// Module flags
		{"M flag", "perl -MData::Dumper -e 'print Dumper(1)'", "perl -M <val> -e <perl-expr>"},
		{"m flag", "perl -mstrict -e 'print 1'", "perl -m <val> -e <perl-expr>"},
		{"M fused with module", "perl -MPOSIX -e 'print ceil(1.5)'", "perl -M <val> -e <perl-expr>"},

		// Include path
		{"I flag separate", "perl -I /usr/lib/perl5 script.pl", "perl -I <path> <script>"},
		{"I flag fused", "perl -I/usr/lib/perl5 script.pl", "perl -I <path> <script>"},

		// Multiple files
		{"multiple input files", "perl -ne 'print if /foo/' a.txt b.txt c.txt", "perl -ne <perl-expr> <path>+"},

		// Redirects
		{"redirect output", "perl -e 'print 42' > out.txt", "perl -e <perl-expr> > <path>"},
		{"pipe input", "perl -pe 's/foo/bar/g' < input.txt", "perl -pe <perl-expr> < <path>"},

		// 0 flag (record separator)
		{"null separator", "perl -0ne 'print' file.txt", "perl -0ne <perl-expr> <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("shellshape.Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	// Collision tests
	t.Run("different scripts collide", func(t *testing.T) {
		a := shellshape.Normalize("perl deploy.pl")
		b := shellshape.Normalize("perl backup.pl")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different expressions collide", func(t *testing.T) {
		a := shellshape.Normalize("perl -pe 's/foo/bar/g' data.txt")
		b := shellshape.Normalize("perl -pe 's/hello/world/g' notes.txt")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	// Subshell safety test
	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("perl script.pl")
		subshell := shellshape.Normalize("perl $(dangerous-command)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
