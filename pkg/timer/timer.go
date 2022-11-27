package timer

import (
	"github.com/bokuweb/gopher-boy/pkg/types"
)

const (
	// TimerRegisterOffset is register offset address
	TimerRegisterOffset types.Word = 0xFF00
	// DIV - Divider Register (R/W)
	// 