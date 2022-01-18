// +build wasm

package main

import (
	"errors"

	// "image/color"
	"log"
	"syscall/js"

	"github.com/bokuweb/gopher-boy/pkg/interrupt"
	"github.com/bokuweb/gopher-boy/pkg/logger"
	"github.com/bokuweb/gopher-boy/pkg/pad"
	"github.com/bokuweb/gopher-boy/pkg/types"
	"github.com/bokuweb/gopher-boy/pkg/window"

	"github.com/bokuweb/gopher-boy/pkg/gpu"
	"github.com/bokuweb/gopher-boy/pkg/timer"

	"github.com/bokuweb/gopher-boy/pkg/cpu"
	"github.com/bokuweb/gopher-boy/pkg/gb"
	"github.com/bokuweb/gopher-boy/pkg/ram"

	"github.com/bokuweb/gopher-boy/pkg/bus"
	"github.com/bokuweb/gopher-boy/pkg/cartridge"
)

func newGB(this js.Value, args []js.Value) interface{} {
	buf := []byte{}
	for i := 0; i < args[0].Get("length").Int(); i++ {
		buf = append(buf, byte(args[0].Index(i).Int()))
	}
	l := logger.NewLogger(logger.LogLevel("INFO"))
	cart, err := cartridge.NewCartridge(buf)
	if err != nil {
