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
	dumpAll := flag.Bool("dump-all", false, "hex dump all sections")
	symtab := flag.Bool("symtab", false, "print out symbol table information")

	flag.Parse()
	objFileName := flag.Arg(0)
	ensureExclusive(*hdr, *dump, *dumpAll, *symtab)

	var obj internal.ObjectFile
	obj, err := internal.LoadMachObjectFile(objFileName)

	if err != nil {
		die(err.Error())
		return
	}

	if *hdr {
		obj.PrintHeaders()
		return
	}

	if *symtab {
		obj.PrintSymbolTable()
		return
	}

	if *dump != "" {
		err = obj.DumpSections(*dump)
		if err != nil {
			die(err.Error())
		}
		return
	}

	if *dumpAll {
		err = obj.DumpSections(".*")
		if err != nil {
			die(err.Error())
		}
		return
	}

	// If nothing else is specified, we just assume disassembly as the
	// default action.
	err = obj.PrintDisassembly()
	if err != nil {
		die(err.Error())
	}
}

func die(msg string) {
	_, _ = fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func ensureExclusive(hdr bool, dump string, dumpAll bool, symtab bool) {
	count := 0

	if hdr {
		count++
	}
	if dump != "" {
		count++
	}
	if dumpAll {
		count++
	}
	if symtab {
		count++
	}

	if count > 1 {
		flag.Usage()
		die("too many arguments provided")
	}
}
