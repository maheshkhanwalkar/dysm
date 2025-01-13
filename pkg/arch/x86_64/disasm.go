package x86_64

import (
	"fmt"
	"golang.org/x/arch/x86/x86asm"
	"strings"
)

type Instruction struct {
	inst *x86asm.Inst
	addr uint64
}

func (inst Instruction) Asm() string {
	var asm string
	if inst.inst.Opcode == 0 {
		asm = "<unknown>"
	} else {
		asm = strings.ToLower(inst.inst.String())
	}
	return fmt.Sprintf("%016x: %s", inst.addr, asm)
}

func Disassemble(code []byte, startAddr uint64) ([]Instruction, error) {
	instructions := make([]Instruction, 0)
	address := startAddr
	bytesProcessed := 0

	for bytesProcessed < len(code) {
		inst, _ := x86asm.Decode(code[bytesProcessed:], 64)
		instructions = append(instructions, Instruction{addr: address, inst: &inst})
		bytesProcessed += inst.Len
		address += uint64(inst.Len)
	}

	return instructions, nil
}
