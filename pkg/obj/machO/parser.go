package machO

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type MachO struct {
	hdr *Header
}

// Test whether magic is a Mach-O magic value
func Test(magic uint32) bool {
	return magic == Magic64
}

// FromSlice reads a Mach-O file from the given byte slice
func FromSlice(arr []byte) (*MachO, error) {
	r := bytes.NewReader(arr)

	hdr, err := readHeader(r)
	if err != nil {
		return nil, err
	}

	return &MachO{hdr: hdr}, nil
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
	case X86_64:
		return "x86_64"
	case ARM64:
		return "ARM64"
	default:
		return "unknown"
	}
}

func readHeader(r io.Reader) (*Header, error) {
	var hdr Header
	if err := binary.Read(r, binary.LittleEndian, &hdr); err != nil {
		return nil, err
	}

	if !Test(hdr.Magic) {
		magic := fmt.Sprintf("0x%x", hdr.Magic)
		return nil, &InvalidMachOMagicError{InvalidMagic: magic}
	}

	return &hdr, nil
}
