package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	shellshape "github.com/micthiesen/shellshape"
	_ "github.com/micthiesen/shellshape/handlers"
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "version" {
			fmt.Println(versionString())
			return
		}
		if os.Args[1] == "list" {
			for _, name := range shellshape.RegisteredHandlers() {
				fmt.Println(name)
			}
			return
		}
		input := strings.Join(os.Args[1:], " ")
		fmt.Println(shellshape.Normalize(input))
		return
	}

	// If stdin is a terminal, print usage instead of hanging.
	stat, _ := os.Stdin.Stat()
	if stat.Mode()&os.ModeCharDevice != 0 {
		fmt.Fprintln(os.Stderr, "Usage: shellshape <command>")
		fmt.Fprintln(os.Stderr, "       echo <command> | shellshape")
		fmt.Fprintln(os.Stderr, "       shellshape list")
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			fmt.Println()
			continue
		}
		fmt.Println(shellshape.Normalize(line))
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "shellshape: %v\n", err)
		os.Exit(1)
	}
}

func versionString() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "shellshape (unknown)"
	}
	var revision, time string
	var dirty bool
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			revision = s.Value
		case "vcs.time":
			time = s.Value
		case "vcs.modified":
			dirty = s.Value == "true"
		}
	}
	if revision == "" {
		return "shellshape (dev)"
	}
	short := revision
	if len(short) > 7 {
		short = short[:7]
	}
	v := "shellshape " + short + " (" + time + ")"
	if dirty {
		v += " dirty"
	}
	return v
}
