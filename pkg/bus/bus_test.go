package bus

import (
	"testing"

	"github.com/bokuweb/gopher-boy/pkg/cartridge"
	"github.com/bokuweb/gopher-boy/pkg/gpu"
	"github.com/bokuweb/gopher-boy/pkg/interrupt"
	"github.com/bokuweb/gopher-boy/pkg/logger"
	"github.com/bokuweb/gopher-boy/pkg/pad"
	"github.com/bokuweb/gopher-boy/pkg/ram"
	"github.com/bokuweb/gopher-boy/pkg/timer"
	"github.com/bokuweb/gopher-boy/pkg/types"
	"github.com/stretchr/testify/assert"
)

func setup() (*Bus, *ram.RAM, *ram.RAM) {
	buf := make([]byte, 0x8000)
	cart, _ := cartridge.NewCartridge(buf)
	vRAM := ram.NewRAM(0x2000)
	wRAM := ram.NewRAM(0x2000)
	hRAM := ram.NewRAM(0x80)
	oamRAM := ram.NewRAM(0xA0)
	gpu := gpu.NewGPU()
	pad := pad.NewPad()
	l := logger.NewLogger(logger.LogLevel("Debug"))
	t := timer.NewTimer()
	irq := interrupt.NewInterrupt()
	return NewBus(l, cart, gpu, vRAM, wRAM, hRAM, oamRAM, t, irq, pad), wRAM, hRAM
}

func TestWRAMReadWrite(t *testing.T) {
	ass