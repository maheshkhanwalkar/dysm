package raw

import "bytes"

// CPU type constants
// There are other possible values, but these are the only mainstream ones
const (
	X86_64 uint32 = 0x1000007
	ARM64  uint32 = 0x100000C
)

// File Type constants
const (
	RelocatableObjectFile     uint32 = 0x1
	DemandPagedExecutableFile uint32 = 0x2
)

const (
	Magic64 uint32 = 0xfeedfacf
)

// Header is the on-disk representation of the Mach-O header for Mac object files
type Header struct {
	Magic       uint32
	CpuType     uint32
	CpuSubType  uint32
	FileType    uint32
	NumLoadCmd  uint32
	SizeLoadCmd uint32
	Flags       uint32
	Reserved    uint32
}

const (
	SymbolTable   uint32 = 0x2
	SegmentLoad64 uint32 = 0x19
)

// LoadCmd is a load command entry, which follows the Mach-O header
type LoadCmd struct {
	CmdType uint32
	CmdSize uint32
}

const (
	MemRead  uint32 = 0x1
	MemWrite uint32 = 0x2
	MemExec  uint32 = 0x4
)

type SegmentLoadCmd64 struct {
	SegmentName [16]byte
	Address     uint64
	AddressSize uint64
	FileOffset  uint64
	FileSize    uint64
	MaxMemProt  uint32
	InitMemProt uint32
	NumSections uint32
	Flag32      uint32
}

type Section64 struct {
	SectionName     [16]byte
	SegmentName     [16]byte
	Address         uint64
	Size            uint64
	FileOffset      uint32
	Alignment       uint32
	RelocFileOffset uint32
	NumRelocations  uint32
	Flag            uint32
	Reserved1       uint32
	Reserved2       uint32
	Reserved3       uint32
}

type SymbolTableLoadCmd struct {
	SymbolFileOffset  uint32
	NumSymbols        uint32
	StringTableOffset uint32
	StringTableSize   uint32
}

const (
	SymbolUndefined uint8 = 0x0
	SymbolExternal  uint8 = 0x1
	SymbolSector    uint8 = 0xe
)

type Symbol64 struct {
	NameOffset    uint32
	SymbolType    uint8
	SectionNumber uint8
	DataInfo      uint16
	SymbolAddress uint64
}

// GetName converts a NULL-padded byte array to a Go string
func GetName(nameBytes []byte) string {
	// The raw name is a NULL-padded string, so we need to find the position of the first NULL
	// and take the slice just before that point and convert it to a Go string
	n := bytes.IndexByte(nameBytes, 0)
	var name []byte

	if n >= 0 {
		name = nameBytes[:n]
	} else {
		// Somehow there is no NULL, so just use the entire byte array
		name = nameBytes
	}

	return string(name)
}

// GetPermissionString returns a string representation of the given permission value
// The format is as follows:
//  1. 'r' (read-allowed)  or '-' (read not allowed)
//  2. 'w' (write-allowed) or '-' (write not allowed)
//  3. 'x' (exec-allowed)  or '-' (exec not allowed)
//
// Together, they made a 3-digit string
//
//	e.g. r-x means read-execute allowed, but write is not
//	     rwx means read-write-execute are allowed
//	     --- means nothing is allowed
//
// As shown above, the length of the string returned is always 3 digits long,
// with any unsupported (denied) modes represented with a '-' character.
func GetPermissionString(perm uint32) string {
	r := getPermDigit(perm, MemRead, "r")
	w := getPermDigit(perm, MemWrite, "w")
	x := getPermDigit(perm, MemExec, "x")
	return r + w + x
}

func getPermDigit(perm uint32, flag uint32, ch string) string {
	if perm&flag != 0 {
		return ch
	} else {
		return "-"
	}
}
