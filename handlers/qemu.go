package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"regexp"
)

// sizeRE matches disk/memory size specs like 20G, +10G, 512M, 1T.
var sizeRE = regexp.MustCompile(`^[+-]?\d+[KMGTkmgt]?$`)

func init() {
	for _, name := range []string{"qemu-system-x86_64", "qemu-system-aarch64"} {
		shellshape.Register(name, handleQemuSystem)
	}
	shellshape.Register("qemu-img", handleQemuImg, shellshape.HandlerOptions{HasSubcommands: true})
}

// handleQemuSystem handles qemu-system-* virtual machine emulators.
// Path flags (-hda, -cdrom, -kernel, etc.) collapse to <path>.
// Structural flags (-cpu, -display) keep their argument verbatim.
// Value flags (-m, -smp, -net, -device, etc.) collapse to <val>.
// Remaining positionals use ClassifyToken.
func handleQemuSystem(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	pathFlags := map[string]bool{
		"-hda": true, "-hdb": true, "-hdc": true, "-hdd": true,
		"-cdrom": true, "-kernel": true, "-initrd": true,
		"-bios": true, "-pflash": true, "-dtb": true,
	}

	structuralFlags := map[string]bool{
		"-cpu": true, "-display": true,
	}

	valFlags := map[string]bool{
		"-m": true, "-smp": true,
		"-boot": true, "-net": true, "-nic": true,
		"-serial": true, "-parallel": true, "-monitor": true,
		"-chardev": true, "-device": true, "-drive": true,
		"-netdev": true, "-object": true, "-audiodev": true,
		"-blockdev": true, "-fsdev": true, "-virtfs": true,
		"-usbdevice": true, "-machine": true, "-M": true,
		"-accel": true, "-watchdog": true, "-watchdog-action": true,
		"-rtc": true, "-icount": true, "-soundhw": true,
		"-vga": true, "-vnc": true, "-k": true, "-name": true,
		"-numa": true, "-append": true, "-global": true,
		"-readconfig": true, "-writeconfig": true,
		"-incoming": true, "-spice": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: valFlags, Placeholder: "<val>"},
	}

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if structuralFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if shellshape.IsSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, args[i])
				}
				i++
			}
			continue
		}

		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}

// handleQemuImg handles qemu-img disk image management.
// Value flags (-f, -O, -o, etc.) collapse to <val>.
// Positionals use ClassifyToken.
func handleQemuImg(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	valFlags := map[string]bool{
		"-f": true, "-O": true, "-o": true,
		"-B": true, "-F": true,
		"--object": true, "--image-opts": true,
		"-l": true, "-S": true, "-t": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: valFlags, Placeholder: "<val>"},
	}

	var result []string
	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		if placeholder, ok := shellshape.MatchFlagCategory(tok, categories); ok {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: sizes like 20G, +10G collapse to <val>
		if sizeRE.MatchString(tok) && !shellshape.NumberRE.MatchString(tok) {
			result = append(result, "<val>")
		} else {
			result = append(result, shellshape.ClassifyToken(tok))
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
