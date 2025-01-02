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

	err = binary.Read(f, binary.LittleEndian, &machOHeader)
	if err != nil {
		return err
	}

	if err = validateMagic(machOHeader.Magic); err != nil {
		return err
	}

	fmt.Printf("Mach-O magic: 0x%x\n", machOHeader.Magic)
	fmt.Printf("Architecture: %s", machO.CpuType(machOHeader.CpuType))

	return nil
}

func validateMagic(magic uint32) error {
	if magic != machO.Magic32 && magic != machO.Magic64 {
		return fmt.Errorf("invalid Mach-O magic: 0x%x", magic)
	}

	return nil
}
