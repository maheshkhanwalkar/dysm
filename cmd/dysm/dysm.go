package main

import (
	"flag"
	"fmt"
	"github.com/maheshkhanwalkar/dysm/internal"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		die("no arguments provided")
	}

	hdr := flag.Bool("hdr", false, "print out section information for the file")
	dump := flag.String("dump", "", "hex dump of specified section")
	symtab := flag.Bool("symtab", false, "print out symbol table information")

	flag.Parse()

	objFileName := flag.Arg(0)
	var obj internal.ObjectFile

	obj, err := internal.LoadMachObjectFile(objFileName)

	if err != nil {
		die(err.Error())
		return
	}

	if *hdr {
		obj.PrintHeaders()
	}

	if *symtab {
		obj.PrintSymbolTable()
	}

	if *dump != "" {
		err = obj.DumpSection(*dump)
		if err != nil {
			die(err.Error())
		}
	}
}

func die(msg string) {
	_, _ = fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
