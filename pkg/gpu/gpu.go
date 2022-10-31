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
	stat            byte
	ly              uint
	lyc             byte
	scrollX         byte
	scrollY         byte
	windowX         byte
	windowY         byte
	bgPalette       byte
	objPalette0     byte
	objPalette1     byte
	disableDisplay  bool
	oamDMAStarted   bool
	oamDMAStartAddr types.Word
}

// GPUMode
type GPUMode = byte

const (
	// HBlankMode is period CPU can access the display RAM ($8000-$9FFF).
	HBlankMode GPUMode = iota
	// period and the CPU can access the display RAM ($8000-$9FFF).
	VBlankMode
	SearchingOAMMode
	TransferingData
)

// G