// +build native

package main

import (
	"errors"
	"log"
	"os"

	"github.com/bokuweb/gopher-boy/pkg/interrupt"
	"github.com/bokuweb/gopher-boy/pkg/pad"

	"github.com/bokuweb/gopher-boy/pkg/gpu"
	"github.com/bokuweb/