package shellshape

import "testing"

func TestDd(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// Basic usage
		{"if and of", "dd if=/dev/sda of=/dev/sdb", "dd if=<path> of=<path>"},
		{"if of bs", "dd if=/dev/sda of=/dev/sdb bs=4096", "dd if=<path> of=<path> bs=N"},
		{"if of bs count", "dd if=disk.img of=/dev/sdb bs=1M count=100", "dd if=<path> of=<path> bs=N count=N"},
		{"if only with bs count", "dd if=/dev/urandom bs=512 count=1", "dd if=<path> bs=N count=N"},
		{"skip and seek", "dd if=input.bin of=output.bin skip=10 seek=20", "dd if=<path> of=<path> skip=N seek=N"},

		// Keyword operands preserved
		{"conv preserved", "dd if=input.bin of=output.bin conv=notrunc", "dd if=<path> of=<path> conv=notrunc"},
		{"conv multiple", "dd if=file.iso of=/dev/disk2 conv=sync,noerror", "dd if=<path> of=<path> conv=sync,noerror"},
		{"status progress", "dd if=/dev/sda of=backup.img status=progress", "dd if=<path> of=<path> status=progress"},
		{"iflag oflag", "dd if=/dev/sda of=/dev/sdb iflag=fullblock oflag=sync", "dd if=<path> of=<path> iflag=fullblock oflag=sync"},

		// All numeric operands
		{"ibs obs", "dd if=a of=b ibs=512 obs=1024", "dd if=<path> of=<path> ibs=N obs=N"},
		{"cbs", "dd if=a of=b cbs=80 conv=block", "dd if=<path> of=<path> cbs=N conv=block"},
		{"iseek oseek", "dd if=a of=b iseek=5 oseek=10", "dd if=<path> of=<path> iseek=N oseek=N"},
		{"files", "dd if=/dev/tape files=3", "dd if=<path> files=N"},
		{"speed", "dd if=a of=b speed=1048576", "dd if=<path> of=<path> speed=N"},

		// Full real-world example
		{"full example", "dd if=/dev/sda of=/tmp/backup.img bs=4m conv=sync,noerror status=progress", "dd if=<path> of=<path> bs=N conv=sync,noerror status=progress"},

		// Redirect
		{"with redirect", "dd if=/dev/zero bs=1024 count=10 > output.bin", "dd if=<path> bs=N count=N > <path>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different paths collide", func(t *testing.T) {
		a := Normalize("dd if=/dev/sda of=/dev/sdb bs=4096")
		b := Normalize("dd if=/tmp/image.iso of=/dev/disk3 bs=4096")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("different sizes collide", func(t *testing.T) {
		a := Normalize("dd if=/dev/zero of=file.dat bs=1024 count=100")
		b := Normalize("dd if=/dev/zero of=file.dat bs=4M count=5000")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := Normalize("dd if=/dev/zero of=output.bin bs=1024")
		subshell := Normalize("dd if=$(dangerous-command) of=output.bin bs=1024")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
