package machO

type cpuType uint32
type fileType uint32
type magicType uint32

// CPU type constants
// There are other possible values, but these are the only two mainstream ones
const (
	X86 cpuType = 0x6
	ARM cpuType = 0xC
)

// File Type constants
// TODO - check if there are other values used frequently
const (
	RelocatableObjectFile     fileType = 0x1
	DemandPagedExecutableFile fileType = 0x2
)

// Magic Value constants
const (
	Magic32 magicType = 0xfeedface
	Magic64 magicType = 0xfeedfecf
)

// Header is the on-disk representation of the Mach-O header for Mac object files
type Header struct {
	Magic       uint32
	CpuType     cpuType
	CpuSubType  uint32
	FileType    fileType
	NumLoadCmd  uint32
	SizeLoadCmd uint32
	Flags       uint32
	Reserved    uint32
}
