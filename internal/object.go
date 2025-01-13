package internal

import (
	"encoding/binary"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/maheshkhanwalkar/dysm/pkg/arch"
	"github.com/maheshkhanwalkar/dysm/pkg/arch/arm64"
	"github.com/maheshkhanwalkar/dysm/pkg/arch/x86_64"
	"github.com/maheshkhanwalkar/dysm/pkg/obj/machO"
	"os"
	"regexp"
)

const (
	hexDumpRowBytes uint64 = 0x10
)

// ObjectFile represents a generic object file (format agnostic) interface that defines function signatures
// to print out information about the file
type ObjectFile interface {
	PrintHeaders()
	PrintSymbolTable()
	DumpSections(sectionNameRegex string) error
	PrintDisassembly() error
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

// DumpSections prints out a hex-dump of all sections that match the given section name regex
// The section name refers to a Mach-O section (not segment)
// Returns an error if there is no section that matches the regex in the object file
func (m *MachObjectFile) DumpSections(sectionNameRegex string) error {
	sections, err := m.findSections(sectionNameRegex)
	if err != nil {
		return err
	}

	for _, section := range sections {
		data := section.Data
		rounds := uint64(len(data)) / hexDumpRowBytes

		// Need an extra partial round -- if not cleanly divisible by hexDumpRowBytes
		if uint64(len(data))%hexDumpRowBytes != 0 {
			rounds += 1
		}

		fmt.Printf("Section %s\n\n", section.Name)

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

// PrintDisassembly prints out the disassembled machine code from the __text section out
// as architecture-specific assembly language along with corresponding addresses.
func (m *MachObjectFile) PrintDisassembly() error {
	sections, err := m.findSections("__text")
	if err != nil {
		return err
	}

	textSection := sections[0]
	var instructions []arch.Inst

	switch m.Obj.CpuArchName() {
	case "ARM64":
		arm64Insts, err := arm64.Disassemble(textSection.Data, textSection.Address)
		if err != nil {
			return err
		}
		instructions = toInstList(arm64Insts)
	case "x86_64":
		x86Insts, err := x86_64.Disassemble(textSection.Data, textSection.Address)
		if err != nil {
			return err
		}
		instructions = toInstList(x86Insts)
	}

	for _, instruction := range instructions {
		fmt.Println(instruction.Asm())
	}

	return nil
}

func (m *MachObjectFile) findSections(sectionNameRegex string) ([]machO.Section, error) {
	var matched = make([]machO.Section, 0)

	// Assumption: a section name is unique across segments -- otherwise, this will return
	// the first section that matches the given name
	for _, seg := range m.Obj.Segments {
		for _, section := range seg.Sections {
			didMatch, err := regexp.MatchString("^"+sectionNameRegex+"$", section.Name)
			if err != nil {
				return nil, err
			}

			if didMatch {
				matched = append(matched, section)
			}
		}
	}

	if len(matched) == 0 {
		return nil, fmt.Errorf("no section found matching %s", sectionNameRegex)
	} else {
		return matched, nil
	}
}

func clumpBytes(arr []byte) []uint16 {
	res := make([]uint16, 0, len(arr)/2+1)

	for i := 0; i < len(arr); i += 2 {
		group := arr[i:min(i+2, len(arr))]
		if len(group) == 1 {
			group = []byte{group[0], 0}
		}

		u16 := binary.LittleEndian.Uint16(group)
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

func toInstList[T arch.Inst](arr []T) []arch.Inst {
	insts := make([]arch.Inst, 0, len(arr))
	for _, inst := range arr {
		insts = append(insts, arch.Inst(inst))
	}
	return insts
}
