package shellshape

func init() {
	Register("dig", handleDig)
}

// handleDig handles the dig DNS lookup utility.
// @server tokens become @<server>. Query options (+short, +trace, etc.) are
// preserved verbatim. DNS record types (A, MX, AAAA, etc.) are preserved.
// The domain name positional becomes <name>. Flags like -p take numeric args,
// -f/-k take file paths, -b/-x take addresses, -q takes a name, -y takes a key.
func handleDig(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	// Flags that consume the next token with a specific placeholder.
	addrFlags := map[string]bool{"-b": true, "-x": true}
	pathFlags := map[string]bool{"-f": true, "-k": true}
	numericFlags := map[string]bool{"-p": true}
	nameFlags := map[string]bool{"-q": true}
	keyFlags := map[string]bool{"-y": true}
	// -t and -c consume next token but preserve it verbatim (type/class).
	verbatimFlags := map[string]bool{"-t": true, "-c": true}

	// DNS record types to preserve verbatim as positionals.
	dnsTypes := map[string]bool{
		"A": true, "AAAA": true, "ANY": true, "AXFR": true, "CAA": true,
		"CNAME": true, "DNSKEY": true, "DS": true, "HINFO": true,
		"IXFR": true, "MX": true, "NAPTR": true, "NS": true, "PTR": true,
		"RP": true, "RRSIG": true, "SOA": true, "SRV": true, "TXT": true,
		// Classes that can appear as positionals.
		"IN": true, "CH": true, "HS": true,
	}

	var result []string

	i := 0
	for i < len(args) {
		tok := args[i]

		if isSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// @server notation
		if len(tok) > 1 && tok[0] == '@' {
			result = append(result, "@<server>")
			i++
			continue
		}

		// Query options: +short, +noall, +answer, +timeout=10, etc.
		if len(tok) > 0 && tok[0] == '+' {
			result = append(result, tok)
			i++
			continue
		}

		// Flags that consume the next token.
		if addrFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<addr>")
				i++
			}
			continue
		}

		if pathFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<path>")
				i++
			}
			continue
		}

		if numericFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "N")
				i++
			}
			continue
		}

		if nameFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<name>")
				i++
			}
			continue
		}

		if keyFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<key>")
				i++
			}
			continue
		}

		if verbatimFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, args[i])
				i++
			}
			continue
		}

		// Boolean flags (including -4 and -6 which isFlagToken misses).
		if isFlagToken(tok) || tok == "-4" || tok == "-6" {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: DNS type/class keywords are preserved verbatim.
		if dnsTypes[tok] {
			result = append(result, tok)
			i++
			continue
		}

		// Any other positional is a domain name.
		result = append(result, "<name>")
		i++
	}

	result = append(result, redirects...)
	return result
}
