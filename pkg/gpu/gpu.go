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
	TILEMAP1             = 0x9C00
	TILEDATA0            = 0x8800
	TILEDATA1            = 0x8000
	OAMSTART             = 0xFE00
)

// NewGPU is GPU constructor
func NewGPU() *GPU {
	return &GPU{
		imageData:       make([]byte, constants.ScreenWidth*constants.ScreenHeight*4),
		mode:            HBlankMode,
		clock:           0,
		lcdc:            0x91,
		ly:              0,
		scrollX:         0,
		scrollY:         0,
		oamDMAStarted:   false,
		oamDMAStartAddr: 0,
	}
}

// Init initialize GPU
func (g *GPU) Init(bus bus.Accessor, irq interrupt.Interrupt) {
	g.bus = bus
	g.irq = irq
}

// Step is run GPU
func (g *GPU) Step(cycles uint) {
	if g.bus == nil {
		panic("Please initialize gpu with Init, before running.")
	}
	g.updateMode()

	g.clock += cycles

	if !g.lcdEnabled() {
		return
	}
	if g.clock >= CyclePerLine {
		if g.ly == constants.ScreenHeight {
			g.buildSprites()
			g.irq.SetIRQ(irq.VerticalBlankFlag)
			if g.vBlankInterruptEnabled() {
				g.irq.SetIRQ(irq.LCDSFlag)
			}
		} else if g.ly >= constants.ScreenHeight+LCDVBlankHeight {
			g.ly = 0
			g.buildBGTile()
		} else if g.ly < constants.ScreenHeight {
			g.buildBGTile()
			if g.windowEnabled() {
				g.buildWindowTile()
			}
		}

		if g.ly == uint(g.lyc) {
			g.stat |= 0x04
			if g.coincidenceInterruptEnabled() {
				g.irq.SetIRQ(irq.LCDSFlag)
			}
		} else {
			g.stat &= 0xFB
		}
		g.ly++
		g.clock -= CyclePerLine
	}
}

func (g *GPU) lcdEnabled() bool {
	return (g.lcdc & 0x80) == 0x80
}

func (g *GPU) longSprite() bool {
	return (g.lcdc & 0x04) == 0x04
}

func (g *GPU) coincidenceInterruptEnabled() bool {
	return (g.stat & 0x40) == 0x40
}

func (g *GPU) vBlankInterruptEnabled() bool {
	return (g.stat & 0x10) == 0x10
}

func (g *GPU) hblankInterruptEnabled() bool {
	return (g.stat & 0x08) == 0x08
}

func (g *GPU) Read(addr types.Word) byte {
	switch addr {
	case LCDC:
		return g.lcdc
	case STAT:
		return g.stat&0xF8 | (byte(g.mode)) | 0x80
	case SCROLLX:
		return g.scrollX
	case SCROLLY:
		return g.scrollY
	case LY:
		return byte(g.ly)
	case BGP:
		return g.bgPalette
	case OBP0:
		return g.objPalette0
	case OBP1:
		return g.objPalette1
	case WX:
		return g.windowX
	case WY:
		return g.windowY
	}
	return 0x00
}

func (g *GPU) updateMode() {
	switch {
	case g.ly > constants.ScreenHeight:
		g.mode = VBlankMode
	case g.clock <= 80:
		g.mode = SearchingOAMMode
	c