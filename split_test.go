package shellshape

import "testing"

func TestSplitTopLevelCompoundCommands(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantSegs int // number of top-level segments
		wantSeps []string
	}{
		// For loop: internal ; should NOT split
		{"for loop single segment", "for f in a b; do echo $f; done", 1, []string{""}},
		// For loop followed by &&: should split into 2 segments
		{"for then and", "for f in a b; do echo $f; done && echo ok", 2, []string{"&&", ""}},
		// Nested for loops: still single segment
		{"nested for", "for x in 1 2; do for y in a b; do echo $x $y; done; done", 1, []string{""}},
		// While loop: internal ; should NOT split
		{"while loop single", "while read line; do echo $line; done", 1, []string{""}},
		// If statement: internal ; should NOT split
		{"if single", "if test -f x; then echo y; fi", 1, []string{""}},
		// If-else: single compound
		{"if else single", "if test -f x; then echo y; else echo n; fi", 1, []string{""}},
		// If chained with outer command
		{"if then and", "if test -f x; then echo y; fi && echo done", 2, []string{"&&", ""}},
		// Plain commands: still split normally
		{"plain semicolons", "echo a; echo b; echo c", 3, []string{";", ";", ""}},
		// Mixed: plain then compound
		{"plain then for", "echo start; for f in a b; do echo $f; done", 2, []string{";", ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			segs := splitTopLevel(tt.input, splitOps)
			if len(segs) != tt.wantSegs {
				t.Errorf("splitTopLevel(%q): got %d segments, want %d", tt.input, len(segs), tt.wantSegs)
				for i, s := range segs {
					t.Logf("  seg[%d] text=%q sep=%q", i, s.text, s.sep)
				}
				return
			}
			for i, wantSep := range tt.wantSeps {
				if segs[i].sep != wantSep {
					t.Errorf("seg[%d].sep = %q, want %q (text=%q)", i, segs[i].sep, wantSep, segs[i].text)
				}
			}
		})
	}
}
