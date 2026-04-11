package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	shellshape "github.com/mthiesen/shellshape"
)

func main() {
	if len(os.Args) > 1 {
		input := strings.Join(os.Args[1:], " ")
		fmt.Println(shellshape.Normalize(input))
		return
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
