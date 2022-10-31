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
	