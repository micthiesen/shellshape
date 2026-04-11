package handlers

import shellshape "github.com/micthiesen/shellshape"

func init() {
	shellshape.Register("wget", handleWget)
}

// handleWget handles the wget command.
// URL positionals are classified via classifyToken (producing <http-uri>, etc.).
// Output/log flags (-o, -a, -O, -i, -P, etc.) collapse the next arg to <path>.
// Numeric flags (-t, -T, -w, -l, -Q, etc.) collapse the next arg to N.
// Header flags (--header, --warc-header) collapse the next arg to <header>.
// Data flags (--post-data, --body-data) collapse the next arg to <data>.
// User-agent/execute/referer flags collapse the next arg to <str>.
// Other argument-consuming flags collapse the next arg to <val>.
func handleWget(subcommand string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	pathFlags := map[string]bool{
		"-o": true, "--output-file": true,
		"-a": true, "--append-output": true,
		"-O": true, "--output-document": true,
		"-i": true, "--input-file": true,
		"-P": true, "--directory-prefix": true,
		"--config": true, "--rejected-log": true,
		"--load-cookies": true, "--save-cookies": true,
		"--post-file": true, "--body-file": true,
		"--private-key": true, "--ca-certificate": true,
		"--ca-directory": true, "--crl-file": true,
		"--random-file": true, "--certificate": true,
		"--warc-file": true, "--warc-dedup": true,
		"--warc-tempdir": true, "--hsts-file": true,
		"--pinnedpubkey": true,
	}

	numericFlags := map[string]bool{
		"-t": true, "--tries": true,
		"-T": true, "--timeout": true,
		"--dns-timeout": true, "--connect-timeout": true, "--read-timeout": true,
		"-w": true, "--wait": true, "--waitretry": true,
		"-l": true, "--level": true,
		"-Q": true, "--quota": true,
		"--limit-rate": true, "--max-redirect": true,
		"--cut-dirs": true, "--backups": true,
		"--start-pos": true, "--warc-max-size": true,
	}

	headerFlags := map[string]bool{
		"--header": true, "--warc-header": true,
	}

	dataFlags := map[string]bool{
		"--post-data": true, "--body-data": true,
	}

	strFlags := map[string]bool{
		"-U": true, "--user-agent": true,
		"-e": true, "--execute": true,
		"--referer": true,
	}

	valFlags := map[string]bool{
		"--user": true, "--password": true,
		"--http-user": true, "--http-password": true,
		"--ftp-user": true, "--ftp-password": true,
		"--proxy-user": true, "--proxy-password": true,
		"-B": true, "--base": true,
		"--method": true,
		"-A":       true, "--accept": true,
		"-R": true, "--reject": true,
		"--accept-regex": true, "--reject-regex": true,
		"-D": true, "--domains": true,
		"--exclude-domains": true,
		"-I":                true, "--include-directories": true,
		"-X": true, "--exclude-directories": true,
		"--follow-tags": true, "--ignore-tags": true,
		"--report-speed": true, "--progress": true,
		"--prefer-family":  true,
		"--local-encoding": true, "--remote-encoding": true,
		"--default-page": true, "--compression": true,
		"--regex-type": true, "--restrict-file-names": true,
		"--bind-address": true, "--ciphers": true,
		"--certificate-type": true, "--private-key-type": true,
		"--use-askpass": true, "--retry-on-http-error": true,
		"--input-metalink": true, "--metalink-index": true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: numericFlags, Placeholder: "N"},
		{Flags: headerFlags, Placeholder: "<header>"},
		{Flags: dataFlags, Placeholder: "<data>"},
		{Flags: strFlags, Placeholder: "<str>"},
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

		// Positional: classify (URLs become <http-uri>, etc.)
		result = append(result, shellshape.ClassifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
