package internal

import (
	"encoding/binary"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/maheshkhanwalkar/dysm/pkg/obj/machO"
	"os"
)

const (
	hexDumpRowBytes uint64 = 0x10
)

// ObjectFile represents a generic object file (format agnostic) interface that defines function signatures
// to print out information about the file
type ObjectFile interface {
	PrintHeaders()
	PrintSymbolTable()
	DumpSection(sectionName string) error
}

// MachObjectFile represents a Mach-O object file and implements the ObjectFile interface
type MachObjectFile struct {
	Obj      *machO.MachO
	FileName string
}

func NewMachObjectFile(obj *machO.MachO, fileName string) *MachObjectFile {
	return &MachObjectFile{obj, fileName}
}

// PrintHeaders prints out Mach-O header meta-data
func (m *MachObjectFile) PrintHeaders() {
	fmt.Println("File: " + m.FileName)
	fmt.Printf("Architecture: %s\n\n", m.Obj.CpuArchName())

	for _, seg := range m.Obj.Segments {
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

// DumpSection prints out a hex-dump of the given section
// The section name refers to a Mach-O section (not segment)
// Returns an error if there is no section with the given name in the object file
func (m *MachObjectFile) DumpSection(sectionName string) error {
	section, err := m.findSection(sectionName)
	if err != nil {
		return err
	}

	data := section.Data
	rounds := uint64(len(data)) / hexDumpRowBytes

	fmt.Printf("Section %s\n\n", sectionName)

	for i := uint64(0); i < rounds; i++ {
		end := min((i+1)*hexDumpRowBytes, uint64(len(data)))
		row := clumpBytes(data[i*hexDumpRowBytes : end])

		address := section.Address + i*hexDumpRowBytes
		fmt.Printf("%016x: ", address)

		for j, u16 := range row {
			fmt.Printf("%04x", u16)
			if j != len(row)-1 {
				fmt.Printf(" ")
			}
		}

		fmt.Printf("\n")
	}

	return nil
}

// PrintSymbolTable prints out the symbol table information from the Mach-O file
// in a tabular format
func (m *MachObjectFile) PrintSymbolTable() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.SetTitle("Mach-O Symbol Table")
	t.AppendHeader(table.Row{"#", "Name", "Address", "Type", "Section"})

	const MaxStringLength = 32

	for i, sym := range m.Obj.Symbols {
		sectionName := "None"
		if sym.Section != nil {
			sectionName = trimString(sym.Section.Name, MaxStringLength)
		}

		symName := trimString(sym.Name, MaxStringLength)

		t.AppendRow(table.Row{
			i + 1,
			symName,
			fmt.Sprintf("%x", sym.Address),
			getTypeString(&sym),
			sectionName,
		})
	}

	t.SetStyle(table.StyleLight)
	t.Render()
}

func (m *MachObjectFile) findSection(sectionName string) (*machO.Section, error) {
	var matched *machO.Section

	// Assumption: a section name is unique across segments -- otherwise, this will return
	// the first section that matches the given name
	for _, seg := range m.Obj.Segments {
		for _, section := range seg.Sections {
			if section.Name == sectionName {
				matched = &section
				return matched, nil
			}
		}
	}

	return nil, fmt.Errorf("section %s not found", sectionName)
}

func clumpBytes(arr []byte) []uint16 {
	res := make([]uint16, 0, len(arr)/2+1)

	for i := 0; i < len(arr); i += 2 {
		u16 := binary.LittleEndian.Uint16(arr[i:min(i+2, len(arr))])
		res = append(res, u16)
	}

	return res
}

func getTypeString(sym *machO.Symbol) string {
	var typeString string

	if sym.IsExported() {
		typeString += "External"
	} else {
		typeString += "Internal (Not Exported)"
	}

	if sym.IsUndefined() {
		typeString += ", Needs linker resolution"
	}

	if typeString == "" {
		typeString = "Unknown"
	}

	return typeString
}

func trimString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	return s[:maxLen-3] + "..."
}
