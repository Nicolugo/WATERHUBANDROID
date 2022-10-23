package cartridge

import (
	"github.com/bokuweb/gopher-boy/pkg/ram"
	"github.com/bokuweb/gopher-boy/pkg/rom"
	"github.com/bokuweb/gopher-boy/pkg/types"
)

// MBC1 is (Memory Bank Controller 1
// MBC1 has two different maximum memory modes: 16Mbit ROM/8KByte RAM or 4Mbit ROM/32KByte RAM.
type MBC1 struct {
	rom             *rom.ROM
	ram             *ram.RAM
	selectedROMBank int
	selectedRAMBank int
	romBanking      bool
	ramEnabled      bool
	hasBattery      bool
	memoryMode      MBC1MemoryMode
	RAMSize         int
}

// MBC1MemoryMode is MBC1 max memory mode
// The MBC1 defaults to 16Mbit ROM/8KByte RAM mode
// on power up. Writing a value (XXXXXXXS - X = Don't care, S = Memory model select) into 6000-7FFF area
// will select the memory model to use.
// S = 0 selects 16/8 mode. S = 1 selects 4/32 mode.
type MBC1MemoryMode = string

const (
	// ROM16mRAM8kMode is 4/32 memory mode
	// Writing a value (XXXXXXBB - X = Don'