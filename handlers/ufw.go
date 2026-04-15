package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"regexp"
	"strings"
)

func init() {
	shellshape.Register("ufw", handleUfw, shellshape.HandlerOptions{HasSubcommands: true})
}

var ufwPortProtoRE = regexp.MustCompile(`^\d[\d:]*/(tcp|udp|any)$`)

// handleUfw handles ufw (Uncomplicated Firewall).
// Subcommands: enable, disable, status, allow, deny, reject, limit, delete,
// reset, reload, app, logging, default.
// Rule keywords (from, to, port, proto, comment, on) are structural.
// After from/to: IPs/CIDRs -> <val>, "any" stays verbatim.
// After port: numbers -> N, ranges -> <val>.
// After proto: protocol name stays verbatim (tcp, udp).
// After comment: value -> <val>.
// After on: interface name -> <val>.
// Port numbers -> N, port/proto combos -> <val>.
// App names -> <val>.
func handleUfw(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	// Structural keywords in rule syntax (consume next token as IP/CIDR or "any").
	ruleKeywords := map[string]bool{
		"from": true, "to": true,
	}

	// Keywords whose next token is a value (<val>).
	valKeywords := map[string]bool{
		"comment": true,
	}

	// Structural tokens that are always kept verbatim.
	structuralTokens := map[string]bool{
		"numbered": true, "verbose": true,
		"allow": true, "deny": true, "reject": true, "limit": true,
		"list": true, "info": true, "update": true,
		"in": true, "out": true,
		"on": true, "off": true,
		"incoming": true, "outgoing": true, "routed": true,
	}

	// Protocols are structural.
	protocols := map[string]bool{
		"tcp": true, "udp": true, "any": true,
		"ah": true, "esp": true, "gre": true, "ipv6": true,
	}

	var result []string
	i := 0

	// Handle a leading flag that the normalizer extracted as the subcommand.
	result, subcommand, i = shellshape.RepairLeadingFlagSubcommand(
		subcommand, args, i, result, nil, nil,
	)

	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Flags
		if shellshape.IsFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// "port" keyword: next token is a port number or range
		if tok == "port" {
			result = append(result, tok)
			i++
			if i < len(args) {
				if shellshape.IsSubshellToken(args[i]) {
					result = append(result, args[i])
				} else if shellshape.NumberRE.MatchString(args[i]) {
					result = append(result, "N")
				} else {
					result = append(result, "<val>")
				}
				i++
			}
			continue
		}

		// "proto" keyword: next token is kept verbatim (protocol name)
		if tok == "proto" {
			result = append(result, tok)
			i++
			if i < len(args) {
				if shellshape.IsSubshellToken(args[i]) {
					result = append(result, args[i])
				} else {
					result = append(result, args[i]) // tcp, udp, etc. verbatim
				}
				i++
			}
			continue
		}

		// "from"/"to" keywords: next token is IP/CIDR or "any"
		if ruleKeywords[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				if shellshape.IsSubshellToken(args[i]) {
					result = append(result, args[i])
				} else if args[i] == "any" {
					result = append(result, "any")
				} else {
					result = append(result, "<val>")
				}
				i++
			}
			continue
		}

		// "comment" keyword: next token is <val>
		if valKeywords[tok] {
			result, i = shellshape.ConsumeFlagArg(tok, args, i, result, "<val>")
			continue
		}

		// "on" keyword: when followed by an interface name, consume it as <val>.
		// When used as on/off toggle (e.g. "logging on"), keep structural.
		if tok == "on" {
			result = append(result, tok)
			i++
			if i < len(args) && !shellshape.IsFlagToken(args[i]) && !shellshape.IsSubshellToken(args[i]) &&
				!structuralTokens[args[i]] && !ruleKeywords[args[i]] && !protocols[args[i]] &&
				args[i] != "port" && args[i] != "proto" && args[i] != "comment" && args[i] != "on" &&
				!shellshape.NumberRE.MatchString(args[i]) {
				result = append(result, "<val>")
				i++
			}
			continue
		}

		// Structural tokens
		if structuralTokens[tok] {
			result = append(result, tok)
			i++
			continue
		}

		// Protocols as bare tokens
		if protocols[tok] {
			result = append(result, tok)
			i++
			continue
		}

		// Port/proto combos like 22/tcp
		if ufwPortProtoRE.MatchString(tok) {
			result = append(result, "<val>")
			i++
			continue
		}

		// "logging" subcommand: on/off are structural, other levels are <val>
		if subcommand == "logging" && (tok == "on" || tok == "off") {
			result = append(result, tok)
			i++
			continue
		}

		// Numbers (port numbers, rule numbers for delete)
		if shellshape.NumberRE.MatchString(tok) {
			result = append(result, "N")
			i++
			continue
		}

		// Port ranges like 8000:9000
		if strings.Contains(tok, ":") {
			result = append(result, "<val>")
			i++
			continue
		}

		// App names and other values
		result = append(result, "<val>")
		i++
	}

	result = append(result, redirects...)
	return result
}
