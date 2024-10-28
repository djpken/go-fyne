package common

import (
	"github.com/djpken/go-fyne"
	"github.com/djpken/go-fyne/internal/cache"
)

// CanvasForObject returns the canvas for the specified object.
func CanvasForObject(obj fyne.CanvasObject) fyne.Canvas {
	return cache.GetCanvasForObject(obj)
}
