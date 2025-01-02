package internal

import (
	"github.com/maheshkhanwalkar/dysm/pkg/obj/machO"
)

func LoadMachObjectFile(name string) (*machO.MachO, error) {
	mach, err := machO.FromFile(name)
	if err != nil {
		return nil, err
	}
	return mach, nil
}
