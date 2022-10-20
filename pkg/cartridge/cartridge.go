package cartridge

import (
	"fmt"
	"strings"

	"github.com/bokuweb/gopher-boy/pkg/types"
)

// Cartridge is GameBoy cartridge
type Cartridge struct {
	mbc     MBC
	Title   string
	ROM     []byte
	RAMSize int
}

/*
  CartridgeType is
  0x00: ROM ONLY
  0x01: ROM+MBC1
  0x02: ROM+MBC1+RAM
  0x03: ROM+MBC1+RAM+BATT
  0x05: ROM+MBC2
  0x06: ROM+MBC2+BATTERY
  0x08: ROM+RAM
  0x09: ROM+RAM+BATTERY
  0x0B: ROM+MMM01
  0x0C: ROM+MMM01+SRAM
  0x0D: ROM+MMM01+SRAM+BATT
  0x12: ROM+MBC3+RAM
  0x13: ROM+MBC3+RAM+BATT
  0x19: ROM+MBC5
  0x1A: ROM+MBC5+RAM
  0x1B: ROM+MBC5+RAM+BATT
  0x1C: ROM+MBC5+RUMBLE
  0x1D: ROM+MBC5+RUMBLE+SRAM
  0x1E: ROM+MBC5+RUMBLE+SRAM+BATT
  