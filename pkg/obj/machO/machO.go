package machO

type cpuType uint32
type fileType uint32

// CPU type constants
// There are other possible values, but these are the only mainstream ones
const (
	X86_32 cpuType = 0x7
	ARM32  cpuType = 0xC
	X86_64 cpuType = 0x1000007
	ARM64  cpuType = 0x100000C
)

// File Type constants
// TODO - check if there are other values used frequently
const (
	RelocatableObjectFile     fileType = 0x1
	DemandPagedExecutableFile fileType = 0x2
)

// Magic Value constants
const (
	Magic32 uint32 = 0xfeedface
	Magic64 uint32 = 0xfeedfacf
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

func CpuType(cpuType cpuType) string {
	switch cpuType {
	case X86_32:
		return "x86_32"
	case ARM32:
		return "ARM32"
	case X86_64:
		return "x86_64"
	case ARM64:
		return "ARM64"
	default:
		return "unknown"
	}
}
