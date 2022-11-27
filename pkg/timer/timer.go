package timer

import (
	"github.com/bokuweb/gopher-boy/pkg/types"
)

const (
	// TimerRegisterOffset is register offset address
	TimerRegisterOffset types.Word = 0xFF00
	// DIV - Divider Register (R/W)
	// This register is incremented 16384 (~16779 on SGB) times a second.
	// Writing any value sets it to $00.
	DIV = 0x04
	// TIMA - Timer counter (R/W)
	// This timer is incremented by a clock frequency specified by the TAC register ($FF07).
	// The timer generates an interrupt when it overflows.
	TIMA = 0x05
	// TMA - Timer Modulo (R/W)
	// When the TIMA overflows, this data will be loaded.
	TMA = 0x06
	// TAC - Timer Control (R/W)
	// Bit 2 - Timer Stop
	//         0: Stop Timer
	//         1: Start Timer
	// Bits 1+0 - Input Clock Select
	//         00: 4