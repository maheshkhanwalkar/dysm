package arm64

import (
	"errors"
	"fmt"
	"golang.org/x/arch/arm64/arm64asm"
	"strings"
)

type Instruction struct {
	inst *arm64asm.Inst
	addr uint64
}

func (inst Instruction) Asm() string {
	asm := strings.ToLower(inst.inst.String())
	return fmt.Sprintf("%016x: %s", inst.addr, asm)
}

type instWord struct {
	// instBytes is always 4-bytes long, as ARM64 has fixed size instructions
	instBytes []byte
}

func Disassemble(code []byte, startAddr uint64) ([]Instruction, error) {
	words, err := groupIntoInstWords(code)
	if err != nil {
		return nil, err
	}

	instructions := make([]Instruction, 0, len(words))
	address := startAddr

	for _, word := range words {
		inst, _ := arm64asm.Decode(word.instBytes)
		instructions = append(instructions, Instruction{addr: address, inst: &inst})
		address += 4
	}

	return instructions, nil
}

func groupIntoInstWords(code []byte) ([]instWord, error) {
	if len(code)%4 != 0 {
		return nil,
			errors.New("code section must be a multiple of 4 bytes, as all ARM64 instructions are 4 bytes long")
	}

	instructions := make([]instWord, 0, len(code)/4)

	for i := 0; i < len(code); i += 4 {
		instructions = append(instructions, instWord{instBytes: code[i : i+4]})
	}

	return instructions, nil
}
