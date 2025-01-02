package internal

import (
	"encoding/binary"
	"fmt"
	"github.com/maheshkhanwalkar/dysm/pkg/obj/machO"
	"os"
)

func LoadObjectFile(name string) error {
	f, err := os.Open(name)
	machOHeader := machO.Header{}

	if err != nil {
		return err
	}
	defer f.Close()

	if err = binary.Read(f, binary.LittleEndian, &machOHeader); err != nil {
		return err
	}

	if err = validateMagic(machOHeader.Magic); err != nil {
		return err
	}

	fmt.Printf("Mach-O magic: 0x%x\n", machOHeader.Magic)
	fmt.Printf("Architecture: %s\n", machO.CpuArchName(machOHeader.CpuType))
	fmt.Printf("Number of load commands: %d\n", machOHeader.NumLoadCmd)

	loadCmd := machO.LoadCmd{}

	if err = binary.Read(f, binary.LittleEndian, &loadCmd); err != nil {
		return err
	}

	if loadCmd.CmdType == machO.SegmentLoad64 {
		_ = parseSectionLoadCommand(f)
	}

	return nil
}

func validateMagic(magic uint32) error {
	if magic != machO.Magic64 {
		return fmt.Errorf("invalid Mach-O magic: 0x%x", magic)
	}

	return nil
}

func parseSectionLoadCommand(f *os.File) error {
	sectionLoadCmd := machO.SegmentLoadCmd64{}
	if err := binary.Read(f, binary.LittleEndian, &sectionLoadCmd); err != nil {
		return err
	}

	fmt.Printf("Segment name: %s\n", machO.GetSegmentName(&sectionLoadCmd))
	return nil
}
