package machO

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/maheshkhanwalkar/dysm/pkg/obj/machO/raw"
	"io"
	"os"
)

// MachO top level structure of Mach-O contents
type MachO struct {
	hdr      *raw.Header
	Segments []Segment
}

// Segment represents a Mach-O segment
type Segment struct {
	Name        string
	Address     uint64
	AddressSize uint64
	Sections    []Section
}

// Section represents a Mach-O section
// TODO actually add relevant fields
type Section struct {
}

// Test whether magic is a Mach-O magic value
func Test(magic uint32) bool {
	return magic == raw.Magic64
}

// FromSlice reads a Mach-O file from the given byte slice
func FromSlice(arr []byte) (*MachO, error) {
	r := bytes.NewReader(arr)

	hdr, err := readHeader(r)
	if err != nil {
		return nil, err
	}

	segments, err := readSegments(r, hdr)
	if err != nil {
		return nil, err
	}

	return &MachO{hdr: hdr, Segments: segments}, nil
}

// FromFile reads a Mach-O file from the given file
func FromFile(name string) (*MachO, error) {
	arr, err := os.ReadFile(name)
	if err != nil {
		return nil, err
	}
	return FromSlice(arr)
}

// CpuArchName returns the Mach-O file's architecture name
func (m *MachO) CpuArchName() string {
	switch m.hdr.CpuType {
	case raw.X86_64:
		return "x86_64"
	case raw.ARM64:
		return "ARM64"
	default:
		return "unknown"
	}
}

func readHeader(r io.Reader) (*raw.Header, error) {
	var hdr raw.Header
	if err := binary.Read(r, binary.LittleEndian, &hdr); err != nil {
		return nil, err
	}

	if !Test(hdr.Magic) {
		magic := fmt.Sprintf("0x%x", hdr.Magic)
		return nil, &InvalidMachOMagicError{InvalidMagic: magic}
	}

	return &hdr, nil
}

func readSegments(r io.Reader, hdr *raw.Header) ([]Segment, error) {
	segments := make([]Segment, 0, hdr.NumLoadCmd)

	for i := uint32(0); i < hdr.NumLoadCmd; i++ {
		var loadCmd raw.LoadCmd
		if err := binary.Read(r, binary.LittleEndian, &loadCmd); err != nil {
			return nil, err
		}

		// Ignore this command, but read 'til the end, so we can correctly process the
		// next segment in the file
		if loadCmd.CmdType != raw.SegmentLoad64 {
			ign := make([]byte, loadCmd.CmdSize-8)
			if _, err := r.Read(ign); err != nil {
				return nil, err
			}
			continue
		}

		var segment raw.SegmentLoadCmd64
		if err := binary.Read(r, binary.LittleEndian, &segment); err != nil {
			return nil, err
		}

		seg := Segment{
			Name:        raw.GetSegmentName(&segment),
			Address:     segment.Address,
			AddressSize: segment.AddressSize,
			Sections:    make([]Section, segment.NumSections),
		}

		segments = append(segments, seg)

		// TODO actually read the sections -- instead of just skipping over them...
		const sectionSize = 80
		ign := make([]byte, sectionSize*segment.NumSections)
		if n, err := r.Read(ign); err != nil {
			if n != len(ign) {
				return nil, err
			}
		}
	}

	return segments, nil
}
