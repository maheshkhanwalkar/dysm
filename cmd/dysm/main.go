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
	flag.Parse()

	objFileName := flag.Arg(0)
	obj, err := internal.LoadMachObjectFile(objFileName)

	if err != nil {
		die(err.Error())
		return
	}

	if *hdr {
		fmt.Println("File: " + objFileName)
		fmt.Printf("Architecture: %s\n\n", obj.CpuArchName())

		for _, seg := range obj.Segments {
			fmt.Printf("Segment name: %s, Start Address: 0x%x, End Address: 0x%x, Permission: %s\n", seg.Name,
				seg.Address, seg.Address+seg.AddressSize, seg.Permission)

			for i, section := range seg.Sections {
				fmt.Printf("\t%d) Section name: %s, Start Address: 0x%x, End Address: 0x%x, Size: %d\n", i+1,
					section.Name, section.Address, section.Address+section.Size, section.Size)
			}

			if len(seg.Sections) == 0 {
				fmt.Println("\tNo sections in this segment")
			}
		}
	}
}

func die(msg string) {
	_, _ = fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
