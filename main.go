package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: kipf <command> [args...]")
		os.Exit(1)
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	var code int
	switch cmd {
	case "git":
		code = gitHandler(args)
	default:
		code = passthrough(cmd, args)
	}
	os.Exit(code)
}
