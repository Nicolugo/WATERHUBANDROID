
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
	&inst{0xC, "RRC H", 0, 2, func(cpu *CPU, operands []byte) { cpu.rrc_n(&cpu.Regs.H) }},
	&inst{0xD, "RRC L", 0, 2, func(cpu *CPU, operands []byte) { cpu.rrc_n(&cpu.Regs.L) }},
	&inst{0xE, "RRC (HL)", 0, 4, func(cpu *CPU, operands []byte) { cpu.rrc_hl() }},
	&inst{0xF, "RRC A", 0, 2, func(cpu *CPU, operands []byte) { cpu.rrc_n(&cpu.Regs.A) }},
	&inst{0x10, "RL B", 0, 2, func(cpu *CPU, operands []byte) { cpu.rl_n(&cpu.Regs.B) }},
	&inst{0x11, "RL C", 0, 2, func(cpu *CPU, operands []byte) { cpu.rl_n(&cpu.Regs.C) }},
	&inst{0x12, "RL D", 0, 2, func(cpu *CPU, operands []byte) { cpu.rl_n(&cpu.Regs.D) }},
	&inst{0x13, "RL E", 0, 2, func(cpu *CPU, operands []byte) { cpu.rl_n(&cpu.Regs.E) }},
	&inst{0x14, "RL H", 0, 2, func(cpu *CPU, operands []byte) { cpu.rl_n(&cpu.Regs.H) }},
	&inst{0x15, "RL L", 0, 2, func(cpu *CPU, operands []byte) { cpu.rl_n(&cpu.Regs.L) }},
	&inst{0x16, "RL (HL)", 0, 4, func(cpu *CPU, operands []byte) { cpu.rl_hl() }},
	&inst{0x17, "RL A", 0, 2, func(cpu *CPU, operands []byte) { cpu.rl_n(&cpu.Regs.A) }},
	&inst{0x18, "RR B", 0, 2, func(cpu *CPU, operands []byte) { cpu.rr_n(&cpu.Regs.B) }},
	&inst{0x19, "RR C", 0, 2, func(cpu *CPU, operands []byte) { cpu.rr_n(&cpu.Regs.C) }},
	&inst{0x1A, "RR D", 0, 2, func(cpu *CPU, operands []byte) { cpu.rr_n(&cpu.Regs.D) }},
	&inst{0x1B, "RR E", 0, 2, func(cpu *CPU, operands []byte) { cpu.rr_n(&cpu.Regs.E) }},
	&inst{0x1C, "RR H", 0, 2, func(cpu *CPU, operands []byte) { cpu.rr_n(&cpu.Regs.H) }},
	&inst{0x1D, "RR L", 0, 2, func(cpu *CPU, operands []byte) { cpu.rr_n(&cpu.Regs.L) }},
	&inst{0x1E, "RR (HL)", 0, 4, func(cpu *CPU, operands []byte) { cpu.rr_hl() }},
	&inst{0x1F, "RR A", 0, 2, func(cpu *CPU, operands []byte) { cpu.rr_n(&cpu.Regs.A) }},
	&inst{0x20, "SLA B", 0, 2, func(cpu *CPU, operands []byte) { cpu.sla_n(&cpu.Regs.B) }},
	&inst{0x21, "SLA C", 0, 2, func(cpu *CPU, operands []byte) { cpu.sla_n(&cpu.Regs.C) }},
	&inst{0x22, "SLA D", 0, 2, func(cpu *CPU, operands []byte) { cpu.sla_n(&cpu.Regs.D) }},
	&inst{0x23, "SLA E", 0, 2, func(cpu *CPU, operands []byte) { cpu.sla_n(&cpu.Regs.E) }},
	&inst{0x24, "SLA H", 0, 2, func(cpu *CPU, operands []byte) { cpu.sla_n(&cpu.Regs.H) }},
	&inst{0x25, "SLA L", 0, 2, func(cpu *CPU, operands []byte) { cpu.sla_n(&cpu.Regs.L) }},
	&inst{0x26, "SLA (HL)", 0, 4, func(cpu *CPU, operands []byte) { cpu.sla_hl() }},
	&inst{0x27, "SLA A", 0, 2, func(cpu *CPU, operands []byte) { cpu.sla_n(&cpu.Regs.A) }},
	&inst{0x28, "SRA B", 0, 2, func(cpu *CPU, operands []byte) { cpu.sra_n(&cpu.Regs.B) }},
	&inst{0x29, "SRA C", 0, 2, func(cpu *CPU, operands []byte) { cpu.sra_n(&cpu.Regs.C) }},
	&inst{0x2A, "SRA D", 0, 2, func(cpu *CPU, operands []byte) { cpu.sra_n(&cpu.Regs.D) }},
	&inst{0x2B, "SRA E", 0, 2, func(cpu *CPU, operands []byte) { cpu.sra_n(&cpu.Regs.E) }},
	&inst{0x2C, "SRA H", 0, 2, func(cpu *CPU, operands []byte) { cpu.sra_n(&cpu.Regs.H) }},
	&inst{0x2D, "SRA L", 0, 2, func(cpu *CPU, operands []byte) { cpu.sra_n(&cpu.Regs.L) }},
	&inst{0x2E, "SRA (HL)", 0, 4, func(cpu *CPU, operands []byte) { cpu.sra_hl() }},
	&inst{0x2F, "SRA A", 0, 2, func(cpu *CPU, operands []byte) { cpu.sra_n(&cpu.Regs.A) }},
	&inst{0x30, "SWAP B", 0, 2, func(cpu *CPU, operands []byte) { cpu.swap_n(&cpu.Regs.B) }},
	&inst{0x31, "SWAP C", 0, 2, func(cpu *CPU, operands []byte) { cpu.swap_n(&cpu.Regs.C) }},
	&inst{0x32, "SWAP D", 0, 2, func(cpu *CPU, operands []byte) { cpu.swap_n(&cpu.Regs.D) }},
	&inst{0x33, "SWAP E", 0, 2, func(cpu *CPU, operands []byte) { cpu.swap_n(&cpu.Regs.E) }},
	&inst{0x34, "SWAP H", 0, 2, func(cpu *CPU, operands []byte) { cpu.swap_n(&cpu.Regs.H) }},
	&inst{0x35, "SWAP L", 0, 2, func(cpu *CPU, operands []byte) { cpu.swap_n(&cpu.Regs.L) }},
	&inst{0x36, "SWAP (HL)", 0, 4, func(cpu *CPU, operands []byte) { cpu.swap_hl() }},
	&inst{0x37, "SWAP A", 0, 2, func(cpu *CPU, operands []byte) { cpu.swap_n(&cpu.Regs.A) }},
	&inst{0x38, "SRL B", 0, 2, func(cpu *CPU, operands []byte) { cpu.srl_n(&cpu.Regs.B) }},
	&inst{0x39, "SRL C", 0, 2, func(cpu *CPU, operands []byte) { cpu.srl_n(&cpu.Regs.C) }},
	&inst{0x3A, "SRL D", 0, 2, func(cpu *CPU, operands []byte) { cpu.srl_n(&cpu.Regs.D) }},
	&inst{0x3B, "SRL E", 0, 2, func(cpu *CPU, operands []byte) { cpu.srl_n(&cpu.Regs.E) }},
	&inst{0x3C, "SRL H", 0, 2, func(cpu *CPU, operands []byte) { cpu.srl_n(&cpu.Regs.H) }},
	&inst{0x3D, "SRL L", 0, 2, func(cpu *CPU, operands []byte) { cpu.srl_n(&cpu.Regs.L) }},
	&inst{0x3E, "SRL (HL)", 0, 4, func(cpu *CPU, operands []byte) { cpu.srl_hl() }},
	&inst{0x3F, "SRL A", 0, 2, func(cpu *CPU, operands []byte) { cpu.srl_n(&cpu.Regs.A) }},
	&inst{0x40, "BIT 0,B", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit0, &cpu.Regs.B) }},
	&inst{0x41, "BIT 0,C", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit0, &cpu.Regs.C) }},
	&inst{0x42, "BIT 0,D", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit0, &cpu.Regs.D) }},
	&inst{0x43, "BIT 0,E", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit0, &cpu.Regs.E) }},
	&inst{0x44, "BIT 0,H", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit0, &cpu.Regs.H) }},
	&inst{0x45, "BIT 0,L", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit0, &cpu.Regs.L) }},
	&inst{0x46, "BIT 0,(HL)", 0, 3, func(cpu *CPU, operands []byte) { cpu.bit_b_hl(types.Bit0) }},
	&inst{0x47, "BIT 0,A", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit0, &cpu.Regs.A) }},
	&inst{0x48, "BIT 1,B", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit1, &cpu.Regs.B) }},
	&inst{0x49, "BIT 1,C", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit1, &cpu.Regs.C) }},
	&inst{0x4A, "BIT 1,D", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit1, &cpu.Regs.D) }},
	&inst{0x4B, "BIT 1,E", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit1, &cpu.Regs.E) }},
	&inst{0x4C, "BIT 1,H", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit1, &cpu.Regs.H) }},
	&inst{0x4D, "BIT 1,L", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit1, &cpu.Regs.L) }},
	&inst{0x4E, "BIT 1,(HL)", 0, 3, func(cpu *CPU, operands []byte) { cpu.bit_b_hl(types.Bit1) }},
	&inst{0x4F, "BIT 1,A", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit1, &cpu.Regs.A) }},
	&inst{0x50, "BIT 2,B", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit2, &cpu.Regs.B) }},
	&inst{0x51, "BIT 2,C", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit2, &cpu.Regs.C) }},
	&inst{0x52, "BIT 2,D", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit2, &cpu.Regs.D) }},
	&inst{0x53, "BIT 2,E", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit2, &cpu.Regs.E) }},
	&inst{0x54, "BIT 2,H", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit2, &cpu.Regs.H) }},
	&inst{0x55, "BIT 2,L", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit2, &cpu.Regs.L) }},
	&inst{0x56, "BIT 2,(HL)", 0, 3, func(cpu *CPU, operands []byte) { cpu.bit_b_hl(types.Bit2) }},
	&inst{0x57, "BIT 2,A", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit2, &cpu.Regs.A) }},
	&inst{0x58, "BIT 3,B", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit3, &cpu.Regs.B) }},
	&inst{0x59, "BIT 3,C", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit3, &cpu.Regs.C) }},
	&inst{0x5A, "BIT 3,D", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit3, &cpu.Regs.D) }},
	&inst{0x5B, "BIT 3,E", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit3, &cpu.Regs.E) }},
	&inst{0x5C, "BIT 3,H", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit3, &cpu.Regs.H) }},
	&inst{0x5D, "BIT 3,L", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit3, &cpu.Regs.L) }},
	&inst{0x5E, "BIT 3,(HL)", 0, 3, func(cpu *CPU, operands []byte) { cpu.bit_b_hl(types.Bit3) }},
	&inst{0x5F, "BIT 3,A", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit3, &cpu.Regs.A) }},
	&inst{0x60, "BIT 4,B", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit4, &cpu.Regs.B) }},
	&inst{0x61, "BIT 4,C", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit4, &cpu.Regs.C) }},
	&inst{0x62, "BIT 4,D", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit4, &cpu.Regs.D) }},
	&inst{0x63, "BIT 4,E", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit4, &cpu.Regs.E) }},
	&inst{0x64, "BIT 4,H", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit4, &cpu.Regs.H) }},
	&inst{0x65, "BIT 4,L", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit4, &cpu.Regs.L) }},
	&inst{0x66, "BIT 4,(HL)", 0, 3, func(cpu *CPU, operands []byte) { cpu.bit_b_hl(types.Bit4) }},
	&inst{0x67, "BIT 4,A", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit4, &cpu.Regs.A) }},
	&inst{0x68, "BIT 5,B", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit5, &cpu.Regs.B) }},
	&inst{0x69, "BIT 5,C", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit5, &cpu.Regs.C) }},
	&inst{0x6A, "BIT 5,D", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit5, &cpu.Regs.D) }},
	&inst{0x6B, "BIT 5,E", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit5, &cpu.Regs.E) }},
	&inst{0x6C, "BIT 5,H", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit5, &cpu.Regs.H) }},
	&inst{0x6D, "BIT 5,L", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit5, &cpu.Regs.L) }},
	&inst{0x6E, "BIT 5,(HL)", 0, 3, func(cpu *CPU, operands []byte) { cpu.bit_b_hl(types.Bit5) }},
	&inst{0x6F, "BIT 5,A", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit5, &cpu.Regs.A) }},
	&inst{0x70, "BIT 6,B", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit6, &cpu.Regs.B) }},
	&inst{0x71, "BIT 6,C", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit6, &cpu.Regs.C) }},
	&inst{0x72, "BIT 6,D", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit6, &cpu.Regs.D) }},
	&inst{0x73, "BIT 6,E", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit6, &cpu.Regs.E) }},
	&inst{0x74, "BIT 6,H", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit6, &cpu.Regs.H) }},
	&inst{0x75, "BIT 6,L", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit6, &cpu.Regs.L) }},
	&inst{0x76, "BIT 6,(HL)", 0, 3, func(cpu *CPU, operands []byte) { cpu.bit_b_hl(types.Bit6) }},
	&inst{0x77, "BIT 6,A", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit6, &cpu.Regs.A) }},
	&inst{0x78, "BIT 7,B", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit7, &cpu.Regs.B) }},
	&inst{0x79, "BIT 7,C", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit7, &cpu.Regs.C) }},
	&inst{0x7A, "BIT 7,D", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit7, &cpu.Regs.D) }},
	&inst{0x7B, "BIT 7,E", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit7, &cpu.Regs.E) }},
	&inst{0x7C, "BIT 7,H", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit7, &cpu.Regs.H) }},
	&inst{0x7D, "BIT 7,L", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit7, &cpu.Regs.L) }},
	&inst{0x7E, "BIT 7,(HL)", 0, 3, func(cpu *CPU, operands []byte) { cpu.bit_b_hl(types.Bit7) }},
	&inst{0x7F, "BIT 7,A", 0, 2, func(cpu *CPU, operands []byte) { cpu.bit_b_r(types.Bit7, &cpu.Regs.A) }},
	&inst{0x80, "RES 0,B", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit0, &cpu.Regs.B) }},
	&inst{0x81, "RES 0,C", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit0, &cpu.Regs.C) }},
	&inst{0x82, "RES 0,D", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit0, &cpu.Regs.D) }},
	&inst{0x83, "RES 0,E", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit0, &cpu.Regs.E) }},
	&inst{0x84, "RES 0,H", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit0, &cpu.Regs.H) }},
	&inst{0x85, "RES 0,L", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit0, &cpu.Regs.L) }},
	&inst{0x86, "RES 0,(HL)", 0, 4, func(cpu *CPU, operands []byte) { cpu.res_b_hl(types.Bit0) }},
	&inst{0x87, "RES 0,A", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit0, &cpu.Regs.A) }},
	&inst{0x88, "RES 1,B", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit1, &cpu.Regs.B) }},
	&inst{0x89, "RES 1,C", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit1, &cpu.Regs.C) }},
	&inst{0x8A, "RES 1,D", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit1, &cpu.Regs.D) }},
	&inst{0x8B, "RES 1,E", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit1, &cpu.Regs.E) }},
	&inst{0x8C, "RES 1,H", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit1, &cpu.Regs.H) }},
	&inst{0x8D, "RES 1,L", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit1, &cpu.Regs.L) }},
	&inst{0x8E, "RES 1,(HL)", 0, 4, func(cpu *CPU, operands []byte) { cpu.res_b_hl(types.Bit1) }},
	&inst{0x8F, "RES 1,A", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit1, &cpu.Regs.A) }},
	&inst{0x90, "RES 2,B", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit2, &cpu.Regs.B) }},
	&inst{0x91, "RES 2,C", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit2, &cpu.Regs.C) }},
	&inst{0x92, "RES 2,D", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit2, &cpu.Regs.D) }},
	&inst{0x93, "RES 2,E", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit2, &cpu.Regs.E) }},
	&inst{0x94, "RES 2,H", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit2, &cpu.Regs.H) }},
	&inst{0x95, "RES 2,L", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit2, &cpu.Regs.L) }},
	&inst{0x96, "RES 2,(HL)", 0, 4, func(cpu *CPU, operands []byte) { cpu.res_b_hl(types.Bit2) }},
	&inst{0x97, "RES 2,A", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit2, &cpu.Regs.A) }},
	&inst{0x98, "RES 3,B", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit3, &cpu.Regs.B) }},
	&inst{0x99, "RES 3,C", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit3, &cpu.Regs.C) }},
	&inst{0x9A, "RES 3,D", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit3, &cpu.Regs.D) }},
	&inst{0x9B, "RES 3,E", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit3, &cpu.Regs.E) }},
	&inst{0x9C, "RES 3,H", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit3, &cpu.Regs.H) }},
	&inst{0x9D, "RES 3,L", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit3, &cpu.Regs.L) }},
	&inst{0x9E, "RES 3,(HL)", 0, 4, func(cpu *CPU, operands []byte) { cpu.res_b_hl(types.Bit3) }},
	&inst{0x9F, "RES 3,A", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit3, &cpu.Regs.A) }},
	&inst{0xA0, "RES 4,B", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit4, &cpu.Regs.B) }},
	&inst{0xA1, "RES 4,C", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit4, &cpu.Regs.C) }},
	&inst{0xA2, "RES 4,D", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit4, &cpu.Regs.D) }},
	&inst{0xA3, "RES 4,E", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit4, &cpu.Regs.E) }},
	&inst{0xA4, "RES 4,H", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit4, &cpu.Regs.H) }},
	&inst{0xA5, "RES 4,L", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit4, &cpu.Regs.L) }},
	&inst{0xA6, "RES 4,(HL)", 0, 4, func(cpu *CPU, operands []byte) { cpu.res_b_hl(types.Bit4) }},
	&inst{0xA7, "RES 4,A", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit4, &cpu.Regs.A) }},
	&inst{0xA8, "RES 5,B", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit5, &cpu.Regs.B) }},
	&inst{0xA9, "RES 5,C", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit5, &cpu.Regs.C) }},
	&inst{0xAA, "RES 5,D", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit5, &cpu.Regs.D) }},
	&inst{0xAB, "RES 5,E", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit5, &cpu.Regs.E) }},
	&inst{0xAC, "RES 5,H", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit5, &cpu.Regs.H) }},
	&inst{0xAD, "RES 5,L", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit5, &cpu.Regs.L) }},
	&inst{0xAE, "RES 5,(HL)", 0, 4, func(cpu *CPU, operands []byte) { cpu.res_b_hl(types.Bit5) }},
	&inst{0xAF, "RES 5,A", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit5, &cpu.Regs.A) }},
	&inst{0xB0, "RES 6,B", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit6, &cpu.Regs.B) }},
	&inst{0xB1, "RES 6,C", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit6, &cpu.Regs.C) }},
	&inst{0xB2, "RES 6,D", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit6, &cpu.Regs.D) }},
	&inst{0xB3, "RES 6,E", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit6, &cpu.Regs.E) }},
	&inst{0xB4, "RES 6,H", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit6, &cpu.Regs.H) }},
	&inst{0xB5, "RES 6,L", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit6, &cpu.Regs.L) }},
	&inst{0xB6, "RES 6,(HL)", 0, 4, func(cpu *CPU, operands []byte) { cpu.res_b_hl(types.Bit6) }},
	&inst{0xB7, "RES 6,A", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit6, &cpu.Regs.A) }},
	&inst{0xB8, "RES 7,B", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit7, &cpu.Regs.B) }},
	&inst{0xB9, "RES 7,C", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit7, &cpu.Regs.C) }},
	&inst{0xBA, "RES 7,D", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit7, &cpu.Regs.D) }},
	&inst{0xBB, "RES 7,E", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit7, &cpu.Regs.E) }},
	&inst{0xBC, "RES 7,H", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit7, &cpu.Regs.H) }},
	&inst{0xBD, "RES 7,L", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit7, &cpu.Regs.L) }},
	&inst{0xBE, "RES 7,(HL)", 0, 4, func(cpu *CPU, operands []byte) { cpu.res_b_hl(types.Bit7) }},
	&inst{0xBF, "RES 7,A", 0, 2, func(cpu *CPU, operands []byte) { cpu.res_b_r(types.Bit7, &cpu.Regs.A) }},
	&inst{0xC0, "SET 0,B", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit0, &cpu.Regs.B) }},
	&inst{0xC1, "SET 0,C", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit0, &cpu.Regs.C) }},
	&inst{0xC2, "SET 0,D", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit0, &cpu.Regs.D) }},
	&inst{0xC3, "SET 0,E", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit0, &cpu.Regs.E) }},
	&inst{0xC4, "SET 0,H", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit0, &cpu.Regs.H) }},
	&inst{0xC5, "SET 0,L", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit0, &cpu.Regs.L) }},
	&inst{0xC6, "SET 0,(HL)", 0, 4, func(cpu *CPU, operands []byte) { cpu.set_b_hl(types.Bit0) }},
	&inst{0xC7, "SET 0,A", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit0, &cpu.Regs.A) }},
	&inst{0xC8, "SET 1,B", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit1, &cpu.Regs.B) }},
	&inst{0xC9, "SET 1,C", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit1, &cpu.Regs.C) }},
	&inst{0xCA, "SET 1,D", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit1, &cpu.Regs.D) }},
	&inst{0xCB, "SET 1,E", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit1, &cpu.Regs.E) }},
	&inst{0xCC, "SET 1,H", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit1, &cpu.Regs.H) }},
	&inst{0xCD, "SET 1,L", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit1, &cpu.Regs.L) }},
	&inst{0xCE, "SET 1,(HL)", 0, 4, func(cpu *CPU, operands []byte) { cpu.set_b_hl(types.Bit1) }},
	&inst{0xCF, "SET 1,A", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit1, &cpu.Regs.A) }},
	&inst{0xD0, "SET 2,B", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit2, &cpu.Regs.B) }},
	&inst{0xD1, "SET 2,C", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit2, &cpu.Regs.C) }},
	&inst{0xD2, "SET 2,D", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit2, &cpu.Regs.D) }},
	&inst{0xD3, "SET 2,E", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit2, &cpu.Regs.E) }},
	&inst{0xD4, "SET 2,H", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit2, &cpu.Regs.H) }},
	&inst{0xD5, "SET 2,L", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit2, &cpu.Regs.L) }},
	&inst{0xD6, "SET 2,(HL)", 0, 4, func(cpu *CPU, operands []byte) { cpu.set_b_hl(types.Bit2) }},
	&inst{0xD7, "SET 2,A", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit2, &cpu.Regs.A) }},
	&inst{0xD8, "SET 3,B", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit3, &cpu.Regs.B) }},
	&inst{0xD9, "SET 3,C", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit3, &cpu.Regs.C) }},
	&inst{0xDA, "SET 3,D", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit3, &cpu.Regs.D) }},
	&inst{0xDB, "SET 3,E", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit3, &cpu.Regs.E) }},
	&inst{0xDC, "SET 3,H", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit3, &cpu.Regs.H) }},
	&inst{0xDD, "SET 3,L", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit3, &cpu.Regs.L) }},
	&inst{0xDE, "SET 3,(HL)", 0, 4, func(cpu *CPU, operands []byte) { cpu.set_b_hl(types.Bit3) }},
	&inst{0xDF, "SET 3,A", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit3, &cpu.Regs.A) }},
	&inst{0xE0, "SET 4,B", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit4, &cpu.Regs.B) }},
	&inst{0xE1, "SET 4,C", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit4, &cpu.Regs.C) }},
	&inst{0xE2, "SET 4,D", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit4, &cpu.Regs.D) }},
	&inst{0xE3, "SET 4,E", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit4, &cpu.Regs.E) }},
	&inst{0xE4, "SET 4,H", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit4, &cpu.Regs.H) }},
	&inst{0xE5, "SET 4,L", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit4, &cpu.Regs.L) }},
	&inst{0xE6, "SET 4,(HL)", 0, 4, func(cpu *CPU, operands []byte) { cpu.set_b_hl(types.Bit4) }},
	&inst{0xE7, "SET 4,A", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit4, &cpu.Regs.A) }},
	&inst{0xE8, "SET 5,B", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit5, &cpu.Regs.B) }},
	&inst{0xE9, "SET 5,C", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit5, &cpu.Regs.C) }},
	&inst{0xEA, "SET 5,D", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit5, &cpu.Regs.D) }},
	&inst{0xEB, "SET 5,E", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit5, &cpu.Regs.E) }},
	&inst{0xEC, "SET 5,H", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit5, &cpu.Regs.H) }},
	&inst{0xED, "SET 5,L", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit5, &cpu.Regs.L) }},
	&inst{0xEE, "SET 5,(HL)", 0, 4, func(cpu *CPU, operands []byte) { cpu.set_b_hl(types.Bit5) }},
	&inst{0xEF, "SET 5,A", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit5, &cpu.Regs.A) }},
	&inst{0xF0, "SET 6,B", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit6, &cpu.Regs.B) }},
	&inst{0xF1, "SET 6,C", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit6, &cpu.Regs.C) }},
	&inst{0xF2, "SET 6,D", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit6, &cpu.Regs.D) }},
	&inst{0xF3, "SET 6,E", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit6, &cpu.Regs.E) }},
	&inst{0xF4, "SET 6,H", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit6, &cpu.Regs.H) }},
	&inst{0xF5, "SET 6,L", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit6, &cpu.Regs.L) }},
	&inst{0xF6, "SET 6,(HL)", 0, 4, func(cpu *CPU, operands []byte) { cpu.set_b_hl(types.Bit6) }},
	&inst{0xF7, "SET 6,A", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit6, &cpu.Regs.A) }},
	&inst{0xF8, "SET 7,B", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit7, &cpu.Regs.B) }},
	&inst{0xF9, "SET 7,C", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit7, &cpu.Regs.C) }},
	&inst{0xFA, "SET 7,D", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit7, &cpu.Regs.D) }},
	&inst{0xFB, "SET 7,E", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit7, &cpu.Regs.E) }},
	&inst{0xFC, "SET 7,H", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit7, &cpu.Regs.H) }},
	&inst{0xFD, "SET 7,L", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit7, &cpu.Regs.L) }},
	&inst{0xFE, "SET 7,(HL)", 0, 4, func(cpu *CPU, operands []byte) { cpu.set_b_hl(types.Bit7) }},
	&inst{0xFF, "SET 7,A", 0, 2, func(cpu *CPU, operands []byte) { cpu.set_b_r(types.Bit7, &cpu.Regs.A) }},
}

var instructions = []*inst{
	&inst{0x0, "NOP", 0, 1, func(cpu *CPU, operands []byte) { cpu.nop() }},
	&inst{0x1, "LD BC,nn", 2, 3, func(cpu *CPU, operands []byte) { cpu.ldn_nn(&cpu.Regs.B, &cpu.Regs.C, operands) }},
	&inst{0x2, "LD (BC),A", 0, 2, func(cpu *CPU, operands []byte) { cpu.ldrr_r(cpu.Regs.B, cpu.Regs.C, cpu.Regs.A) }},
	&inst{0x3, "INC BC", 0, 2, func(cpu *CPU, operands []byte) { cpu.inc_nn(&cpu.Regs.B, &cpu.Regs.C) }},
	&inst{0x4, "INC B", 0, 1, func(cpu *CPU, operands []byte) { cpu.inc_n(&cpu.Regs.B) }},
	&inst{0x5, "DEC B", 0, 1, func(cpu *CPU, operands []byte) { cpu.dec_n(&cpu.Regs.B) }},
	&inst{0x6, "LD B,n", 1, 2, func(cpu *CPU, operands []byte) { cpu.ldnn_n(&cpu.Regs.B, operands) }},
	&inst{0x7, "RLCA", 0, 1, func(cpu *CPU, operands []byte) { cpu.rlca() }},
	&inst{0x8, "LD (nn),SP", 2, 5, func(cpu *CPU, operands []byte) { cpu.ldnn_sp(operands) }},
	&inst{0x9, "ADD HL,BC", 0, 2, func(cpu *CPU, operands []byte) { cpu.addhl_rr(&cpu.Regs.B, &cpu.Regs.C) }},
	&inst{0xA, "LD A,(BC)", 0, 2, func(cpu *CPU, operands []byte) { cpu.ldr_rr(cpu.Regs.B, cpu.Regs.C, &cpu.Regs.A) }},
	&inst{0xB, "DEC BC", 0, 2, func(cpu *CPU, operands []byte) { cpu.dec_nn(&cpu.Regs.B, &cpu.Regs.C) }},
	&inst{0xC, "INC C", 0, 1, func(cpu *CPU, operands []byte) { cpu.inc_n(&cpu.Regs.C) }},
	&inst{0xD, "DEC C", 0, 1, func(cpu *CPU, operands []byte) { cpu.dec_n(&cpu.Regs.C) }},
	&inst{0xE, "LD C,n", 1, 2, func(cpu *CPU, operands []byte) { cpu.ldnn_n(&cpu.Regs.C, operands) }},
	&inst{0xF, "RRCA", 0, 1, func(cpu *CPU, operands []byte) { cpu.rrca() }},
	&inst{0x10, "STOP", 1, 0, func(cpu *CPU, operands []byte) { cpu.stop() }},
	&inst{0x11, "LD DE,(nn)", 2, 3, func(cpu *CPU, operands []byte) { cpu.ldn_nn(&cpu.Regs.D, &cpu.Regs.E, operands) }},
	&inst{0x12, "LD (DE),A", 0, 2, func(cpu *CPU, operands []byte) { cpu.ldrr_r(cpu.Regs.D, cpu.Regs.E, cpu.Regs.A) }},
	&inst{0x13, "INC DE", 0, 2, func(cpu *CPU, operands []byte) { cpu.inc_nn(&cpu.Regs.D, &cpu.Regs.E) }},
	&inst{0x14, "INC D", 0, 1, func(cpu *CPU, operands []byte) { cpu.inc_n(&cpu.Regs.D) }},
	&inst{0x15, "DEC D", 0, 1, func(cpu *CPU, operands []byte) { cpu.dec_n(&cpu.Regs.D) }},
	&inst{0x16, "LD D,n", 1, 2, func(cpu *CPU, operands []byte) { cpu.ldnn_n(&cpu.Regs.D, operands) }},
	&inst{0x17, "RLA", 0, 1, func(cpu *CPU, operands []byte) { cpu.rla() }},
	&inst{0x18, "JR n", 1, 3, func(cpu *CPU, operands []byte) { cpu.jr_n(operands) }},
	&inst{0x19, "ADD HL,DE", 0, 2, func(cpu *CPU, operands []byte) { cpu.addhl_rr(&cpu.Regs.D, &cpu.Regs.E) }},
	&inst{0x1A, "LD A,(DE)", 0, 2, func(cpu *CPU, operands []byte) { cpu.ldr_rr(cpu.Regs.D, cpu.Regs.E, &cpu.Regs.A) }},
	&inst{0x1B, "DEC DE", 0, 2, func(cpu *CPU, operands []byte) { cpu.dec_nn(&cpu.Regs.D, &cpu.Regs.E) }},
	&inst{0x1C, "INC E", 0, 1, func(cpu *CPU, operands []byte) { cpu.inc_n(&cpu.Regs.E) }},
	&inst{0x1D, "DEC E", 0, 1, func(cpu *CPU, operands []byte) { cpu.dec_n(&cpu.Regs.E) }},
	&inst{0x1E, "LD E,n", 1, 2, func(cpu *CPU, operands []byte) { cpu.ldnn_n(&cpu.Regs.E, operands) }},
	&inst{0x1F, "RRA", 0, 1, func(cpu *CPU, operands []byte) { cpu.rra() }},
	&inst{0x20, "JR NZ,*", 1, 2, func(cpu *CPU, operands []byte) { cpu.jrcc_n(Z, false, operands) }},
	&inst{0x21, "LD HL,nn", 2, 3, func(cpu *CPU, operands []byte) { cpu.ldn_nn(&cpu.Regs.H, &cpu.Regs.L, operands) }},
	&inst{0x22, "LD (HL+),A", 0, 2, func(cpu *CPU, operands []byte) { cpu.ldihl_a() }},
	&inst{0x23, "INC HL", 0, 2, func(cpu *CPU, operands []byte) { cpu.inc_nn(&cpu.Regs.H, &cpu.Regs.L) }},
	&inst{0x24, "INC H", 0, 1, func(cpu *CPU, operands []byte) { cpu.inc_n(&cpu.Regs.H) }},
	&inst{0x25, "DEC H", 0, 1, func(cpu *CPU, operands []byte) { cpu.dec_n(&cpu.Regs.H) }},
	&inst{0x26, "LD H,n", 1, 2, func(cpu *CPU, operands []byte) { cpu.ldnn_n(&cpu.Regs.H, operands) }},
	&inst{0x27, "DAA", 0, 1, func(cpu *CPU, operands []byte) { cpu.daa() }},
	&inst{0x28, "JR Z,*", 1, 2, func(cpu *CPU, operands []byte) { cpu.jrcc_n(Z, true, operands) }},