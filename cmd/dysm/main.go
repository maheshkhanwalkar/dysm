package main

import (
	"fmt"
	"github.com/maheshkhanwalkar/dysm/internal"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		die("no arguments provided")
	}

	objFile := os.Args[1]
	_ = internal.LoadObjectFile(objFile)
}

func die(msg string) {
	_, _ = fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
