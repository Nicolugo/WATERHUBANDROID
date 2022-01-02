// +build native

package main

import (
	"errors"
	"log"
	"os"

	"github.com/bokuweb/gopher-boy/pkg/interrupt"
	"github.com/bokuweb/gopher-boy/pkg/pad"

	"github.com/bokuweb/gopher-boy/pkg/gpu"
	"github.com/bokuweb/gopher-boy/pkg/timer"
	"github.com/bokuweb/gopher-boy/pkg/utils"

	"github.com/bokuweb/gopher-boy/pkg/cpu"
	"github.com/bokuweb/gopher-boy/pkg/gb"
	"github.com/bokuweb/gopher-boy/pkg/logger"
	"github.com/bokuweb/gopher-boy/pkg/ram"

	"github.com/bokuweb/gopher-boy/pkg/bus"
	"github.com/bokuweb/gopher-boy/pkg/cartridge"
	"github.com/bokuweb/gopher-boy/pkg/window"
)

func main() {
	level := "Debug"
	if os.Getenv("LEVEL") != "" {
		level = os.Getenv("LEVEL")
	}
	l := logger.NewLogger(logger.LogLevel(level))
	if len(os.Args) != 2 {
		log.Fatalf("ERRO