package gb

import (
	"time"

	"github.com/bokuweb/gopher-boy/pkg/cpu"
	"github.com/bokuweb/gopher-boy/pkg/gpu"
	"github.com/bokuweb/gopher-boy/pkg/interfaces/window"
	"github.com/bokuweb/gopher-boy/pkg/interrupt"
	"github.com/bokuweb/gopher-boy/pkg/timer"
)

// CyclesPerFrame is cpu clock num for 1frame.
const CyclesPerFrame = 70224

// GB is gameboy emulator struct
type GB struct {
	currentCycle uint
	cpu          *cpu.CPU
	gpu          *gpu.GPU
	timer        *timer.Timer
	irq          *interrupt.Interrupt
	win          window.Window
}

// NewGB is gb initializer
func NewGB(cpu *cpu.CPU, gpu *gpu.GPU, timer *timer.Timer, irq *interrupt.Interrupt, win window.Window) *GB {
	return &GB{
		currentCycle: 0,
		cpu:          cpu,
		gpu:          gpu,
		timer:        timer,
		irq:          irq,
		win:          win,
	}
}

// Start is
func (g *GB) Start() {
	t := time.NewTicker(16 * time.Millisecond)
	for {
		select {
		case <-t.C:
			buf := g.Next()
			g.win.Render(buf)
		}
	}
	t.Stop()
}
func (g *GB) Next() []byte {
	for {
		var cycles uint
		if 