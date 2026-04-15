package shellshape

import (
	"reflect"
	"testing"
)

func TestEmitPositional(t *testing.T) {
	tests := []struct {
		name        string
		start       []string
		tok         string
		placeholder string
		want        []string
	}{
		{"plain token", nil, "foo", "<val>", []string{"<val>"}},
		{"preserves prefix", []string{"pre"}, "foo", "<path>", []string{"pre", "<path>"}},
		{"subshell preserved", nil, "$(date)", "<val>", []string{"$(date)"}},
		{"subshell in middle", []string{"pre"}, "$(pwd)", "<path>", []string{"pre", "$(pwd)"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EmitPositional(tt.start, tt.tok, tt.placeholder)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EmitPositional(%v, %q, %q) = %v, want %v", tt.start, tt.tok, tt.placeholder, got, tt.want)
			}
		})
	}
}

func TestConsumeUntil(t *testing.T) {
	terminators := map[string]bool{";": true, "+": true}

	tests := []struct {
		name           string
		args           []string
		i              int
		wantTerminator string
		wantNext       int
	}{
		{"semicolon terminator", []string{"grep", "-l", "foo", "{}", ";"}, 0, ";", 5},
		{"plus terminator", []string{"rm", "{}", "+"}, 0, "+", 3},
		{"mid-stream start", []string{"a", "b", "c", ";", "d"}, 2, ";", 4},
		{"no terminator", []string{"grep", "-l", "foo"}, 0, "", 3},
		{"empty slice", nil, 0, "", 0},
		{"start past end", []string{"a"}, 5, "", 1},
		{"terminator immediately", []string{";", "a"}, 0, ";", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			term, next := ConsumeUntil(tt.args, tt.i, terminators)
			if term != tt.wantTerminator || next != tt.wantNext {
				t.Errorf("ConsumeUntil(%v, %d) = (%q, %d), want (%q, %d)",
					tt.args, tt.i, term, next, tt.wantTerminator, tt.wantNext)
			}
		})
	}
}

func TestRepairLeadingFlagSubcommand(t *testing.T) {
	// Typical setup: a handler with --filter (val), -C (path), and some
	// boolean flags. Mimics the pnpm/systemctl family.
	valFlags := map[string]bool{"--filter": true, "-F": true}
	pathFlags := map[string]bool{"-C": true, "--config": true}
	categories := []FlagCategory{
		{Flags: valFlags, Placeholder: "<val>"},
		{Flags: pathFlags, Placeholder: "<path>"},
	}
	verbatimFlags := map[string]bool{"-t": true, "--type": true}

	tests := []struct {
		name               string
		subcommand         string
		args               []string
		startResult        []string
		categories         []FlagCategory
		verbatim           map[string]bool
		wantResult         []string
		wantRealSubcommand string
		wantNextI          int
	}{
		{
			name:               "not a flag: no-op",
			subcommand:         "install",
			args:               []string{"react"},
			startResult:        []string{"exe", "install"},
			categories:         categories,
			wantResult:         []string{"exe", "install"},
			wantRealSubcommand: "install",
			wantNextI:          0,
		},
		{
			name:               "leading val flag then real subcommand",
			subcommand:         "--filter",
			args:               []string{"infra", "exec", "tsgo", "--noEmit"},
			startResult:        []string{"exe", "--filter"},
			categories:         categories,
			wantResult:         []string{"exe", "--filter", "<val>", "exec"},
			wantRealSubcommand: "exec",
			wantNextI:          2,
		},
		{
			name:               "leading path flag then real subcommand",
			subcommand:         "-C",
			args:               []string{"/app", "add", "react"},
			startResult:        []string{"exe", "-C"},
			categories:         categories,
			wantResult:         []string{"exe", "-C", "<path>", "add"},
			wantRealSubcommand: "add",
			wantNextI:          2,
		},
		{
			name:               "two leading flags before subcommand",
			subcommand:         "--filter",
			args:               []string{"infra", "-C", "/app", "add"},
			startResult:        []string{"exe", "--filter"},
			categories:         categories,
			wantResult:         []string{"exe", "--filter", "<val>", "-C", "<path>", "add"},
			wantRealSubcommand: "add",
			wantNextI:          4,
		},
		{
			name:               "boolean flag leaked (no value)",
			subcommand:         "--dry-run",
			args:               []string{"enable"},
			startResult:        []string{"exe", "--dry-run"},
			categories:         nil,
			wantResult:         []string{"exe", "--dry-run", "enable"},
			wantRealSubcommand: "enable",
			wantNextI:          1,
		},
		{
			name:               "subshell as leading flag value preserved",
			subcommand:         "--filter",
			args:               []string{"$(pick-workspace)", "exec", "tsgo"},
			startResult:        []string{"exe", "--filter"},
			categories:         categories,
			wantResult:         []string{"exe", "--filter", "$(pick-workspace)", "exec"},
			wantRealSubcommand: "exec",
			wantNextI:          2,
		},
		{
			name:               "verbatim-value flag leaked",
			subcommand:         "-t",
			args:               []string{"service", "list-units"},
			startResult:        []string{"exe", "-t"},
			categories:         nil,
			verbatim:           verbatimFlags,
			wantResult:         []string{"exe", "-t", "service", "list-units"},
			wantRealSubcommand: "list-units",
			wantNextI:          2,
		},
		{
			name:               "fused flag leaked",
			subcommand:         "--filter=infra",
			args:               []string{"exec", "tsgo"},
			startResult:        []string{"exe", "--filter=infra"},
			categories:         categories,
			wantResult:         []string{"exe", "--filter=infra", "exec"},
			wantRealSubcommand: "exec",
			wantNextI:          1,
		},
		{
			name:               "exhausted without finding subcommand",
			subcommand:         "--filter",
			args:               []string{"infra"},
			startResult:        []string{"exe", "--filter"},
			categories:         categories,
			wantResult:         []string{"exe", "--filter", "<val>"},
			wantRealSubcommand: "",
			wantNextI:          1,
		},
		{
			name:               "intermediate subshell before subcommand",
			subcommand:         "--filter",
			args:               []string{"infra", "$(date)", "exec"},
			startResult:        []string{"exe", "--filter"},
			categories:         categories,
			wantResult:         []string{"exe", "--filter", "<val>", "$(date)", "exec"},
			wantRealSubcommand: "exec",
			wantNextI:          3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, realSub, nextI := RepairLeadingFlagSubcommand(
				tt.subcommand, tt.args, 0, tt.startResult, tt.categories, tt.verbatim,
			)
			if !reflect.DeepEqual(got, tt.wantResult) {
				t.Errorf("result = %v, want %v", got, tt.wantResult)
			}
			if realSub != tt.wantRealSubcommand {
				t.Errorf("realSubcommand = %q, want %q", realSub, tt.wantRealSubcommand)
			}
			if nextI != tt.wantNextI {
				t.Errorf("nextI = %d, want %d", nextI, tt.wantNextI)
			}
		})
	}
}
