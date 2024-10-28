//go:build !windows

package glfw

import (
	"github.com/djpken/go-fyne"
	"github.com/djpken/go-fyne/internal/scale"
)

func (w *window) setDarkMode() {
}

func (w *window) computeCanvasSize(width, height int) fyne.Size {
	return fyne.NewSize(scale.ToFyneCoordinate(w.canvas, width), scale.ToFyneCoordinate(w.canvas, height))
}
