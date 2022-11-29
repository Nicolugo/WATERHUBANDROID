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
	//         00: 4.096 KHz (~4.194 KHz SGB)
	//         01: 262.144 Khz (~268.4 KHz SGB)
	//         10: 65.536 KHz (~67.11 KHz SGB)
	//         11: 16.384 KHz (~16.78 KHz SGB)
	TAC = 0x07
)

// Timer has 4 registers.
type Timer struct {
	internalCounter uint16
	TIMA            byte
	TAC             byte
	TMA             byte
}

// NewTimer constructs timer peripheral.
func NewTimer() *Timer {
	return &Timer{
		// 4.194304MHz / 256 = 16.384KHz
		internalCounter: 0,
		TIMA:            0x00,
		TAC:             0x00,
		TMA:             0x00,
	}
}

// Update timer counter registers
// If timer is overflowed return true
func (timer *Timer) Update(cycles uint) bool {
	r := false
	for cycles > 0 {
		cycles--
		old := timer.internalCounter
		timer.internalCounter += 4

		if !timer.isStarted()