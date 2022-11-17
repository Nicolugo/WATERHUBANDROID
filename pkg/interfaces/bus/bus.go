package bus

import "github.com/bokuweb/gopher-boy/pkg/types"

// Accessor bus accessor interface
type Accessor interface {
	WriteByte(addr types.Word, data byte)
	WriteWord(addr types