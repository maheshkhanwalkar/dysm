package machO

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
	SegmentLoad64 uint32 = 0x19
)

// LoadCmd is a load command entry, which follows the Mach-O header
type LoadCmd struct {
	CmdType uint32
	CmdSize uint32
}

type SegmentLoadCmd64 struct {
	SegmentName [16]byte
	Address     uint64
	AddressSize uint64
	FileOffset  uint64
	FileSize    uint64
	Unused1     uint32 // This is actually Maximum virtual memory protections -- but we don't care
	Unused2     uint32 // Initial virtual memory protections  [again, don't care]
	NumSections uint32
	Flag32      uint32
}

func GetSegmentName(seg *SegmentLoadCmd64) string {
	// The raw segment name is a NULL-padded string, so we need to find the position of the first NULL
	// and take the slice just before that point and convert it to a Go string
	n := bytes.IndexByte(seg.SegmentName[:], 0)
	var name []byte

	if n >= 0 {
		name = seg.SegmentName[:n]
	} else {
		// Somehow there is no NULL, so just use the entire byte array
		name = seg.SegmentName[:]
	}

	return string(name)
}
