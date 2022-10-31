package gpu

import (
	"image/color"

	"github.com/bokuweb/gopher-boy/pkg/constants"
	"github.com/bokuweb/gopher-boy/pkg/interfaces/bus"
	"github.com/bokuweb/gopher-boy/pkg/interfaces/interrupt"
	irq "github.com/bokuweb/gopher-boy/pkg/interrupt"
	"github.com/bokuweb/gopher-boy/pkg/types"
)

// CyclePerLine is gpu clock count per line
const CyclePerLine uint = 456

// LCDVBlankHeight means vblank height
const LCDVBlankHeight uint = 10

const spriteNum = 40

// GPU is
type GPU struct {
	bus             bus.Accessor
	irq             interrupt.Interrupt
	imageData       []byte
	mode            GPUMode
	clock           uint
	lcdc            byte