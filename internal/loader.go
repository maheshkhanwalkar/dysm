package internal

import (
	"github.com/maheshkhanwalkar/dysm/pkg/obj/machO"
)

func LoadMachObjectFile(name string) (*MachObjectFile, error) {
	mach, err := machO.FromFile(name)
	if err != nil {
		return nil, err
	}

	objFile := NewMachObjectFile(mach, name)
	return objFile, nil
}
