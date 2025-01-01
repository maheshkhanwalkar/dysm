package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		die("no arguments provided")
	}
}

func die(msg string) {
	_, _ = fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
