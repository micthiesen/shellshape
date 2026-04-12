package shellshape

import "sort"

// HandlerOptions configures how a command is normalized.
type HandlerOptions struct {
	// HasSubcommands indicates the second positional token is a subcommand
	// (e.g. "docker run", "git commit"). The subcommand is extracted and
	// passed to the handler before remaining tokens.
	HasSubcommands bool
}

var handlers = map[string]HandlerFunc{}
var handlerOptions = map[string]HandlerOptions{}

// Register adds a command handler to the registry.
// If fn is nil the command is registered for its options only (e.g. subcommand
// detection) without a per-command handler.
func Register(name string, fn HandlerFunc, opts ...HandlerOptions) {
	if fn != nil {
		handlers[name] = fn
	}
	if len(opts) > 0 {
		handlerOptions[name] = opts[0]
	}
}

func hasSubcommands(exe string) bool {
	return handlerOptions[exe].HasSubcommands
}

// RegisteredHandlers returns a sorted list of all registered command names.
func RegisteredHandlers() []string {
	seen := make(map[string]bool)
	for name := range handlers {
		seen[name] = true
	}
	for name := range handlerOptions {
		seen[name] = true
	}
	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func init() {
	// Commands that need subcommand detection but have no per-command handler.
	for _, name := range []string{
		"git", "cargo", "gcloud", "az",
		"pip", "pip3", "uv", "poetry",
		"snap",
		"launchctl",
		"pulumi", "ollama", "lms",
	} {
		Register(name, nil, HandlerOptions{HasSubcommands: true})
	}
}
