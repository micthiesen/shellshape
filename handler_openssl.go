package shellshape

func init() {
	Register("openssl", handleOpenssl)
}

// handleOpenssl handles the openssl command.
// The subcommand (req, x509, s_client, enc, etc.) is already extracted by
// normalizeSingleCommand. This handler processes the remaining tokens after
// the subcommand.
//
// Path flags (-in, -out, -key, -cert, etc.) consume the next token as <path>.
// Value flags (-subj, -passin, -connect, etc.) consume the next token as <val>.
// Numeric flags (-days, -iter, etc.) consume the next token as N.
// Boolean flags are kept verbatim.
// Remaining positionals use classifyToken.
func handleOpenssl(subcommand string, tokens []string) []string {
	args, redirects := splitRedirects(tokens)

	pathFlags := map[string]bool{
		"-in": true, "-out": true, "-key": true, "-signkey": true,
		"-keyout": true, "-cert": true, "-certfile": true,
		"-CA": true, "-CAkey": true, "-CAfile": true, "-CApath": true,
		"-CAstore": true, "-config": true, "-inkey": true,
		"-signature": true, "-writerand": true, "-rand": true,
		"-sign": true, "-verify": true, "-prverify": true,
		"-sess_in": true, "-sess_out": true,
		"-cert_chain": true, "-requestCAfile": true,
		"-psk_session": true, "-ctlogfile": true,
	}

	valueFlags := map[string]bool{
		"-subj": true, "-passin": true, "-passout": true,
		"-pass": true, "-password": true,
		"-k": true, "-K": true, "-S": true, "-iv": true,
		"-hmac": true, "-mac": true, "-macopt": true,
		"-connect": true, "-host": true, "-bind": true,
		"-proxy": true, "-proxy_user": true, "-proxy_pass": true,
		"-unix": true, "-psk": true, "-psk_identity": true, "-name": true,
		"-engine": true, "-md": true, "-newkey": true,
		"-pkeyopt": true, "-sigopt": true, "-vfyopt": true,
		"-addext": true, "-extensions": true, "-reqexts": true,
		"-set_serial": true, "-nameopt": true, "-certopt": true,
		"-reqopt": true, "-dateopt": true,
		"-provider": true, "-provparam": true, "-propquery": true,
		"-provider-path": true,
		"-cipher":        true, "-keygen_engine": true,
		"-copy_extensions": true, "-ext": true,
		"-checkhost": true, "-checkemail": true, "-checkip": true,
		"-starttls": true, "-dane_tlsa_domain": true, "-dane_tlsa_rrdata": true,
		"-bufsize": true, "-inform": true, "-outform": true,
		"-keyform": true, "-certform": true,
		"-ssl_config": true, "-ssl_client_engine": true,
		"-skeyopt": true, "-skeymgmt": true,
		"-section": true, "-not_before": true, "-not_after": true,
	}

	numericFlags := map[string]bool{
		"-days": true, "-iter": true, "-saltlen": true,
		"-primes": true, "-port": true,
		"-maxfraglen": true, "-max_send_frag": true,
		"-split_send_frag": true, "-max_pipelines": true,
		"-read_buf": true, "-xoflen": true, "-checkend": true,
		"-multi": true,
	}

	categories := []flagCategory{
		{pathFlags, "<path>"},
		{valueFlags, "<val>"},
		{numericFlags, "N"},
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

		if placeholder, ok := matchFlagCategory(tok, categories); ok {
			result, i = consumeFlagArg(tok, args, i, result, placeholder)
			continue
		}

		if isFlagToken(tok) {
			result = append(result, tok)
			i++
			continue
		}

		// Positional
		result = append(result, classifyToken(tok))
		i++
	}

	result = append(result, redirects...)
	return result
}
