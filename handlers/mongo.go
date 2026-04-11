package handlers

import (
	shellshape "github.com/micthiesen/shellshape"
	"strings"
)

func init() {
	for _, name := range []string{"mongo", "mongosh"} {
		shellshape.Register(name, handleMongo)
	}
}

// handleMongo handles mongo and mongosh (MongoDB shell).
// The first positional is the database name or connection string (<db> or <mongodb-uri>).
// The second positional (mongo legacy) is a script file (<path>).
// --eval takes a JavaScript expression (<expr>).
// Host flags collapse to <host>, port to N, auth/credential flags to <val>.
// TLS/SSL certificate path flags collapse to <path>.
func handleMongo(_ string, tokens []string) []string {
	args, redirects := shellshape.SplitRedirects(tokens)

	hostFlags := map[string]bool{
		"-h": true, "--host": true,
	}

	numericFlags := map[string]bool{
		"--port": true,
	}

	evalFlags := map[string]bool{
		"--eval": true,
	}

	pathFlags := map[string]bool{
		"-f": true, "--file": true,
		"--tlsCertificateKeyFile": true, "--tlsCAFile": true,
		"--tlsCRLFile": true, "--tlsCertificateSelector": true,
		"--sslCAFile": true, "--sslCRLFile": true,
		"--sslPEMKeyFile": true,
		"--keyFile":       true,
	}

	valFlags := map[string]bool{
		"-u": true, "--username": true,
		"-p": true, "--password": true,
		"--authenticationDatabase":  true,
		"--authenticationMechanism": true,
		"--gssapiServiceName":       true, "--gssapiHostName": true,
		"--awsIamSessionToken":            true,
		"--tlsCertificateKeyFilePassword": true,
		"--sslPEMKeyPassword":             true,
		"--networkMessageCompressors":     true,
		"--apiVersion":                    true,
	}

	categories := []shellshape.FlagCategory{
		{Flags: hostFlags, Placeholder: "<host>"},
		{Flags: numericFlags, Placeholder: "N"},
		{Flags: evalFlags, Placeholder: "<expr>"},
		{Flags: pathFlags, Placeholder: "<path>"},
		{Flags: valFlags, Placeholder: "<val>"},
	}

	var result []string
	dbAssigned := false
	scriptAssigned := false

	i := 0
	for i < len(args) {
		tok := args[i]

		if shellshape.IsSubshellToken(tok) {
			result = append(result, tok)
			if !dbAssigned {
				dbAssigned = true
			} else {
				scriptAssigned = true
			}
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

		// Positionals: first is database/connection string, second is script file.
		if !dbAssigned {
			if strings.HasPrefix(tok, "mongodb://") || strings.HasPrefix(tok, "mongodb+srv://") {
				result = append(result, "<mongodb-uri>")
			} else {
				result = append(result, "<db>")
			}
			dbAssigned = true
		} else if !scriptAssigned {
			result = append(result, "<path>")
			scriptAssigned = true
		} else {
			result = append(result, shellshape.ClassifyToken(tok))
		}
		i++
	}

	result = append(result, redirects...)
	return result
}
