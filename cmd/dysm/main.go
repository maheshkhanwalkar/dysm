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

	objFileName := os.Args[1]
	obj, err := internal.LoadMachObjectFile(objFileName)

	if err != nil {
		die(err.Error())
		return
	}

	fmt.Println("File: " + objFileName)
	fmt.Printf("Architecture: %s\n", obj.CpuArchName())

	for _, seg := range obj.Segments {
		fmt.Println("Segment name: " + seg.Name)
	}
}

func die(msg string) {
	_, _ = fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
