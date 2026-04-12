package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"testing"
)

func TestQemuSystem(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"boot from image", "qemu-system-x86_64 -hda /vm/disk.img", "qemu-system-x86_64 -hda <path>"},
		{"memory and cdrom", "qemu-system-x86_64 -m 2G -hda /vm/disk.img -cdrom /iso/os.iso -boot d", "qemu-system-x86_64 -m <val> -hda <path> -cdrom <path> -boot <val>"},
		{"smp and cpu", "qemu-system-x86_64 -smp 4 -cpu host -m 4096 -hda /vm/disk.qcow2", "qemu-system-x86_64 -smp <val> -cpu host -m <val> -hda <path>"},
		{"display type", "qemu-system-x86_64 -display gtk -hda /vm/disk.img", "qemu-system-x86_64 -display gtk -hda <path>"},
		{"network and serial", "qemu-system-x86_64 -net nic -net user -serial stdio -hda /vm/disk.img", "qemu-system-x86_64 -net <val> -net <val> -serial <val> -hda <path>"},
		{"drive flag", "qemu-system-x86_64 -drive file=/vm/disk.qcow2,format=qcow2", "qemu-system-x86_64 -drive <val>"},
		{"aarch64 variant", "qemu-system-aarch64 -M virt -cpu cortex-a57 -m 1024 -hda /vm/arm.img", "qemu-system-aarch64 -M <val> -cpu cortex-a57 -m <val> -hda <path>"},
		{"boolean flags", "qemu-system-x86_64 -enable-kvm -nographic -hda /vm/disk.img", "qemu-system-x86_64 -enable-kvm -nographic -hda <path>"},
		{"device passthrough", "qemu-system-x86_64 -device virtio-net-pci,netdev=net0 -netdev user,id=net0 -hda /vm/disk.img", "qemu-system-x86_64 -device <val> -netdev <val> -hda <path>"},
		{"kernel boot", "qemu-system-x86_64 -kernel /boot/vmlinuz -initrd /boot/initrd.img -append root=/dev/sda1", "qemu-system-x86_64 -kernel <path> -initrd <path> -append <val>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different images collide", func(t *testing.T) {
		a := shellshape.Normalize("qemu-system-x86_64 -m 2G -hda /vm/ubuntu.qcow2")
		b := shellshape.Normalize("qemu-system-x86_64 -m 4G -hda /vm/debian.qcow2")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("qemu-system-x86_64 -hda /vm/disk.img")
		subshell := shellshape.Normalize("qemu-system-x86_64 -hda $(echo evil)")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}

func TestQemuImg(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"create image", "qemu-img create -f qcow2 /vm/disk.qcow2 20G", "qemu-img create -f <val> <path> <val>"},
		{"info", "qemu-img info /vm/disk.qcow2", "qemu-img info <path>"},
		{"convert format", "qemu-img convert -f vmdk -O qcow2 /vm/input.vmdk /vm/output.qcow2", "qemu-img convert -f <val> -O <val> <path>+"},
		{"resize", "qemu-img resize /vm/disk.qcow2 +10G", "qemu-img resize <path> <val>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shellshape.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}

	t.Run("different paths collide", func(t *testing.T) {
		a := shellshape.Normalize("qemu-img info /vm/ubuntu.qcow2")
		b := shellshape.Normalize("qemu-img info /vm/debian.qcow2")
		if a != b {
			t.Errorf("expected %q == %q", a, b)
		}
	})

	t.Run("subshell not collapsed", func(t *testing.T) {
		benign := shellshape.Normalize("qemu-img create -f qcow2 /vm/disk.qcow2 20G")
		subshell := shellshape.Normalize("qemu-img create -f qcow2 $(echo evil) 20G")
		if benign == subshell {
			t.Error("subshell must produce different shape than literal")
		}
	})
}
