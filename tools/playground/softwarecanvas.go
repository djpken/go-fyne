package playground

import (
	"github.com/djpken/go-fyne/driver/software"
	"github.com/djpken/go-fyne/test"
)

// NewSoftwareCanvas creates a new canvas in memory that can render without hardware support
func NewSoftwareCanvas() test.WindowlessCanvas {
	return software.NewCanvas()
}
