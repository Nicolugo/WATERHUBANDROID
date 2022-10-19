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
  0x09: ROM+RAM