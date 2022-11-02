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

// GPU register addresses
const (
	LCDC types.Word = 0x00
	STAT            = 0x01
	// Scroll Y (R/W)
	// 8 Bit value $00-$FF to scroll BG Y screen
	// position.
	SCROLLY = 0x02
	// Scroll X (R/W)
	// 8 Bit value $00-$FF to scroll BG X screen
	// position.
	SCROLLX = 0x03
	// LY Y-Coordinate (R)
	// The LY indicates the vertical line to which
	// the present data is transferred to the LCD
	// Driver. The LY can take on any value
	// between 0 through 153. The values between
	// 144 and 153 indicate the V-Blank period.
	// Writing will reset the counter.
	LY  = 0x04
	LYC = 0x05
	// BGP - BG & Window Palette Data (R/W)
	// Bit 7-6 - Data for Dot Data 11
	// (Normally darkest color)
	// Bit 5-4 - Data for Dot Data 10
	// Bit 3-2 - Data for Dot Data 01
	// Bit 1-0 - Data for Dot Data 00
	// (Normally lightest color)
	// This selects the shade of grays to use
	// for the background (BG) & window pixels.
	// Since each pixel uses 2 bits, the
	// corresponding shade will be selected from here.
	DMA  = 0x06
	BGP  = 0x07
	OBP0 = 0x08
	OBP1 = 0x09
	WX   = 0x0B
	WY   = 0x0A
)

const (
	TILEMAP0  types.Word = 0x9800
	TI