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
	// Writing a value (XXXXXXBB - X = Don't care, B = bank select bits) into 4000-5FFF area
	// will set the two most significant ROM address lines.
	// * NOTE: The Super Smart Card doesn't require this operation because it's RAM bank is ALWAYS enabled.
	// Include this operation anyway to allow your code to work with both
	ROM16mRAM8kMode MBC1MemoryMode = "ROM16M/RAM8K"
	// ROM4mRAM32kMode is 4/32 memory mode
	// Writing a value (XXXXXXBB - X = Don't care, B = bank select bits) into 4000-5FFF area
	// will select an appropriate RAM bank at A000-C000.
	// Before you can read 