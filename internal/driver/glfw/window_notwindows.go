//go:build !windows

package glfw

import (
	"djpken/go-fyne"
	"djpken/go-fyne/internal/scale"
)

func (w *window) setDarkMode() {
}

func (w *window) computeCanvasSize(width, height int) fyne.Size {
	return fyne.NewSize(scale.ToFyneCoordinate(w.canvas, width), scale.ToFyneCoordinate(w.canvas, height))
}
