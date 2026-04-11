package shellshape

func init() {
	Register("curl", handleCurl)
}

// handleCurl handles the curl command.
// URL positionals are classified via classifyToken (producing <http-uri>, etc.).
// Data flags (-d, --data, --json, -F, etc.) collapse the next arg to <data>.
// Header flags (-H, --header) collapse the next arg to <header>.
// Method flag (-X, --request) collapses the next arg to <method>.
// Path/file flags (-o, -T, -b, -c, etc.) collapse the next arg to <path>.
// Numeric flags (--retry, -m, --connect-timeout, etc.) collapse the next arg to N.
// User-agent flag (-A) collapses the next arg to <str>.
// Other argument-consuming flags collapse the next arg to <val>.
func handleCurl(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	dataFlags := map[string]bool{
		"-d": true, "--data": true, "--data-ascii": true, "--data-binary": true,
		"--data-raw": true, "--data-urlencode": true, "--json": true,
		"-F": true, "--form": true, "--form-string": true,
		"-u": true, "--user": true, "-U": true, "--proxy-user": true,
		"--url-query": true,
	}

	headerFlags := map[string]bool{
		"-H": true, "--header": true, "--proxy-header": true,
	}

	methodFlags := map[string]bool{
		"-X": true, "--request": true,
	}

	pathFlags := map[string]bool{
		"-o": true, "--output": true, "--output-dir": true,
		"-T": true, "--upload-file": true,
		"-b": true, "--cookie": true,
		"-c": true, "--cookie-jar": true,
		"-D": true, "--dump-header": true,
		"-K": true, "--config": true,
		"-E": true, "--cert": true,
		"--cert-type": true, "--cacert": true, "--capath": true,
		"--key": true, "--crlfile": true,
		"--trace": true, "--trace-ascii": true,
		"--stderr": true, "--random-file": true, "--egd-file": true,
		"--netrc-file": true, "--libcurl": true,
		"--unix-socket": true, "--abstract-unix-socket": true,
		"--etag-compare": true, "--etag-save": true,
		"--hsts": true, "--alt-svc": true,
		"--pubkey": true, "--proxy-key": true,
	}

	numericFlags := map[string]bool{
		"-m": true, "--max-time": true,
		"--connect-timeout": true, "--expect100-timeout": true,
		"--happy-eyeballs-timeout-ms": true,
		"--keepalive-time":            true,
		"--max-redirs":                true, "--max-filesize": true,
		"--retry": true, "--retry-delay": true, "--retry-max-time": true,
		"-Y": true, "--speed-limit": true,
		"-y": true, "--speed-time": true,
		"--parallel-max": true,
		"--tftp-blksize": true,
		"-C":             true, "--continue-at": true,
	}

	strFlags := map[string]bool{
		"-A": true, "--user-agent": true,
	}

	// Flags that take an argument but get generic <val>
	valFlags := map[string]bool{
		"-w": true, "--write-out": true,
		"-e": true, "--referer": true,
		"-x": true, "--proxy": true,
		"--preproxy": true, "--proxy1.0": true,
		"--url": true, "--doh-url": true, "--ipfs-gateway": true,
		"--resolve": true, "--connect-to": true,
		"--interface": true, "--dns-interface": true,
		"--dns-servers": true, "--dns-ipv4-addr": true, "--dns-ipv6-addr": true,
		"--proto": true, "--proto-default": true, "--proto-redir": true,
		"--ciphers": true, "--curves": true,
		"--tls-max": true, "--tls13-ciphers": true,
		"--proxy-tls13-ciphers": true,
		"--tlsauthtype":         true, "--tlsuser": true, "--tlspassword": true,
		"--proxy-tlsauthtype": true, "--proxy-tlsuser": true, "--proxy-tlspassword": true,
		"--limit-rate": true, "--rate": true,
		"--local-port": true, "-r": true, "--range": true,
		"--noproxy": true, "--request-target": true,
		"-t": true, "--telnet-option": true,
		"-z": true, "--time-cond": true,
		"--ftp-method": true, "--ftp-ssl-ccc-mode": true,
		"-P": true, "--ftp-port": true, "--ftp-account": true,
		"--ftp-alternative-to-user": true,
		"--login-options":           true, "--delegation": true,
		"--mail-from": true, "--mail-rcpt": true, "--mail-auth": true,
		"--sasl-authzid": true, "--service-name": true,
		"--proxy-service-name": true, "--socks5-gssapi-service": true,
		"--engine": true, "--key-type": true, "--proxy-key-type": true,
		"--krb": true, "--oauth2-bearer": true,
		"--aws-sigv4": true, "--hostpubmd5": true, "--hostpubsha256": true,
		"--pinnedpubkey": true, "--proxy-pinnedpubkey": true,
		"--proxy-cert": true, "--proxy-cert-type": true,
		"--proxy-cacert": true, "--proxy-capath": true,
		"--proxy-crlfile": true, "--proxy-pass": true,
		"--pass": true, "--trace-config": true,
		"--create-file-mode": true, "--variable": true,
		"--haproxy-clientip": true,
		"-Q":                 true, "--quote": true,
		"--socks4": true, "--socks4a": true,
		"--socks5": true, "--socks5-hostname": true,
		"-h": true, "--help": true,
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

		if dataFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<data>")
				i++
			}
			continue
		}

		if headerFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<header>")
				i++
			}
			continue
		}

		if methodFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<method>")
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

		if strFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<str>")
				i++
			}
			continue
		}

		if valFlags[tok] {
			result = append(result, tok)
			i++
			if i < len(args) {
				result = append(result, "<val>")
				i++
			}
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional: classify (URLs become <http-uri>, etc.)
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
