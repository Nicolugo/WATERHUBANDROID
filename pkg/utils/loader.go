package utils

import (
	"bufio"
	"os"
)

// LoadROM loads gameboy ROMfile to buf
func LoadROM(filename string) ([]byte, error) {
	file, err 