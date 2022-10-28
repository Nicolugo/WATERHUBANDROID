
package cpu

import (
	"github.com/bokuweb/gopher-boy/pkg/interfaces/bus"
	"github.com/bokuweb/gopher-boy/pkg/interfaces/interrupt"
	"github.com/bokuweb/gopher-boy/pkg/interfaces/logger"
	"github.com/bokuweb/gopher-boy/pkg/types"
	"github.com/bokuweb/gopher-boy/pkg/utils"
)

// Registers is generic cpu registers
type Registers struct {
	A types.Register
	B types.Register
	C types.Register
	D types.Register
	E types.Register
	H types.Register
	L types.Register
	F types.Register
}

type flags int

const (
	// C is carry flag
	C flags = iota + 1
	// H is half carry flag
	H
	// N is negative flag
	N
	// Z is zero flag
	Z
)

// CPU is cpu state struct
type CPU struct {
	logger  logger.Logger
	PC      types.Word
	SP      types.Word
	Regs    Registers
	bus     bus.Accessor
	irq     interrupt.Interrupt
	stopped bool
	halted  bool
}

type Cycle = uint

// NewCPU is CPU constructor
func NewCPU(logger logger.Logger, bus bus.Accessor, irq interrupt.Interrupt) *CPU {
	cpu := &CPU{
		logger: logger,
		PC:     0x100, // INFO: Skip
		SP:     0xFFFE,
		Regs: Registers{
			A: 0x11,
			B: 0x00,
			C: 0x00,
			D: 0xFF,
			E: 0x56,
			F: 0x80,
			H: 0x00,
			L: 0x0D,
		},
		bus:     bus,
		irq:     irq,
		stopped: false,
		halted:  false,
	}
	return cpu
}

func (cpu *CPU) fetch() byte {
	d := cpu.bus.ReadByte(cpu.PC)
	cpu.PC++
	return d
}

// Step execute an instruction
func (cpu *CPU) Step() Cycle {

	if cpu.halted {
		if cpu.irq.HasIRQ() {
			cpu.halted = false
		}
		return 0x01
	}
	// cpc := cpu.PC
	if hasIRQ := cpu.resolveIRQ(); hasIRQ {
		return 0x01
	}
	opcode := cpu.fetch()
	var inst *inst
	if opcode == 0xCB {
		next := cpu.fetch()
		inst = cbPrefixedInstructions[next]
	} else {
		inst = instructions[opcode]
	}

	operands := cpu.fetchOperands(inst.OperandsSize)
	// cpu.logger.Info(fmt.Sprintf("PC = %X Opcode = %X %+v %+v %+v", cpc, opcode, cpu.Regs, inst, operands))
	inst.Execute(cpu, operands)
	return inst.Cycles
}

func (cpu *CPU) fetchOperands(size uint) []byte {
	operands := []byte{}
	for i := 0; i < int(size); i++ {
		operands = append(operands, cpu.fetch())
	}
	return operands
}

type inst struct {
	Opcode       byte
	Description  string
	OperandsSize uint
	Cycles       uint
	Execute      func(cpu *CPU, operands []byte)
}

// EMPTY is empty instruction
var EMPTY = &inst{0xFF, "EMPTY", 0, 1, func(cpu *CPU, operands []byte) {
}}

var cbPrefixedInstructions = []*inst{
	&inst{0x0, "RLC B", 0, 2, func(cpu *CPU, operands []byte) { cpu.rlc_n(&cpu.Regs.B) }},
	&inst{0x1, "RLC C", 0, 2, func(cpu *CPU, operands []byte) { cpu.rlc_n(&cpu.Regs.C) }},
	&inst{0x2, "RLC D", 0, 2, func(cpu *CPU, operands []byte) { cpu.rlc_n(&cpu.Regs.D) }},
	&inst{0x3, "RLC E", 0, 2, func(cpu *CPU, operands []byte) { cpu.rlc_n(&cpu.Regs.E) }},
	&inst{0x4, "RLC H", 0, 2, func(cpu *CPU, operands []byte) { cpu.rlc_n(&cpu.Regs.H) }},
	&inst{0x5, "RLC L", 0, 2, func(cpu *CPU, operands []byte) { cpu.rlc_n(&cpu.Regs.L) }},
	&inst{0x6, "RLC (HL)", 0, 4, func(cpu *CPU, operands []byte) { cpu.rlc_hl() }},
	&inst{0x7, "RLC A", 0, 2, func(cpu *CPU, operands []byte) { cpu.rlc_n(&cpu.Regs.A) }},
	&inst{0x8, "RRC B", 0, 2, func(cpu *CPU, operands []byte) { cpu.rrc_n(&cpu.Regs.B) }},
	&inst{0x9, "RRC C", 0, 2, func(cpu *CPU, operands []byte) { cpu.rrc_n(&cpu.Regs.C) }},
	&inst{0xA, "RRC D", 0, 2, func(cpu *CPU, operands []byte) { cpu.rrc_n(&cpu.Regs.D) }},
	&inst{0xB, "RRC E", 0, 2, func(cpu *CPU, operands []byte) { cpu.rrc_n(&cpu.Regs.E) }},