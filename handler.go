package shellshape

// handlerFunc processes tokens after the executable name.
// subcommand is the subcommand token (empty if none).
type handlerFunc func(subcommand string, tokens []string) []string

// handler registry - maps executable names to their handler
var handlers = map[string]handlerFunc{
	// text output
	"echo":   handleEcho,
	"printf": handleEcho,

	// text processing
	"sed":   handleSed,
	"grep":  handleGrep,
	"egrep": handleGrep,
	"fgrep": handleGrep,
	"rg":    handleGrep,
	"ag":    handleGrep,
	"ack":   handleGrep,
	"awk":   handleAwk,
	"gawk":  handleAwk,
	"mawk":  handleAwk,
	"nawk":  handleAwk,
	"jq":    handleJq,
	"yq":    handleYq,
	"tr":    handleTr,
	"cut":   handleCut,
	"sort":  handleSort,
	"uniq":  handleUniq,
	"wc":    handleWc,
	"nl":    handleNl,
	"comm":  handleComm,
	"paste": handlePaste,
	"column": handleColumn,
	"envsubst": handleEnvsubst,
	"strings": handleStrings,

	// file viewing
	"cat":  handleCat,
	"head": handleTail,
	"tail": handleTail,
	"tee":  handleTee,
	"diff": handleDiff,
	"xxd":  handleXxd,

	// file operations
	"find":  handleFind,
	"ls":    handleLs,
	"cp":    handleCp,
	"mv":    handleMv,
	"rm":    handleRm,
	"unlink": handleRm,
	"ln":    handleLn,
	"link":  handleLn,
	"mkdir": handleMkdir,
	"touch": handleTouch,
	"file":  handleFile,
	"stat":  handleStat,

	// permissions/ownership
	"chmod": handleChmod,
	"chown": handleChown,

	// disk/compression
	"du":     handleDu,
	"df":     handleDf,
	"dd":     handleDd,
	"tar":    handleTar,
	"gzip":   handleGzip,
	"gunzip": handleGzip,
	"zcat":   handleGzip,
	"zip":    handleZip,
	"unzip":  handleUnzip,
	"base64":    handleBase64,
	"b64encode": handleBase64,
	"b64decode": handleBase64,

	// network/dns
	"curl":       handleCurl,
	"wget":       handleWget,
	"ssh":        handleSsh,
	"scp":        handleScp,
	"rsync":      handleRsync,
	"ping":       handlePing,
	"ping6":      handlePing,
	"dig":        handleDig,
	"host":       handleHost,
	"whois":      handleWhois,
	"nc":         handleNc,
	"netcat":     handleNc,
	"ncat":       handleNc,
	"traceroute": handleTraceroute,
	"tracepath":  handleTraceroute,

	// process management
	"ps":   handlePs,
	"kill": handleKill,
	"lsof": handleLsof,

	// crypto/hashing
	"openssl": handleOpenssl,
	// shasum/md5sum/sha*sum registered via init() in handler_shasum.go

	// databases
	"sqlite3": handleSqlite3,
	"sqlite":  handleSqlite3,

	// documentation
	"man":     handleMan,
	"apropos": handleMan,
	"whatis":  handleMan,

	// JS/TS toolchain
	"tsc":       handleTsc,
	"tsgo":      handleTsc,
	"tsx":       handleTsx,
	"node":      handleNode,
	"eslint":    handleEslint,
	"biome":     handleBiome,
	"vitest":    handleVitest,
	"prettier":  handlePrettier,
	"npx":       handleNpx,
	"bunx":      handleNpx,
	"pnpx":      handleNpx,

	// package managers
	"pnpm": handlePnpm,
	"npm":  handlePnpm,
	"yarn": handlePnpm,
	"brew": handleBrew,

	// env management
	"dotenvx": handleDotenvx,
	"dotenv":  handleDotenvx,

	// infrastructure/cloud
	"docker":    handleDocker,
	"aws":       handleAws,
	"gh":        handleGh,
	"sst":       handleSst,
	"kubectl":   handleKubectl,
	"terraform": handleTerraform,
	"helm":      handleHelm,

	// build tools
	"make":  handleMake,
	"xargs": handleXargs,

	// media
	"ffmpeg":  handleFfmpeg,
	"ffprobe": handleFfmpeg,

	// version managers
	"fnm": handleFnm,

	// system utilities
	"stow":  handleStow,
	"watch": handleWatch,
	"tmux":  handleTmux,

	// runtime/subcommand
	"bun": handleBun,
}

var redirectConsumeNext = map[string]bool{
	">": true, ">>": true, "<": true,
	"&>": true, "&>>": true,
	"2>": true, "2>>": true, "1>": true,
}

var redirectStandalone = map[string]bool{
	"2>&1": true, "1>&2": true,
	"&>/dev/null": true, "2>/dev/null": true,
}

// splitRedirects partitions tokens into non-redirect args and redirect tail.
func splitRedirects(tokens []string) (args, redirects []string) {
	i := 0
	for i < len(tokens) {
		tok := tokens[i]

		if redirectStandalone[tok] {
			redirects = append(redirects, tok)
			i++
			continue
		}

		if redirectConsumeNext[tok] {
			if i+1 < len(tokens) {
				redirects = append(redirects, tok, "<path>")
				i += 2
			} else {
				redirects = append(redirects, tok)
				i++
			}
			continue
		}

		args = append(args, tok)
		i++
	}
	return
}
