package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	shellshape "github.com/micthiesen/shellshape"
	_ "github.com/micthiesen/shellshape/handlers"
)

func main() {
	if len(os.Args) > 1 {
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
