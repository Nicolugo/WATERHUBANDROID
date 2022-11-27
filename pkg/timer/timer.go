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
	// The timer generates an interrupt when it overfl