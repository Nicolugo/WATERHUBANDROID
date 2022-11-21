package mocks

import (
	"github.com/bokuweb/gopher-boy/pkg/types"
)

type MockBus struct {
	MockMemory [0x10000]byte
}

func (b *MockBus) WriteByte(addr types.Word, data byte) {
	b.MockMemory[addr] = data
}

func (b *MockBus) WriteWord(addr types.Word, data types.Word) {
	b.MockMemory[addr] = b